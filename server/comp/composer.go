package comp

import (
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/utils"
	"reflect"
	"strings"
)

var (
	preComposerType        = MetaType[PreComposer]()
	postComposerType       = MetaType[PostComposer]()
	initializingNodeType   = MetaType[InitializingNode]()
	unInitializingNodeType = MetaType[UnInitializingNode]()
	linkPointType          = MetaType[LinkPoint]()
)

const (
	structTagKey = "comp"
)

// Composer creates every node and config their properties based on the collected info, then
// link and negotiate for them to put them into working state. the order of function called when no error happens:
//
// Node.New(created by factory method and injected by nmd properties)
// Node.Init (aop)
// Node.BeforeCompose (aop)
// Node.OnEnter
// Node.AfterCompose (aop)
// Node.UnInit (aop)
// Node.OnExit
type Composer struct {
	sessionId  string
	instanceId string

	graphTopo      *nmd.GraphTopology
	nodeSortedList []SessionAware // topographical sorted nodes, first one has no receiver
	nodeMap        map[string]SessionAware

	initiator  CommandInitiator
	linkPoints []LinkPoint
	nodeExited bool // ensure node UnInit called only once
}

func NewSessionComposer(sessionId, instanceId string) *Composer {
	sc := &Composer{
		sessionId:  sessionId,
		instanceId: instanceId,
		nodeMap:    make(map[string]SessionAware),
	}
	sc.initiator = &builtinCommandInitiator{composer: sc}
	return sc
}

func (c *Composer) IterateNode(iter func(name string, node SessionAware)) {
	for name, node := range c.nodeMap {
		iter(name, node)
	}
}

func (c *Composer) GetSessionId() string {
	return c.sessionId
}

func (c *Composer) GetInstanceId() string {
	return c.instanceId
}

func (c *Composer) GetNode(name string) SessionAware {
	if sa, exist := c.nodeMap[name]; exist {
		return sa
	} else {
		return nil
	}
}

func (c *Composer) ParseGraphDescription(desc string) (err error) {
	gt := nmd.NewGraphTopology()
	err = gt.ParseGraph(c.sessionId, desc, filterGatewayNode)
	c.graphTopo = gt
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

// currently, only allow gateway nodes to create loops in the graph, all you have to do is NAME the node with
// "gateway" suffix. gateway node diffs from default filter node by acting as a message exchanger to outside,
// so a node or a sub-graph can write message to gateway and read from it at the same time.
func filterGatewayNode(nodeName string) bool {
	return strings.HasSuffix(strings.ToLower(nodeName), "gateway")
}

func (c *Composer) preConnectNodes() error {
	for _, node := range c.nodeSortedList {
		utils.AopCall(node, []interface{}{c, node}, preComposerType, "BeforeCompose")
	}
	return nil
}

func (c *Composer) postConnectNodes() error {
	for _, node := range c.nodeSortedList {
		utils.AopCall(node, []interface{}{c, node}, postComposerType, "AfterCompose")
	}
	return nil
}

func (c *Composer) connectNodes(graph *event.Graph) (lps []LinkPoint, err error) {
	nodeDefs := c.graphTopo.GetSortedNodeDefs()
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
		var slps []LinkPoint = nil
		sender := c.nodeSortedList[i]
		for _, receiverDef := range nodeDefs[i].Deps {
			receiver := c.nodeMap[receiverDef.Name]
			// TODO: nmd language add support for specifying preferred offer
			// TODO: use ssa to analyze Offer() of every node, check the message type is statically or dynamically defined
			if lp, err = c.Connect(sender, receiver, nil); err != nil {
				return
			} else {
				lps = append(lps, lp)
				slps = append(slps, lp)
			}
		}
		// the sender has created link points, check the node's field and try to inject them to field variables
		value := reflect.ValueOf(sender).Elem()
		typ := value.Type()
		for j := 0; j < typ.NumField(); j++ {
			field := typ.Field(j)
			if field.Type == linkPointType {
				if tag, exist := field.Tag.Lookup(structTagKey); exist {
					if err = injectLinkPoint(field, value.Field(j), slps, tag); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

// ComposeNodes create node instances by type, add them to graph, and link them
func (c *Composer) ComposeNodes(graph *event.Graph) (err error) {
	var nodeIds []*Id
	var ok bool
	nodeDefs := c.graphTopo.GetSortedNodeDefs()

	defer func() {
		if err != nil {
			// undo AddNode
			c.ExitGraph()
		}
	}()

	// create node instances
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
		returnValues := utils.AopCall(sn, nil, initializingNodeType, "Init")
		for _, rv := range returnValues {
			// any Init error would prevent composing
			initRetValue := rv[0]
			if !initRetValue.IsNil() {
				if initRetValue.IsValid() && initRetValue.CanInterface() {
					if err, ok = initRetValue.Interface().(error); !ok {
						err = fmt.Errorf("node init error: %v", initRetValue)
					}
				} else {
					err = fmt.Errorf("node init problem: %v", initRetValue)
				}
				logger.Errorf("node %v init failed with error: %v", sn, initRetValue)
				return
			}
		}

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
	return c.graphTopo.GetSortedNodeDefs()
}

func (c *Composer) GetCommandInitiator() CommandInitiator {
	return c.initiator
}

func (c *Composer) ExitGraph() {
	if c.nodeExited {
		return
	}
	for _, n := range c.nodeSortedList {
		utils.AopCall(n, nil, unInitializingNodeType, "UnInit")
	}
	c.nodeExited = true
}
