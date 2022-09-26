package comp

import (
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/utils"
)

var preInitializerMetaType = MetaType[PreInitializer]()
var postInitializerMetaType = MetaType[PostInitializer]()
var preComposerMetaType = MetaType[PreComposer]()
var postComposerMetaType = MetaType[PostComposer]()

// Composer creates every node and config their properties based on the collected info, then
// link and negotiate for them to put them into working state.
type Composer struct {
	sessionId string

	gt             *nmd.GraphTopology
	nodeSortedList []SessionAware // topographical sorted nodes, first one has no receiver
	nodeMap        map[string]SessionAware

	initiator  CommandInitiator
	linkPoints []LinkPoint
}

func NewSessionComposer(sessionId string) *Composer {
	sc := &Composer{
		sessionId: sessionId,
		nodeMap:   make(map[string]SessionAware),
	}
	sc.initiator = &builtinCommandInitiator{composer: sc}
	return sc
}

func (c *Composer) IterateNode(iter func(name string, node SessionAware)) {
	for name, node := range c.nodeMap {
		iter(name, node)
	}
}

func (c *Composer) GetNode(name string) SessionAware {
	return c.nodeMap[name]
}

func (c *Composer) ParseGraphDescription(desc string) (err error) {
	gt := nmd.NewGraphTopology()
	err = gt.ParseGraph(c.sessionId, desc)
	c.gt = gt
	return
}

func (c *Composer) Connect(sender, receiver SessionAware, preferredOffer []MessageType) (lp LinkPoint, err error) {
	receiverSession := receiver.GetNodeScope()
	receiverName := receiver.GetNodeName()

	if preferredOffer == nil {
		preferredOffer = sender.Offer()
	}

	lp, err = sender.StreamTo(receiverSession, receiverName, preferredOffer)
	return
}

func (c *Composer) preConnectNodes() error {
	for _, node := range c.nodeSortedList {
		utils.AopCall(node, []interface{}{c, node}, preComposerMetaType, "BeforeCompose")
	}
	return nil
}

func (c *Composer) postConnectNodes() error {
	for _, node := range c.nodeSortedList {
		utils.AopCall(node, []interface{}{c, node}, postComposerMetaType, "AfterCompose")
	}
	return nil
}

func (c *Composer) connectNodes(graph *event.Graph) (lps []LinkPoint, err error) {
	nodeDefs := c.gt.GetSortedNodeDefs()
	// add all nodes to graph
	var lp LinkPoint
	for _, node := range c.nodeSortedList {
		if !graph.AddNode(node) {
			err = fmt.Errorf("failed to add node %v to graph", node.GetNodeName())
			return
		}
	}
	// all nodes are added but not connected, as connection is built only when sender and receiver agreed on the
	// message type of their link. however, for some types of node (pubsub), what message type they output depends on
	// their input message type, which means message types propagate from senders(at the end of sorted list) to
	// receivers(at the start of sorted list), so iterate the sorted list reversely until all message types are
	// determined as required by link creation
	for i := len(c.nodeSortedList) - 1; i >= 0; i-- {
		sender := c.nodeSortedList[i]
		for _, receiverDef := range nodeDefs[i].Deps {
			receiver := c.nodeMap[receiverDef.Name]
			// TODO: nmd language add support for specifying preferred offer
			// TODO: use ssa to analyze Offer() of every node, check the message type is statically or dynamically defined
			if lp, err = c.Connect(sender, receiver, nil); err != nil {
				return
			} else {
				lps = append(lps, lp)
			}
		}
	}
	return
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
		utils.AopCall(sn, nil, preInitializerMetaType, "PreInit")
		if err = sn.Init(); err != nil {
			return
		}
		utils.AopCall(sn, nil, postInitializerMetaType, "PostInit")

		c.nodeSortedList = append(c.nodeSortedList, sn)
		c.nodeMap[sn.GetNodeName()] = sn
		id := NewId(sn.GetNodeScope(), sn.GetNodeName())
		nodeIds = append(nodeIds, id)
	}

	if err = c.preConnectNodes(); err != nil {
		return
	}
	if c.linkPoints, err = c.connectNodes(graph); err != nil {
		return
	}
	if err = c.postConnectNodes(); err != nil {
		return
	}

	return
}

func (c *Composer) GetSortedNodes() (ni []*nmd.NodeDef) {
	return c.gt.GetSortedNodeDefs()
}

func (c *Composer) GetCommandInitiator() CommandInitiator {
	return c.initiator
}

func (c *Composer) ExitGraph() {
	for _, n := range c.nodeSortedList {
		n.ExitGraph()
	}
}
