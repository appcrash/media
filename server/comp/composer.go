package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"strings"
	"sync"
)

// we create every node and config their properties based on the collected info, then
// link them and send initial events as described before putting them into working state.

type Composer struct {
	sessionId       string
	gt              *nmd.GraphTopology
	messageProvider []MessageProvider
	dispatch        *Dispatch
	nodeList        []SessionAware // topographical sorted nodes

	// channel handling, can static or dynamically add channels
	mutex              sync.Mutex
	pendingChannelNode map[string]*PubSubNode
	namedChannel       map[string]chan<- *event.Event
}

func NewSessionComposer(sessionId string) *Composer {
	sc := &Composer{
		sessionId:          sessionId,
		pendingChannelNode: make(map[string]*PubSubNode),
		namedChannel:       make(map[string]chan<- *event.Event),
	}
	return sc
}

func (c *Composer) ParseGraphDescription(desc string) (err error) {
	gt := nmd.NewGraphTopology()
	err = gt.ParseGraph(c.sessionId, desc)
	c.gt = gt
	return
}

// PrepareNodes create node instances by type, add them to graph, and link them
func (c *Composer) PrepareNodes(graph *event.Graph) (err error) {
	var nodeIds []*Id
	nodeDefs := c.gt.GetSortedNodeDefs()
	nbNode := len(nodeDefs)

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
			err = errors.New("can not make unknown node")
			return
		}
		if err = sn.Init(); err != nil {
			return
		}
		c.nodeList = append(c.nodeList, sn)
		if provider, ok := sn.(MessageProvider); ok {
			c.messageProvider = append(c.messageProvider, provider)
		}

		id := NewId(sn.GetNodeScope(), sn.GetNodeName())
		nodeIds = append(nodeIds, id)
	}

	// add all nodes to graph, create links between them, as nodes are already topographical sorted,
	// for each node, its dependent nodes are in graph when adding it to graph
	for i, n := range c.nodeList {
		if !graph.AddNode(n) {
			err = errors.New(fmt.Sprintf("failed to add node %v to graph", n.GetNodeName()))
			return
		}
		deps := nodeDefs[i].Deps
		for _, ni := range deps {
			// set pipe end to local session nodes
			if n.SetPipeOut(c.sessionId, ni.Name) != nil {
				err = errors.New(fmt.Sprintf("failed to link %v => %v", n.GetNodeName(), ni.Name))
				return
			}
		}
	}

	// now every node is added to graph and linked
	// create dispatch node which links to all nodes in the session
	c.dispatch = MakeSessionNode(TYPE_DISPATCH, c.sessionId, nil).(*Dispatch)
	c.dispatch.SetMaxLink(nbNode * 2) // reserved nbNode for dynamical link requests
	if !graph.AddNode(c.dispatch) {
		err = errors.New("fail to add send-node to graph")
		return
	}
	if err = c.dispatch.connectTo(nodeIds); err != nil {
		return
	}

	// again, let all nodes reference this dispatch
	for _, n := range c.nodeList {
		n.SetController(c.dispatch)
	}

	// subscribe channels, for all nodes of type pubsub, find the registered channel with same name
	// as specified in pubsub's "channel" property
	if len(c.namedChannel) > 0 {
		for i, n := range nodeDefs {
			if n.Type != TYPE_PUBSUB {
				continue
			}
			var chNameList string
			var ok bool
			for _, p := range n.Props {
				if p.Key == "channel" {
					if p.Type != "str" || p.Value == nil {
						err = errors.New(fmt.Sprintf("pubsub channel value is not string: %v", p.Value))
						return
					}
					if chNameList, ok = p.Value.(string); ok {
					} else {
						err = errors.New(fmt.Sprintf("pubsub channel value can not converted to string: %v", p.Value))
						return
					}
					break
				}
			}
			if chNameList == "" {
				logger.Debugf("pubsub channel props is nil")
				continue
			}

			psNode := c.nodeList[i].(*PubSubNode)
			// pubsub property, for example: channel=a,b,c ...
			logger.Println("chNameList ", chNameList)
			for _, chName := range strings.Split(chNameList, ",") {
				if ch, exist := c.namedChannel[chName]; exist {
					// the channel is already registered, just subscribe it to pubsub node now
					if err = psNode.SubscribeChannel(chName, ch); err != nil {
						return
					}
				} else {
					// the channel is required by pubsub, but has not registered yet, add it to pending list,
					// so we can use composer to register it later dynamically
					c.pendingChannelNode[chName] = psNode
				}
			}
		}
	}

	return
}

func (c *Composer) GetSortedNodes() (ni []*nmd.NodeDef) {
	return c.gt.GetSortedNodeDefs()
}

// GetMessageProvider get entry by its name
func (c *Composer) GetMessageProvider(name string) MessageProvider {
	for _, provider := range c.messageProvider {
		if provider.GetName() == name {
			return provider
		}
	}
	return nil
}

func (c *Composer) GetController() Controller {
	return c.dispatch
}

func (c *Composer) RegisterChannel(name string, ch chan<- *event.Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// check if the channel is in pending list, if found, a pubsub node is waiting for it
	if psNode, exist := c.pendingChannelNode[name]; exist {
		psNode.SubscribeChannel(name, ch)
		delete(c.pendingChannelNode, name)
	} else {
		// statically register this channel to used in PrepareNodes phase
		c.namedChannel[name] = ch
	}
}

func (c *Composer) ExitGraph() {
	for _, n := range c.nodeList {
		n.ExitGraph()
	}
	if c.dispatch != nil {
		c.dispatch.ExitGraph()
	}
}
