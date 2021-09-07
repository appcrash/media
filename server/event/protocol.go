package event

// these commands(request/response) only used for graph's internal communication
const (
	req_link_up = iota
	req_link_down
	req_node_add
	req_node_exit
)

const (
	resp_link_up = iota + 10000
	resp_link_down
	resp_node_add
	resp_node_exit
)

const state_success = 0
const state_link_not_exist = state_success - 1
const state_link_refuse = state_link_not_exist - 1
const state_link_duplicated = state_link_refuse - 1
const state_node_not_exist = state_link_duplicated - 1
const state_node_exceed_max_link = state_node_not_exist - 1

/* ------- request structs ------- */
type linkUpRequest struct {
	fromNode *NodeDelegate
	scope    string
	nodeName string
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
func newLinkUpRequest(nd *NodeDelegate, scope string, nodeName string) *Event {
	return NewEvent(req_link_up, &linkUpRequest{nd, scope, nodeName})
}

func newLinkDownRequest(link *dlink) *Event {
	return NewEvent(req_link_down, &linkDownRequest{link})
}

func newNodeAddRequest(req nodeAddRequest) *Event {
	return NewEvent(req_node_add, &req)
}

func newNodeExitRequest(node *NodeDelegate) *Event {
	return NewEvent(req_node_exit, &nodeExitRequest{node})
}

/* ---------------RESPONSE------------------- */
func newLinkUpResponse(resp *dlink, state int, scope string, name string) *Event {
	return NewEvent(resp_link_up, &linkUpResponse{state, resp, scope, name})
}

func newLinkDownResponse(state int, link *dlink) *Event {
	return NewEvent(resp_link_down, &linkDownResponse{state, link})
}

func newNodeAddResponse(delegate *NodeDelegate, cb Callback) *Event {
	return NewEvent(resp_node_add, &nodeAddResponse{delegate, cb})
}

func newNodeExitResponse() *Event {
	return NewEvent(resp_node_exit, &nodeExitResponse{})
}
