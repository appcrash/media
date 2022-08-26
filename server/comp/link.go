package comp

// LinkPad is the default LinkPoint impl
type LinkPad struct {
	owner        SessionAware
	id           int
	messageTrait MessageTrait
	peer         LinkPoint
	sendFunc     func(msg Message) error
}

func (l *LinkPad) LinkId() int {
	return l.id
}

func (l *LinkPad) Owner() SessionAware {
	return l.owner
}

func (l *LinkPad) SetPeer(lp LinkPoint) {
	l.peer = lp
}

func (l *LinkPad) Peer() LinkPoint {
	return l.peer
}

func (l *LinkPad) SendMessage(msg Message) error {
	return l.sendFunc(msg)
}

func (l *LinkPad) MessageTrait() MessageTrait {
	return l.messageTrait
}
