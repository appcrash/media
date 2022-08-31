package comp

import (
	"hash/fnv"
)

// LinkPad is the default LinkPoint impl
type LinkPad struct {
	owner        SessionAware
	peer         SessionAware
	linkId       int
	identity     uint64
	messageTrait *MessageTrait
	sendFunc     func(msg Message) error
}

func (l *LinkPad) LinkId() int {
	return l.linkId
}

func (l *LinkPad) Identity() uint64 {
	return l.identity
}

func (l *LinkPad) Owner() SessionAware {
	return l.owner
}

func (l *LinkPad) Peer() SessionAware {
	return l.peer
}

func (l *LinkPad) SetPeer(s SessionAware) {
	l.peer = s
}

func (l *LinkPad) SendMessage(msg Message) error {
	return l.sendFunc(msg)
}

func (l *LinkPad) MessageTrait() *MessageTrait {
	return l.messageTrait
}

func MakeLinkIdentity(session, name string, linkId int) uint64 {
	h := fnv.New32a()
	h.Write([]byte(session))
	h.Write([]byte(name))
	msb := uint64(h.Sum32())
	return (msb << 32) | uint64(linkId)
}
