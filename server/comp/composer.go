package comp

import (
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
)

// Composer creates every node and config their properties based on the collected info, then
// link and negotiate for them to put them into working state.
type Composer struct {
	sessionId      string
	gt             *nmd.GraphTopology
	nodeSortedList []SessionAware // topographical sorted nodes
	nodeMap        map[string]SessionAware

	linkPoints []LinkPoint
}

func NewSessionComposer(sessionId string) *Composer {
	sc := &Composer{
		sessionId: sessionId,
	}
	return sc
}

func (c *Composer) ParseGraphDescription(desc string) (err error) {
	gt := nmd.NewGraphTopology()
	err = gt.ParseGraph(c.sessionId, desc)
	c.gt = gt
	return
}

func (c *Composer) Connect(sender, receiver SessionAware) (lp LinkPoint, err error) {
	receiverSession := receiver.GetNodeScope()
	receiverName := receiver.GetNodeName()

	lp, err = sender.StreamTo(receiverSession, receiverName)
	return
}

func (c *Composer) preConnectNodes() {
	for _, node := range c.nodeSortedList {
		trait, _ := NodeTraitOfType(node.GetNodeTypeName())
		if trait.IsPreComposer() {
			node.(PreComposer).BeforeCompose(c)
		}
	}
}

func (c *Composer) postConnectNodes() {
	for _, node := range c.nodeSortedList {
		trait, _ := NodeTraitOfType(node.GetNodeTypeName())
		if trait.IsPostComposer() {
			node.(PostComposer).AfterCompose(c)
		}
	}
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
	// determined as required to link creation
	for i := len(c.nodeSortedList) - 1; i >= 0; i-- {
		sender := c.nodeSortedList[i]
		for _, receiverDef := range nodeDefs[i].Deps {
			receiver := c.nodeMap[receiverDef.Name]
			if lp, err = c.Connect(sender, receiver); err != nil {
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
		if err = sn.Init(); err != nil {
			return
		}
		c.nodeSortedList = append(c.nodeSortedList, sn)
		c.nodeMap[sn.GetNodeName()] = sn

		id := NewId(sn.GetNodeScope(), sn.GetNodeName())
		nodeIds = append(nodeIds, id)
	}

	c.preConnectNodes()
	if c.linkPoints, err = c.connectNodes(graph); err != nil {
		return
	}
	c.postConnectNodes()

	return
}

func (c *Composer) GetSortedNodes() (ni []*nmd.NodeDef) {
	return c.gt.GetSortedNodeDefs()
}

func (c *Composer) GetController() CommandInitiator {
	return c
}

func (c *Composer) ExitGraph() {
	for _, n := range c.nodeSortedList {
		n.ExitGraph()
	}
}

func (c *Composer) Call(fromNode, toNode string, args []string) (resp []string) {
	if to, ok := c.nodeMap[toNode]; ok {
		resp = to.OnCall(fromNode, args)
	} else {
		resp = WithError("no such node")
	}
	return
}

func (c *Composer) Cast(fromNode, toNode string, args []string) {
	if to, ok := c.nodeMap[toNode]; ok {
		to.OnCast(fromNode, args)
	}
}
