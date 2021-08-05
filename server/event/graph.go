package event

import (
	"reflect"
)

type scopeMapType map[string][]*NodeDelegate
type linkSetType map[string]bool
type nodeMapType map[string]*nodeInfo

type EventGraph struct {
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

func (eg *EventGraph) findNode(scope string, name string) *NodeDelegate {
	if nodeList, ok := eg.scopeMap[scope]; ok {
		for _, node := range nodeList {
			if node.getNodeName() == name {
				return node
			}
		}
	}
	return nil
}

func (eg *EventGraph) addNode(nd *NodeDelegate, maxLink int) {
	scope := nd.getNodeScope()
	if nodeList, ok := eg.scopeMap[scope]; !ok {
		eg.scopeMap[scope] = []*NodeDelegate{nd}
	} else {
		eg.scopeMap[scope] = append(nodeList, nd)
	}
	eg.nodeMap[nd.getId()] = &nodeInfo{maxLink: maxLink}
}

func (eg *EventGraph) delNode(nd *NodeDelegate) {
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
}

func (eg *EventGraph) getNodeInfo(nodeId string) *nodeInfo {
	if info, exist := eg.nodeMap[nodeId]; exist {
		return info
	}
	return nil
}

// associate dlink to a node delegate
func (eg *EventGraph) addLink(nd *NodeDelegate, l *dlink) {
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
		// TODO: log error
	}
}

// tear down a dlink in node delegate
func (eg *EventGraph) delLink(nd *NodeDelegate, l *dlink) {
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
		// TODO: log error
		return
	}
	// remove dlink by exchanging with last one
	lastLink := (*linkArray)[length-1]
	(*linkArray)[index] = lastLink
	*linkArray = (*linkArray)[:length-1]
}

func (eg *EventGraph) deliveryEvent(evt *Event) {
	eg.eventChannel <- evt
}

// simply loop forever
func (eg *EventGraph) startEventLoop(c chan int) {
	go func(g *EventGraph) {
		c <- 0
		for {
			evt := <-g.eventChannel
			g.onEvent(evt)
		}
	}(eg)
}

func (eg *EventGraph) onEvent(evt *Event) {
	var ok bool
	switch evt.cmd {
	case req_node_add:
		var req *nodeAddRequest
		if req, ok = evt.obj.(*nodeAddRequest); !ok {
			return
		}
		eg.onAddNode(req)
	case req_node_exit:
		var req *nodeExitRequest
		if req, ok = evt.obj.(*nodeExitRequest); !ok {
			return
		}
		eg.onExitNode(req)
	case req_link_up:
		var req *linkUpRequest
		if req, ok = evt.obj.(*linkUpRequest); !ok {
			return
		}
		eg.onLinkUp(req)
	case req_link_down:
		var req *linkDownRequest
		if req, ok = evt.obj.(*linkDownRequest); !ok {
			return
		}
		eg.onLinkDown(req)
	}
}

// add node to graph, and send node-add response to this node immediately
func (eg *EventGraph) onAddNode(req *nodeAddRequest) {
	maxLink := defaultMaxLink
	node := req.node
	ps := reflect.ValueOf(node)

	elem := ps.Elem()
	if elem.Kind() == reflect.Struct {
		field := elem.FieldByName("maxLink")
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
	eg.addNode(delegate,maxLink)
	// all gears up, rock it
	go func(nd *NodeDelegate) {
		nd.startEventLoop()
		resp := newNodeAddResponse(nd)
		nd.receiveCtrl(resp)
	}(delegate)

}

// node requests exiting the graph, notify all senders linking
// to this node
func (eg *EventGraph) onExitNode(req *nodeExitRequest) {
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
			link.fromNode.receiveCtrl(newLinkDownResponse(state_success, link))
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
		// TODO: log error
		panic("node still have active links")
	}
	eg.delNode(nd)
	// finally, send the last ctrl message for this node
	nd.receiveCtrl(newNodeExitResponse())
}

// request dlink to other node, decline if that node is exiting or
// dlink is duplicated, otherwise create dlink between them
func (eg *EventGraph) onLinkUp(req *linkUpRequest) {
	nodeName := req.nodeName
	scope := req.scope
	fromNode := req.fromNode

	ni :=  eg.nodeMap[fromNode.getId()]
	if ni.maxLink == len(ni.outputLinks) {
		req.fromNode.receiveCtrl(newLinkUpResponse(nil,state_node_exceed_max_link,scope,nodeName))
		return
	}
	toNode := eg.findNode(scope, nodeName)
	if toNode == nil {
		req.fromNode.receiveCtrl(newLinkUpResponse(nil, state_node_not_exist, scope, nodeName))
		return
	}

	link := newLink(eg, fromNode, toNode)
	if _, exist := eg.linkSet[link.name]; exist {
		// duplicated dlink, notify sender
		fromNode.receiveCtrl(newLinkUpResponse(nil, state_link_duplicated, scope, nodeName))
		return
	}
	if toNode.isExiting() {
		// the requested node wouldn't accept this dlink-up request
		fromNode.receiveCtrl(newLinkUpResponse(nil, state_link_refuse, scope, nodeName))
		return
	}
	eg.addLink(fromNode, link)
	eg.addLink(toNode, link)
	eg.linkSet[link.name] = true
	fromNode.receiveCtrl(newLinkUpResponse(link, state_success, scope, nodeName))
}

// request breaking a dlink, such as A ----> B
// this request can come from A(user code) who initiates the operation
// when A don't want to send message to B anymore or
// come from B who wants to exit the event graph, and node delegate will
// silently break all links pointing to B before B really exited
func (eg *EventGraph) onLinkDown(req *linkDownRequest) {
	link := req.link
	fromNode := link.fromNode
	toNode := link.toNode
	// ensure every dlink can be torn down only once
	if _, exist := eg.linkSet[link.name]; !exist {
		fromNode.receiveCtrl(newLinkDownResponse(state_link_not_exist, link))
		return
	}
	delete(eg.linkSet, link.name)
	eg.delLink(fromNode, link)
	eg.delLink(toNode, link)
	fromNode.receiveCtrl(newLinkDownResponse(state_success, link))
}

// public APIs for end user

func NewEventGraph() *EventGraph {
	eg := &EventGraph{
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

func (eg *EventGraph) AddNode(node Node) {
	evt := newNodeAddRequest(nodeAddRequest{node})
	eg.deliveryEvent(evt)
}
