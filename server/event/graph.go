package event

import (
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/utils"
	"reflect"
	"time"
)

const (
	graphAddNodeTimeout = 5 * time.Second
)

type scopeMapType map[string][]*NodeDelegate
type linkSetType map[string]bool
type nodeMapType map[string]*nodeInfo

type Graph struct {
	scopeMap scopeMapType
	nodeMap  nodeMapType // nodeId -> nodeInfo
	linkSet  linkSetType // links that still alive

	eventChannel chan *Event
}

type nodeInfo struct {
	inputLinks  []*dlink
	outputLinks []*dlink
	maxLink     int
}

func (eg *Graph) updateNodeStats() {
	nb := float64(len(eg.nodeMap))
	prom.NodeGraphNodes.Set(nb)
}

func (eg *Graph) updateLinkStats() {
	nb := float64(len(eg.linkSet))
	prom.NodeGraphLinks.Set(nb)
}

func (eg *Graph) findNode(scope string, name string) *NodeDelegate {
	if nodeList, ok := eg.scopeMap[scope]; ok {
		for _, node := range nodeList {
			if node.getNodeName() == name {
				return node
			}
		}
	}
	return nil
}

func (eg *Graph) addNode(nd *NodeDelegate, maxLink int) {
	scope := nd.getNodeScope()
	if nodeList, ok := eg.scopeMap[scope]; !ok {
		eg.scopeMap[scope] = []*NodeDelegate{nd}
	} else {
		eg.scopeMap[scope] = append(nodeList, nd)
	}
	eg.nodeMap[nd.getId()] = &nodeInfo{maxLink: maxLink}
	eg.updateNodeStats()
}

func (eg *Graph) delNode(nd *NodeDelegate) {
	scope := nd.getNodeScope()
	if nodeList, ok := eg.scopeMap[scope]; ok {
		index := -1
		for i, node := range nodeList {
			if node == nd {
				index = i
				break
			}
		}
		length := len(nodeList)
		if index >= 0 {
			if length > 1 {
				nodeList[index] = nodeList[length-1]
				eg.scopeMap[scope] = nodeList[:length-1]
			} else {
				// if length == 1, just remove entire array
				delete(eg.scopeMap, scope)
			}
		}
	}
	delete(eg.nodeMap, nd.getId())
	eg.updateNodeStats()
}

func (eg *Graph) getNodeInfo(nodeId string) *nodeInfo {
	if info, exist := eg.nodeMap[nodeId]; exist {
		return info
	}
	return nil
}

// associate dlink to a node delegate
func (eg *Graph) addLink(nd *NodeDelegate, l *dlink) {
	info := eg.getNodeInfo(nd.getId())
	if info == nil {
		return
	}
	if l.fromNode == nd {
		l.fromIndex = len(info.outputLinks)
		info.outputLinks = append(info.outputLinks, l)
	} else if l.toNode == nd {
		l.toIndex = len(info.inputLinks)
		info.inputLinks = append(info.inputLinks, l)
	} else {
		logger.Errorf("[graph]: can not addLink %v\n", l)
	}
	eg.updateLinkStats()
}

// tear down a dlink in node delegate
func (eg *Graph) delLink(nd *NodeDelegate, l *dlink) {
	info := eg.getNodeInfo(nd.getId())
	if info == nil {
		return
	}
	var linkArray *[]*dlink
	var index int
	var length int
	if l.fromNode == nd {
		if l.fromIndex >= len(info.outputLinks) {
			// wrong index, out of range
			return
		}
		linkArray = &info.outputLinks
		length = len(*linkArray)
		index = l.fromIndex
		info.outputLinks[length-1].fromIndex = l.fromIndex
		l.fromIndex = -1
	} else if l.toNode == nd {
		if l.toIndex >= len(info.inputLinks) {
			// wrong index, out of range
			return
		}
		linkArray = &info.inputLinks
		length = len(*linkArray)
		index = l.toIndex
		info.inputLinks[length-1].toIndex = l.toIndex
		l.toIndex = -1
	} else {
		logger.Errorf("[graph]: can not delete link: %v\n", l)
		return
	}
	// remove dlink by exchanging with last one
	lastLink := (*linkArray)[length-1]
	(*linkArray)[index] = lastLink
	*linkArray = (*linkArray)[:length-1]
	eg.updateLinkStats()
}

func (eg *Graph) deliveryEvent(evt *Event) {
	eg.eventChannel <- evt
}

// simply loop forever
func (eg *Graph) startEventLoop(c chan int) {
	go func(g *Graph) {
		c <- 0
		for {
			evt := <-g.eventChannel
			g.onEvent(evt)
		}
	}(eg)
}

func (eg *Graph) onEvent(evt *Event) {
	var ok bool
	switch evt.cmd {
	case reqNodeAdd:
		var req *nodeAddRequest
		if req, ok = evt.obj.(*nodeAddRequest); !ok {
			return
		}
		eg.onAddNode(req)
	case reqNodeExit:
		var req *nodeExitRequest
		if req, ok = evt.obj.(*nodeExitRequest); !ok {
			return
		}
		eg.onExitNode(req)
	case reqLinkUp:
		var req *linkUpRequest
		if req, ok = evt.obj.(*linkUpRequest); !ok {
			return
		}
		eg.onLinkUp(req)
	case reqLinkDown:
		var req *linkDownRequest
		if req, ok = evt.obj.(*linkDownRequest); !ok {
			return
		}
		eg.onLinkDown(req)
	}
}

// add node to graph, and send node-add response to this node immediately
func (eg *Graph) onAddNode(req *nodeAddRequest) {
	maxLink := defaultMaxLink
	node := req.node
	ps := reflect.ValueOf(node)

	elem := ps.Elem()
	if elem.Kind() == reflect.Struct {
		field := elem.FieldByName("maxLink") // CAVEAT: change the name once NodeProperty field change accordingly
		if field.IsValid() {
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if ml := int(field.Int()); ml > 0 {
					maxLink = ml
				}
			}
		}
	}
	delegate := newNodeDelegate(eg, node, maxLink)
	eg.addNode(delegate, maxLink)
	// all gears up, rock it
	go func(nd *NodeDelegate, cb Callback) {
		nd.startEventLoop()
		resp := newNodeAddResponse(nd, cb)
		nd.receiveCtrl(resp)
	}(delegate, req.cb)

}

// node requests exiting the graph, notify all senders linking to this node
func (eg *Graph) onExitNode(req *nodeExitRequest) {
	nd := req.delegate
	nodeInfo := eg.getNodeInfo(nd.getId())
	if nodeInfo == nil {
		// concurrently call RequestNodeExit can have more than one
		// exit request, but it is ok as only one exit response would
		// be sent
		return
	}

	// tear down all input/output links of the node, but only notify
	// senders link-down state, as for receivers just remove links
	// from their nodeInfo without any notification
	for _, link := range nodeInfo.inputLinks {
		if _, exist := eg.linkSet[link.name]; exist {
			link.fromNode.receiveCtrl(newLinkDownResponse(stateSuccess, link))
			eg.delLink(link.fromNode, link)
			eg.delLink(link.toNode, link)
			delete(eg.linkSet, link.name)
		}
	}
	for _, link := range nodeInfo.outputLinks {
		if _, exist := eg.linkSet[link.name]; exist {
			eg.delLink(link.fromNode, link)
			eg.delLink(link.toNode, link)
			delete(eg.linkSet, link.name)
		}
	}
	if len(nodeInfo.inputLinks) > 0 || len(nodeInfo.outputLinks) > 0 {
		logger.Errorf("[node]: (%v) is exiting but have inputLinks:%v,outputLinks:%v\n",
			nd.getNodeName(), len(nodeInfo.inputLinks), len(nodeInfo.outputLinks))
		panic("node still have active links")
	}
	eg.delNode(nd)
	// finally, send the last ctrl message for this node
	nd.receiveCtrl(newNodeExitResponse())
}

// request dlink to other node, decline if that node is exiting or
// dlink is duplicated, otherwise create dlink between them
func (eg *Graph) onLinkUp(req *linkUpRequest) {
	nodeName := req.nodeName
	scope := req.scope
	fromNode := req.fromNode

	ni := eg.nodeMap[fromNode.getId()]
	if ni.maxLink == len(ni.outputLinks) {
		req.fromNode.receiveCtrl(newLinkUpResponse(nil, stateNodeExceedMaxLink, scope, nodeName, req.c))
		return
	}
	toNode := eg.findNode(scope, nodeName)
	if toNode == nil {
		req.fromNode.receiveCtrl(newLinkUpResponse(nil, stateNodeNotExist, scope, nodeName, req.c))
		return
	}

	link := newLink(eg, fromNode, toNode)
	if _, exist := eg.linkSet[link.name]; exist {
		// duplicated dlink, notify sender
		fromNode.receiveCtrl(newLinkUpResponse(nil, stateLinkDuplicated, scope, nodeName, req.c))
		return
	}
	if toNode.isExiting() {
		// the requested node wouldn't accept this dlink-up request
		fromNode.receiveCtrl(newLinkUpResponse(nil, stateLinkRefuse, scope, nodeName, req.c))
		return
	}
	eg.addLink(fromNode, link)
	eg.addLink(toNode, link)
	eg.linkSet[link.name] = true
	fromNode.receiveCtrl(newLinkUpResponse(link, stateSuccess, scope, nodeName, req.c))
}

// request breaking a dlink, such as A ----> B
// this request can come from A(user code) who initiates the operation
// when A don't want to send message to B anymore or
// come from B who wants to exit the event graph, and node delegate will
// silently break all links pointing to B before B really exited
func (eg *Graph) onLinkDown(req *linkDownRequest) {
	link := req.link
	fromNode := link.fromNode
	toNode := link.toNode
	// ensure every dlink can be torn down only once
	if _, exist := eg.linkSet[link.name]; !exist {
		fromNode.receiveCtrl(newLinkDownResponse(stateLinkNotExist, link))
		return
	}
	delete(eg.linkSet, link.name)
	eg.delLink(fromNode, link)
	eg.delLink(toNode, link)
	fromNode.receiveCtrl(newLinkDownResponse(stateSuccess, link))
}

// public APIs for end user

func NewEventGraph() *Graph {
	eg := &Graph{
		scopeMap:     make(scopeMapType),
		nodeMap:      make(nodeMapType),
		linkSet:      make(linkSetType),
		eventChannel: make(chan *Event),
	}

	// ensure loop started before return
	c := make(chan int)
	eg.startEventLoop(c)
	<-c
	return eg
}

// AddNode [SYNC] add a node to graph and wait until completion, i.e. the node's OnEnter is invoked
func (eg *Graph) AddNode(node Node) (success bool) {
	c := make(chan bool, 1)
	cb := func() { c <- true }
	evt := newNodeAddRequest(nodeAddRequest{node, cb})
	eg.deliveryEvent(evt)
	if err := utils.WaitChannelWithTimeout(c, 1, graphAddNodeTimeout); err == nil {
		success = true
	}
	return
}
