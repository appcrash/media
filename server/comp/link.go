package comp

import (
	"fmt"
	"github.com/appcrash/media/server/utils"
	"hash/fnv"
	"reflect"
	"strings"
	"sync/atomic"
)

type sendFuncType func(msg Message) error

// LinkPad is the default LinkPoint impl
type LinkPad struct {
	owner        SessionAware
	linkId       int
	identity     LinkIdentityType
	enabled      atomic.Value
	messageTrait *MessageTrait
	sendFunc     sendFuncType
}

func (l *LinkPad) LinkId() int {
	return l.linkId
}

func (l *LinkPad) Identity() LinkIdentityType {
	return l.identity
}

func (l *LinkPad) Owner() SessionAware {
	return l.owner
}

func (l *LinkPad) SendMessage(msg Message) (err error) {
	if !l.enabled.Load().(bool) {
		return nil
	}
	if err = l.sendFunc(msg); err != nil {
		//logger.Debugf("disable linkpoint %v of %v as send message failed", l.identity, l.owner)
		//l.SetEnabled(false)
	}
	return
}

func (l *LinkPad) SetEnabled(e bool) {
	l.enabled.Store(e)
}

func (l *LinkPad) MessageTrait() *MessageTrait {
	return l.messageTrait
}

func NewLinkPad(owner SessionAware, linkId int, identity LinkIdentityType, messageTrait *MessageTrait, sendFunc sendFuncType) *LinkPad {
	pad := &LinkPad{
		owner:        owner,
		linkId:       linkId,
		identity:     identity,
		messageTrait: messageTrait,
		sendFunc:     sendFunc,
	}
	pad.enabled.Store(true)
	return pad
}

func MakeLinkIdentity(session, name string, linkId int) LinkIdentityType {
	h := fnv.New32a()
	h.Write([]byte(session))
	h.Write([]byte(name))
	msb := uint64(h.Sum32())
	return LinkIdentityType((msb << 32) | uint64(linkId))
}

// tag is of "key1=value1,key_only,...",supported keys:
//
// type={msg_snake_name}
// nullable
func injectLinkPoint(structField reflect.StructField, field reflect.Value, lps []LinkPoint, tag string) error {
	var nullable bool
	var msgType string
	for _, prop := range strings.Split(tag, ",") {
		if prop == "nullable" {
			nullable = true
		} else if strings.HasPrefix(prop, "type=") {
			msgType = prop[5:]
		}
	}
	if !nullable && len(msgType) == 0 {
		return fmt.Errorf("field %v is not nullable but no message type set", structField)
	}
	for _, lp := range lps {
		if lp.MessageTrait().Name() == msgType {
			// only support injecting one link point for each message type
			utils.SetField(field, reflect.ValueOf(lp))
			return nil
		}
	}
	if nullable {
		return nil
	} else {
		return fmt.Errorf("field %v can not be inject as link point of message type %v not found", structField, msgType)
	}
}
