package comp

import (
	"fmt"
	"github.com/appcrash/media/server/utils"
	"reflect"
)

// trait is the means of class meta-info bookkeeping that props up runtime polymorphism
// it is a workaround because of limitation of golang runtime feature

// MessageTraitTag is a tag interface, if an interface is used only for extend message's behaviour, embed it, then
// use MessageTo() to quickly convert it
type MessageTraitTag interface{}

type Cloneable interface {
	MessageTraitTag
	Clone() Cloneable
}

// MessageTo convert message to specific trait object
func MessageTo[T MessageTraitTag](m Message) (v T) {
	if msg, ok := m.(T); ok {
		v = msg
	}
	return
}

const (
	maxMessageType = 512
	MessagePostfix = "Message"
)

var (
	// initialized by generated code when package loaded
	nbMessageTrait                int
	messageTraitRegistry          = make([]*MessageTrait, maxMessageType)
	messageNameQuery              = make(map[string]*MessageTrait)
	messageConvertibilityRegistry = make([]bool, maxMessageType*maxMessageType)
)

// MessageTrait is used to record all possible message kinds that flow among nodes of known types, ensure node links
// are compatible, that is Node A output is accepted by Node B input if there would be a link.
type MessageTrait struct {
	TypeId MessageType

	// For FooBarMessage struct
	// PtrType: *FooBarMessage
	// Type: FooBarMessage
	// ConvertType: a meta-type, the type of interface "FoobarMessageConvertable"
	PtrType, Type, ConvertType reflect.Type
}

func (m *MessageTrait) Clone() (cloned *MessageTrait) {
	cloned = new(MessageTrait)
	cloned.TypeId = m.TypeId
	cloned.PtrType = m.PtrType
	cloned.Type = m.Type
	cloned.ConvertType = m.ConvertType

	return
}

// ConvertFrom dynamically convert a message to message of this trait's type as long as the original message type
// implements the corresponding As***Message() method
func (m *MessageTrait) ConvertFrom(from Message) (to Message, err error) {
	if !CanConvertMessage(from.Type(), m.TypeId) {
		err = fmt.Errorf("can not convert message from type: %v to %v", from.Type(), m.TypeId)
		return
	}
	// following reflect actions don't check null ptr or any other errors, as the static analysis and start up
	// code should ensure the correctness. if panic do happens, ask user code author (it's all your faults!)

	method := reflect.ValueOf(from).
		Convert(m.ConvertType).Method(0) // get the method value of As***()
	returnValues := method.Call(nil) // Call As***() method to get the required message of this trait's type
	msgValue := returnValues[0]      // only one return value for method in interface ***MessageConvertable
	msg := msgValue.Interface()
	if msg == nil {
		logger.Warnf("get nil when converting message from type: %v to %v", from.Type(), m.TypeId)
		// don't transform nil value, return it directly
		return
	}
	to = msg.(Message)
	return
}

func (m *MessageTrait) String() string {
	return fmt.Sprintf("[message trait: id:%v type:%v", m.TypeId, m.Type.Name())
}

func (m *MessageTrait) Match(peer *MessageTrait) bool {
	return m.TypeId == peer.TypeId
}

// Name is a string that used in nmd language or node's struct tag to identify this trait
func (m *MessageTrait) Name() string {
	typeName := m.Type.Name()
	baseName := typeName[:len(typeName)-len(MessagePostfix)]
	return utils.CamelCaseToSnake(baseName)
}

func MT[T any](convertType reflect.Type) *MessageTrait {
	ptrType := reflect.TypeOf(new(T))
	structType := ptrType.Elem()

	msg, ok := reflect.New(structType).Interface().(Message)
	if !ok {
		panic(fmt.Errorf("%v doesn't implement Message", ptrType.String()))
	}
	trait := &MessageTrait{
		TypeId:      msg.Type(),
		PtrType:     ptrType,
		Type:        structType,
		ConvertType: convertType,
	}

	return trait
}

func MessageTraitOfObject(model Message) (mt *MessageTrait, exist bool) {
	return MessageTraitOfType(model.Type())
}

func MessageTraitOfType(typeId MessageType) (mt *MessageTrait, exist bool) {
	i := int(typeId)
	if i >= maxMessageType {
		return
	}
	trait := messageTraitRegistry[i]
	if trait == nil {
		return
	}
	mt = trait.Clone()
	exist = true
	return
}

func MessageTraitOfName(name string) (mt *MessageTrait, exist bool) {
	if trait, ok := messageNameQuery[name]; ok {
		mt = trait.Clone()
		exist = true
	}
	return
}

func AddMessageTrait(traits ...*MessageTrait) {
	for _, t := range traits {
		if int(t.TypeId) != nbMessageTrait {
			panic(fmt.Errorf("add message trait of type:%v with wrong type id: %v which should be: %v",
				t.Type.Name(), t.TypeId, nbMessageTrait))

		}
		if _, exist := MessageTraitOfName(t.Name()); exist {
			panic(fmt.Errorf("add message trait of type:%v with duplicated trait name", t.Type.Name()))
		}
		logger.Infof("[MESSAGE TRAIT]: name:%v, type_id:%v, type:%v",
			t.Name(), t.TypeId, t.Type.String())
		messageNameQuery[t.Name()] = t
		messageTraitRegistry[nbMessageTrait] = t
		nbMessageTrait++
	}
}

func SetMessageConvertable(from, to MessageType) {
	if from > maxMessageType || to > maxMessageType {
		panic("SetMessageConvertable failed due to message type too large")
	}
	messageConvertibilityRegistry[from*maxMessageType+to] = true
}

func CanConvertMessage(from, to MessageType) bool {
	if from > maxMessageType || to > maxMessageType {
		return false
	}
	return messageConvertibilityRegistry[from*maxMessageType+to]
}
