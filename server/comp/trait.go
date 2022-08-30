package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/utils"
	"reflect"
)

// trait is the means of class meta-info bookkeeping that props up runtime polymorphism

func MetaType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

var (
	cloneableMetaType    = MetaType[Cloneable]()
	messageMetaType      = MetaType[Message]()
	sessionAwareMetaType = MetaType[SessionAware]()
)

// MessageTrait is used to record all possible message kinds that flow among nodes of known types, ensure node links
// are compatible, that is Node A output is accepted by Node B input if there would be a link.
type MessageTrait struct {
	TypeId        MessageType
	PtrType, Type reflect.Type
}

func (m *MessageTrait) Clone() (cloned *MessageTrait) {
	cloned = new(MessageTrait)
	cloned.TypeId = m.TypeId
	cloned.PtrType = m.PtrType
	cloned.Type = m.Type
	return
}

func (m *MessageTrait) Match(peer *MessageTrait) bool {
	return m.TypeId == peer.TypeId ||
		m.TypeId == ANY ||
		peer.TypeId == ANY
}

func (m *MessageTrait) IsCloneable() bool {
	return m.PtrType.Implements(cloneableMetaType)
}

var messageTraitRegistry = make(map[MessageType]*MessageTrait)

func MT[T any]() *MessageTrait {
	ptrType := reflect.TypeOf(new(T))
	structType := ptrType.Elem()
	if !ptrType.Implements(messageMetaType) {
		panic(fmt.Errorf("type %v doesn't implements message interface", ptrType.String()))
	}
	msg := reflect.New(structType).Interface().(Message)
	return &MessageTrait{
		TypeId:  msg.Type(),
		PtrType: ptrType,
		Type:    structType,
	}
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

func MessageTraitOfObject(model Message) (mt MessageTrait, exist bool) {
	return MessageTraitOfType(model.Type())
}

func MessageTraitOfType(typeId MessageType) (mt MessageTrait, exist bool) {
	var pmt *MessageTrait
	if pmt, exist = messageTraitRegistry[typeId]; exist {
		mt = *pmt
	}
	return
}

// NodeTrait record factory method, negotiation infos
type NodeTrait struct {
	Name          string
	NewFunc       func() SessionAware // new function only alloc node, not initialize it
	Accept        []*MessageTrait
	Offer         []*MessageTrait
	PtrType, Type reflect.Type
}

func (nt *NodeTrait) Clone() (cloned *NodeTrait) {
	cloned = new(NodeTrait)
	cloned.Name = nt.Name
	cloned.NewFunc = nt.NewFunc
	cloned.PtrType = nt.PtrType
	cloned.Type = nt.Type
	for _, mt := range nt.Accept {
		cloned.Accept = append(cloned.Accept, mt.Clone())
	}
	for _, mt := range nt.Offer {
		cloned.Offer = append(cloned.Offer, mt.Clone())
	}
	return
}

var nodeTraitRegistry = make(map[string]*NodeTrait)

// NT is a template method to make node trait
func NT[T any]() *NodeTrait {
	ptrType := reflect.TypeOf(new(T))
	structType := ptrType.Elem()
	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("node type %v should be struct kind", structType.String()))
	}
	if !ptrType.Implements(sessionAwareMetaType) {
		panic(fmt.Errorf("node type %v doesn't implements session aware", ptrType.String()))
	}

	name := utils.CamelCaseToSnake(structType.Name())
	var accept, offer []*MessageTrait
	trait := &NodeTrait{}
	newFunc := func() SessionAware {
		node := reflect.New(structType).Interface()
		// this is SessionNode specific hack, other session aware implementation must define these fields to
		// proceed successfully
		nodeValue := reflect.ValueOf(node).Elem()
		nodeValue.FieldByName("Name").Set(reflect.ValueOf(name))
		nodeValue.FieldByName("Trait").Set(reflect.ValueOf(trait.Clone()))
		return node.(SessionAware)
	}
	nodeObj := newFunc()

	for _, messageType := range nodeObj.Accept() {
		if tr, ok := MessageTraitOfType(messageType); !ok {
			panic(fmt.Errorf("node type(%v) accept unknown message type %v", name, messageType))
		} else {
			accept = append(accept, &tr)
		}
	}
	for _, messageType := range nodeObj.Offer() {
		if tr, ok := MessageTraitOfType(messageType); !ok {
			panic(fmt.Errorf("node type(%v) offer unknown message type %v", name, messageType))
		} else {
			offer = append(offer, &tr)
		}
	}
	trait.Name = name
	trait.NewFunc = newFunc
	trait.Accept = accept
	trait.Offer = offer
	trait.PtrType = ptrType
	trait.Type = structType
	return trait
}

// RegisterNodeTrait enable node can be configed by name required by nmd graph
// the name is simply got by converting node struct name to snake case, i.e. ChanSink -> chan_sink
func RegisterNodeTrait(traits ...*NodeTrait) error {
	// sanity check at start-up time to avoid runtime checking
	for _, trait := range traits {
		name := trait.Name
		if name == "" {
			return errors.New("empty trait name")
		}
		if _, ok := nodeTraitRegistry[name]; ok {
			return fmt.Errorf("node traiit(%v) already registered", name)
		}
		if trait.NewFunc == nil {
			return fmt.Errorf("node trait(%v) has no factory method", name)
		}

		logger.Infof("register node trait: name: %v, type: %v", name, trait.Type.String())
		nodeTraitRegistry[name] = trait
	}
	return nil
}

func NodeTraitOfName(name string) (nt *NodeTrait, exist bool) {
	var pnt *NodeTrait
	if pnt, exist = nodeTraitRegistry[name]; exist {
		nt = pnt.Clone()
	}
	return
}
