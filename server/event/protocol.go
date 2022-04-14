package event

// these commands(request/response) only used for graph's internal communication
const (
	reqLinkUp = iota
	reqLinkDown
	reqNodeAdd
	reqNodeExit
)

const (
	respLinkUp = iota + 10000
	respLinkDown
	respNodeAdd
	respNodeExit
)

const stateSuccess = 0
const stateLinkNotExist = stateSuccess - 1
const stateLinkRefuse = stateLinkNotExist - 1
const stateLinkDuplicated = stateLinkRefuse - 1
const stateNodeNotExist = stateLinkDuplicated - 1
const stateNodeExceedMaxLink = stateNodeNotExist - 1

/* ------- request structs ------- */
type linkUpRequest struct {
	fromNode *NodeDelegate
	scope    string
	nodeName string
	c        chan int
}

type linkDownRequest struct {
	link *dlink
}

type nodeAddRequest struct {
	node Node
	cb   Callback
}

type nodeExitRequest struct {
	delegate *NodeDelegate
}

/* ------- response structs ------- */
type linkUpResponse struct {
	state    int
	link     *dlink
	scope    string
	nodeName string
	c        chan int
}

type linkDownResponse struct {
	state int
	link  *dlink
}

type nodeAddResponse struct {
	delegate *NodeDelegate
	cb       Callback
}

type nodeExitResponse struct {
}

// event creators
// unluckily, golang doesn't support macro or meta-programming, we have to
// craft each factory method by hand :(

func NewEventWithCallback(cmd int, obj interface{}, cb Callback) *Event {
	return &Event{cmd: cmd, obj: obj, cb: cb}
}

func NewEvent(cmd int, obj interface{}) *Event {
	return NewEventWithCallback(cmd, obj, nil)
}

/* ---------------REQUEST------------------- */
func newLinkUpRequest(nd *NodeDelegate, scope string, nodeName string, c chan int) *Event {
	return NewEvent(reqLinkUp, &linkUpRequest{nd, scope, nodeName, c})
}

func newLinkDownRequest(link *dlink) *Event {
	return NewEvent(reqLinkDown, &linkDownRequest{link})
}

func newNodeAddRequest(req nodeAddRequest) *Event {
	return NewEvent(reqNodeAdd, &req)
}

func newNodeExitRequest(node *NodeDelegate) *Event {
	return NewEvent(reqNodeExit, &nodeExitRequest{node})
}

/* ---------------RESPONSE------------------- */
func newLinkUpResponse(resp *dlink, state int, scope string, name string, c chan int) *Event {
	return NewEvent(respLinkUp, &linkUpResponse{state, resp, scope, name, c})
}

func newLinkDownResponse(state int, link *dlink) *Event {
	return NewEvent(respLinkDown, &linkDownResponse{state, link})
}

func newNodeAddResponse(delegate *NodeDelegate, cb Callback) *Event {
	return NewEvent(respNodeAdd, &nodeAddResponse{delegate, cb})
}

func newNodeExitResponse() *Event {
	return NewEvent(respNodeExit, &nodeExitResponse{})
}
