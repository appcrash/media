package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/utils"
	"reflect"
)

type Channelable[T any] interface {
	ChannelLink(c chan T)
}

var (
	sessionAwareMetaType = MetaType[SessionAware]()
	channelableMetaType  = MetaType[Channelable[[]byte]]() // this can be more generic if golang improves
)

const (
	nodeTraitChannelable = 1
)

// NodeTrait record factory method, negotiation infos
type NodeTrait struct {
	utils.Flag[uint32]
	Name          string
	NewFunc       func() SessionAware // new function only alloc node, not initialize it
	Accept        []*MessageTrait
	Offer         []*MessageTrait
	PtrType, Type reflect.Type
}

func (nt *NodeTrait) Clone() (cloned *NodeTrait) {
	cloned = new(NodeTrait)
	cloned.Flag = nt.Flag
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

func (nt *NodeTrait) IsChannelable() bool {
	return nt.HasFlag(nodeTraitChannelable)
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
			accept = append(accept, tr)
		}
	}
	for _, messageType := range nodeObj.Offer() {
		if tr, ok := MessageTraitOfType(messageType); !ok {
			panic(fmt.Errorf("node type(%v) offer unknown message type %v", name, messageType))
		} else {
			offer = append(offer, tr)
		}
	}
	trait.Name = name
	trait.NewFunc = newFunc
	trait.Accept = accept
	trait.Offer = offer
	trait.PtrType = ptrType
	trait.Type = structType

	// inspect interface trait
	if ptrType.Implements(channelableMetaType) {
		trait.SetFlag(nodeTraitChannelable)
	}
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
