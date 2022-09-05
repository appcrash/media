package comp

import (
	"fmt"
	"github.com/appcrash/media/server/utils"
	"reflect"
)

//go:generate go run ../../cmd/gentrait -t message -o trait_message_generated.go

// trait is the means of class meta-info bookkeeping that props up runtime polymorphism
// it is a workaround because of limitation of golang runtime feature

type Cloneable interface {
	Clone() Cloneable
}

var (
	messageMetaType   = MetaType[Message]()
	cloneableMetaType = MetaType[Cloneable]()
)

const (
	messageTraitCloneable = 1
)

// MessageTrait is used to record all possible message kinds that flow among nodes of known types, ensure node links
// are compatible, that is Node A output is accepted by Node B input if there would be a link.
type MessageTrait struct {
	utils.Flag[uint32]
	TypeId        MessageType
	PtrType, Type reflect.Type
}

func (m *MessageTrait) Clone() (cloned *MessageTrait) {
	cloned = new(MessageTrait)
	cloned.Flag = m.Flag
	cloned.TypeId = m.TypeId
	cloned.PtrType = m.PtrType
	cloned.Type = m.Type
	return
}

func (m *MessageTrait) String() string {
	return fmt.Sprintf("[message trait: id:%v type:%v", m.TypeId, m.Type.Name())
}

func (m *MessageTrait) Match(peer *MessageTrait) bool {
	return m.TypeId == peer.TypeId
}

func (m *MessageTrait) IsCloneable() bool {
	return m.HasFlag(messageTraitCloneable)
}

var messageTraitRegistry = make(map[MessageType]*MessageTrait)

func MT[T any]() *MessageTrait {
	ptrType := reflect.TypeOf(new(T))
	structType := ptrType.Elem()
	if !ptrType.Implements(messageMetaType) {
		panic(fmt.Errorf("type %v doesn't implements message interface", ptrType.String()))
	}

	msg := reflect.New(structType).Interface().(Message)
	trait := &MessageTrait{
		TypeId:  msg.Type(),
		PtrType: ptrType,
		Type:    structType,
	}

	// inspect interface trait
	if ptrType.Implements(cloneableMetaType) {
		trait.SetFlag(messageTraitCloneable)
	}
	return trait
}

func RegisterMessageTrait(traits ...*MessageTrait) {
	for _, trait := range traits {
		typ := trait.Type
		typeId := trait.TypeId
		if anotherTrait, exist := messageTraitRegistry[typeId]; exist {
			panic(fmt.Sprintf("message(%v) of type %v already registered by %v", typ.String(), typeId, anotherTrait.Type.String()))
		}
		logger.Infof("register message anotherTrait (id:%v) => (type:%v)", typeId, typ.String())
		messageTraitRegistry[typeId] = trait
	}
}

func MessageTraitOfObject(model Message) (mt *MessageTrait, exist bool) {
	return MessageTraitOfType(model.Type())
}

func MessageTraitOfType(typeId MessageType) (mt *MessageTrait, exist bool) {
	var pmt *MessageTrait
	if pmt, exist = messageTraitRegistry[typeId]; exist {
		mt = pmt.Clone()
	}
	return
}
