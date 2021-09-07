package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/utils"
	"strings"
	"time"
)

type Composer struct {
	sessionId string
	gt        *GraphTopology
	entryNode *EntryNode // currently, only support one entry
	sendNode  *Dispatch

	namedChannel map[string]chan<- *event.Event
}

func NewSessionComposer(sessionId string) *Composer {
	sc := &Composer{
		sessionId:    sessionId,
		namedChannel: make(map[string]chan<- *event.Event),
	}
	return sc
}

func (c *Composer) ParseGraphDescription(desc string) (err error) {
	gt := newGraphTopology()
	lines := strings.Split(desc, "\n")
	for _, l := range lines {
		if l == "" {
			continue
		}
		gt.parseLine(l)
	}
	if gt.nbParseError > 0 {
		errStr := fmt.Sprintf("there are total %v error in graph description:\n%v", gt.nbParseError, desc)
		err = errors.New(errStr)
		return
	}
	err = gt.topographicalSort()
	c.gt = gt
	return
}

func (c *Composer) GetSortedNodes() (ni []*NodeInfo) {
	return c.gt.sortedNodeList
}

// PrepareNodes create node instances by type, add them to graph, and link them
func (c *Composer) PrepareNodes(graph *event.EventGraph) (err error) {
	var nodeList []SessionAware
	var nodeIds []*Id
	var sendNode *Dispatch
	nbNode := len(c.gt.sortedNodeList)

	defer func() {
		if err != nil {
			// undo AddNode
			for _, n := range nodeList {
				n.ExitGraph()
			}
			if sendNode != nil {
				sendNode.ExitGraph()
			}
		}
	}()

	// create node instances
	for _, n := range c.gt.sortedNodeList {
		n.Props.Set("Name", n.Name)
		sn := MakeSessionNode(n.Name, c.sessionId, n.Props)
		if sn == nil {
			logger.Errorf("unknown node type: %v\n", n.Name)
			err = errors.New("can not make unknown node")
			return
		}
		nodeList = append(nodeList, sn)
		if n.Name == TYPE_ENTRY {
			c.entryNode = sn.(*EntryNode)
		}
		id := NewId(sn.GetNodeScope(), sn.GetNodeName())
		nodeIds = append(nodeIds, id)
	}

	// add all nodes to graph, take action until all of them are ready
	for _, n := range nodeList {
		if !graph.AddNode(n) {
			err = errors.New(fmt.Sprintf("failed to add node %v to graph", n.GetNodeName()))
			return
		}
	}

	// now every node is added to graph, build link between them:
	// first create send node which links to all nodes in the session,
	// then send route command through send node
	ci := make(ConfigItems)
	sendNode = MakeSessionNode(TYPE_SEND, c.sessionId, ci).(*Dispatch)
	sendNode.SetMaxLink(nbNode) // enlarge it if we need supporting dynamically add-node
	if !graph.AddNode(sendNode) {
		err = errors.New("fail to add send-node to graph")
		return
	}
	if err = sendNode.ConnectTo(nodeIds); err != nil {
		return
	}

	nbLink := len(c.gt.getLinkInfo())
	linkReadyC := make(chan bool, nbLink)
	linkReadyFunc := func() { linkReadyC <- true }
	for _, li := range c.gt.getLinkInfo() {
		from, to := li.From, li.To
		cmd := &GenericRouteCommand{}
		cmd.SessionId = c.sessionId
		cmd.Name = to.Name
		evt := event.NewEventWithCallback(CMD_GENERIC_SET_ROUTE, cmd, linkReadyFunc)
		sendNode.SendTo(c.sessionId, from.Name, evt) // set From-node's route to To-node
	}
	// wait until all set route command executed
	if err = utils.WaitChannelWithTimeout(linkReadyC, nbLink, 2*time.Second); err != nil {
		errStr := "waiting for node ready failed when setting route"
		logger.Errorln(errStr)
		err = errors.New(errStr)
		return
	}
	c.sendNode = sendNode

	// subscribe channels
	if len(c.namedChannel) > 0 {
		for i, n := range c.gt.sortedNodeList {
			if n.Type == TYPE_PUBSUB {
				if name, ok := n.Props["channel"]; ok {
					chName, ok1 := name.(string)
					if !ok1 {
						break
					}
					if ch, exist := c.namedChannel[chName]; exist {
						nodeList[i].(*PubSubNode).SubscribeChannel(chName, ch)
					}
				}
			}
		}
	}
	return
}

func (c *Composer) GetEntry() *EntryNode {
	return c.entryNode
}

func (c *Composer) RegisterChannel(name string, ch chan<- *event.Event) {
	c.namedChannel[name] = ch
}
