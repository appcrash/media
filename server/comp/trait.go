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

var cloneableType = MetaType[Cloneable]()
var sessionAwareType = MetaType[SessionAware]()

// MessageTrait is used to record all possible message kinds that flow among nodes of known types, ensure node links
// are compatible, that is Node A output is accepted by Node B input if there would be a link.
type MessageTrait struct {
	TypeId MessageType
	Type   reflect.Type
}

func (m *MessageTrait) Match(peer *MessageTrait) bool {
	return m.TypeId == peer.TypeId ||
		m.TypeId == ANY ||
		peer.TypeId == ANY
}

func (m *MessageTrait) IsCloneable() bool {
	return m.Type.Implements(cloneableType)
}

var messageTraitRegistry = make(map[MessageType]*MessageTrait)

func RegisterMessageTrait(models ...Message) {
	for _, model := range models {
		typ := reflect.TypeOf(model)
		typeId := model.Type()
		mt := &MessageTrait{
			TypeId: typeId,
			Type:   typ,
		}
		if trait, exist := messageTraitRegistry[typeId]; exist {
			panic(fmt.Sprintf("message(%v) of type %v already registered by %v", typ.String(), typeId, trait.Type.String()))
		}
		logger.Infof("register message trait (id:%v) => (type:%v)", typeId, typ.String())
		messageTraitRegistry[typeId] = mt
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
	// user filled fields
	Name    string
	NewFunc func() SessionAware // new function only alloc node, not initialize it
	Accept  []MessageType
	Offer   []MessageType

	// generated when registering
	Type reflect.Type
}

var nodeTraitRegistry = make(map[string]*NodeTrait)

func getNodeByName(typeName string) SessionAware {
	if trait, ok := nodeTraitRegistry[typeName]; ok {
		return trait.NewFunc()
	}
	logger.Errorf("unknown node type:%v\n", typeName)
	return nil
}

// NT is a template method to make node trait
func NT[T any]() *NodeTrait {
	ptrType := reflect.TypeOf(new(T))
	structType := ptrType.Elem()
	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("type %v should be struct kind", structType.String()))
	}
	if !ptrType.Implements(sessionAwareType) {
		panic(fmt.Errorf("type %v doesn't implements session aware", ptrType.String()))
	}
	name := utils.CamelCaseToSnake(structType.Name())
	newFunc := func() SessionAware {
		node := reflect.New(structType).Interface()
		return node.(SessionAware)
	}
	obj := newFunc()
	return &NodeTrait{
		Name:    name,
		NewFunc: newFunc,
		Accept:  obj.Accept(),
		Offer:   obj.Offer(),
		Type:    structType,
	}
}

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
		for _, messageType := range trait.Accept {
			if _, ok := MessageTraitOfType(messageType); !ok {
				return fmt.Errorf("node trait(%v) accept unknown message type %v", name, messageType)
			}
		}
		for _, messageType := range trait.Offer {
			if _, ok := MessageTraitOfType(messageType); !ok {
				return fmt.Errorf("node trait(%v) offer unknown message type %v", name, messageType)
			}
		}

		logger.Infof("register node trait: name: %v, type: %v", name, trait.Type.String())
		nodeTraitRegistry[name] = trait
	}
	return nil
}

func NodeTraitOfName(name string) (nt NodeTrait, exist bool) {
	var pnt *NodeTrait
	if pnt, exist = nodeTraitRegistry[name]; exist {
		nt = *pnt
	}
	return
}
