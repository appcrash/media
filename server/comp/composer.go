package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"sync"
)

// Composer creates every node and config their properties based on the collected info, then
// link them and send initial events as described before putting them into working state.
type Composer struct {
	sessionId      string
	gt             *nmd.GraphTopology
	nodeSortedList []SessionAware // topographical sorted nodes
	nodeMap        map[string]SessionAware

	// channel handling, channels are registered by parsing nmd graph description, then linked dynamically
	mutex        sync.Mutex
	namedChannel map[string]*channelInfo
}

const (
	maxNegotiationIteration = 5
)

type channelInfo struct {
	isConnected bool
	peerNode    *PubSubNode
	ch          chan<- *event.Event
}

func NewSessionComposer(sessionId string) *Composer {
	sc := &Composer{
		sessionId:    sessionId,
		namedChannel: make(map[string]*channelInfo),
	}
	return sc
}

func (c *Composer) ParseGraphDescription(desc string) (err error) {
	gt := nmd.NewGraphTopology()
	err = gt.ParseGraph(c.sessionId, desc)
	c.gt = gt
	return
}

func (c *Composer) Connect(sender, receiver SessionAware) (llp, rlp LinkPoint, err error) {
	receiverName := receiver.GetNodeName()
	outputTraits := sender.ProvideOffer()
	commonTrait := receiver.AnswerOffer(outputTraits)
	if len(commonTrait) == 0 {
		err = fmt.Errorf("can not create link between %v => %v, no common trait",
			sender.GetNodeName(), receiver.GetNodeName())
		return
	}
	trait := commonTrait[0]
	if err, llp = sender.StreamTo(receiverName); err != nil {
		return
	}
	if err, rlp = receiver.StreamBy(sender.GetNodeName()); err != nil {
		return
	}
	llp.SetPeer(rlp)
	rlp.SetPeer(llp)

	return
}

func (c *Composer) connectNodes(graph *event.Graph) (lps []LinkPoint, err error) {
	nodeDefs := c.gt.GetSortedNodeDefs()
	// add all nodes to graph, create links between them, as nodes are already topographical sorted,
	// for each node, its dependent nodes are in graph when adding it to graph
	var llp, rlp LinkPoint
	for i, sender := range c.nodeSortedList {
		if !graph.AddNode(sender) {
			err = fmt.Errorf("failed to add node %v to graph", sender.GetNodeName())
			return
		}
		deps := nodeDefs[i].Deps
		for _, receiverDef := range deps {
			// just create link point for each other
			receiver := c.nodeMap[receiverDef.Name]
			if llp, rlp, err = c.Connect(sender, receiver); err != nil {
				return
			} else {
				lps = append(lps, llp, rlp)
			}
		}
	}
	return
}

func (c *Composer) Negotiate(lps []LinkPoint) (unresolved []LinkPoint, err error) {
	var iteration int
	for iteration < maxNegotiationIteration {

	}
}

// ComposeNodes create node instances by type, add them to graph, and link them
func (c *Composer) ComposeNodes(graph *event.Graph) (err error) {
	var nodeIds []*Id
	nodeDefs := c.gt.GetSortedNodeDefs()

	defer func() {
		if err != nil {
			// undo AddNode
			c.ExitGraph()
		}
	}()

	// create node instances, collect message providers if any
	for _, n := range nodeDefs {
		n.Props = append(n.Props, &nmd.NodeProp{
			Key:   "Name",
			Type:  "str",
			Value: n.Name,
		})
		sn := MakeSessionNode(n.Type, c.sessionId, n.Props)
		if sn == nil {
			logger.Errorf("unknown node type: %v\n", n.Name)
			err = fmt.Errorf("can not make unknown node: %v", n.Type)
			return
		}
		if err = sn.Init(); err != nil {
			return
		}
		c.nodeSortedList = append(c.nodeSortedList, sn)
		c.nodeMap[sn.GetNodeName()] = sn

		id := NewId(sn.GetNodeScope(), sn.GetNodeName())
		nodeIds = append(nodeIds, id)
	}

	var lps []LinkPoint
	if lps, err = c.connectNodes(graph); err != nil {
		return
	}

	// again, let all nodes reference this dispatch
	for _, n := range c.nodeSortedList {
		n.SetController(c)
	}

	// subscribe channels, for all nodes of type pubsub, find the registered channel with same Name
	// as specified in pubsub's "channel" property
	//for i, n := range nodeDefs {
	//	if n.Type != TypePUBSUB {
	//		continue
	//	}
	//	var chNameList string
	//	var ok bool
	//	for _, p := range n.Props {
	//		if p.Key == "channel" {
	//			if p.Type != "str" || p.Value == nil {
	//				err = fmt.Errorf("pubsub channel value is not string: %v", p.Value)
	//				return
	//			}
	//			if chNameList, ok = p.Value.(string); ok {
	//			} else {
	//				err = fmt.Errorf("pubsub channel value can not converted to string: %v", p.Value)
	//				return
	//			}
	//			break
	//		}
	//	}
	//	if chNameList == "" {
	//		logger.Debugf("pubsub channel props is nil")
	//		continue
	//	}
	//
	//	psNode := c.nodeSortedList[i].(*PubSubNode)
	//	// pubsub property, for example: channel=a,b,c ...
	//	logger.Debugln("chNameList ", chNameList)
	//	for _, chName := range strings.Split(chNameList, ",") {
	//		if _, exist := c.namedChannel[chName]; exist {
	//			// the channel is already registered
	//			err = fmt.Errorf("channel:%v can only subscribe to one pubsub node", chName)
	//			return
	//		} else {
	//			c.namedChannel[chName] = &channelInfo{peerNode: psNode}
	//		}
	//	}
	//}

	return
}

func (c *Composer) GetSortedNodes() (ni []*nmd.NodeDef) {
	return c.gt.GetSortedNodeDefs()
}

func (c *Composer) GetController() CommandInitiator {
	return c
}

func (c *Composer) LinkChannel(name string, ch chan<- *event.Event) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// check if the channel is connected, if found and not connected, a pubsub node is waiting for it
	if ci, exist := c.namedChannel[name]; exist {
		if ci.isConnected {
			return errors.New("channel is already connected")
		}
		if err := ci.peerNode.SubscribeChannel(name, ch); err != nil {
			return err
		}
		ci.ch = ch
		ci.isConnected = true
		return nil
	}
	return fmt.Errorf("no such channel: %v", name)
}

func (c *Composer) UnlinkChannel(name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if ci, exist := c.namedChannel[name]; exist {
		if !ci.isConnected {
			return errors.New("channel is not connected")
		}
		if err := ci.peerNode.UnsubscribeChannel(name); err != nil {
			return err
		}
		ci.ch = nil
		ci.isConnected = false
		return nil
	}
	return fmt.Errorf("no such channel: %v", name)
}

func (c *Composer) ExitGraph() {
	for _, n := range c.nodeSortedList {
		n.ExitGraph()
	}
}

func (c *Composer) Call(fromNode, toNode string, args []string) (resp []string) {
	if to, ok := c.nodeMap[toNode]; ok {
		resp = to.OnCall(c.sessionId, fromNode, args)
	} else {
		resp = WithError("no such node")
	}
	return
}

func (c *Composer) Cast(fromNode, toNode string, args []string) {
	if to, ok := c.nodeMap[toNode]; ok {
		to.OnCast(c.sessionId, fromNode, args)
	}
}
