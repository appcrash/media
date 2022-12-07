package comp

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	sessionAwareMetaType = MetaType[SessionAware]()
	nodeTraitRegistry    = make(map[string]*NodeTrait)
)

// NodeTraitTag is a tag interface, if an interface is used only for extend node's behaviour, embed it, then
// use NodeTo() to quickly convert it
type NodeTraitTag interface{}

type Channelable[T any] interface {
	NodeTraitTag
	ChannelLink(c chan T)
}

type PreComposer interface {
	NodeTraitTag
	// BeforeCompose is called after nodes initialized but not added to graph yet
	BeforeCompose(c *Composer, node SessionAware) error
}

type PostComposer interface {
	NodeTraitTag
	// AfterCompose is called after nodes are negotiated and connected
	AfterCompose(c *Composer, node SessionAware) error
}

type InitializingNode interface {
	NodeTraitTag
	// Init do initialization after node is allocated and configured
	Init() error
}

type UnInitializingNode interface {
	NodeTraitTag
	// UnInit is called when initialization failed or session terminated, before node exit graph
	UnInit()
}

// NodeTo convert session node to node with specific trait object
func NodeTo[T NodeTraitTag](n SessionAware) (v T) {
	if node, ok := n.(T); ok {
		v = node
	}
	return
}

// NodeTrait record static information of node such as factory method, negotiation acceptance
type NodeTrait struct {
	NodeType      string
	FactoryFunc   func() SessionAware
	PtrType, Type reflect.Type
}

func (nt *NodeTrait) Clone() (cloned *NodeTrait) {
	cloned = new(NodeTrait)
	cloned.NodeType = nt.NodeType
	cloned.FactoryFunc = nt.FactoryFunc
	cloned.PtrType = nt.PtrType
	cloned.Type = nt.Type
	return
}

// NT is a template method to make node trait
func NT[T any](typeName string, factoryFunc func() SessionAware) *NodeTrait {
	ptrType := reflect.TypeOf(new(T))
	structType := ptrType.Elem()
	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("node type %v should be struct kind", structType.String()))
	}
	if !ptrType.Implements(sessionAwareMetaType) {
		panic(fmt.Errorf("node type %v doesn't implements session aware", ptrType.String()))
	}

	//name := utils.CamelCaseToSnake(structType.Name())
	trait := &NodeTrait{}
	//newFunc := func() SessionAware {
	//	node := reflect.New(structType).Interface()
	//	// this is SessionNode specific hack, other session aware implementation must define these fields to
	//	// proceed successfully
	//	nodeValue := reflect.ValueOf(node).Elem()
	//	nodeValue.FieldByName("Trait").Set(reflect.ValueOf(trait.Clone()))
	//	return node.(SessionAware)
	//}
	//nodeObj := reflect.New(structType).Interface().(SessionAware)

	trait.NodeType = typeName
	trait.FactoryFunc = factoryFunc
	trait.PtrType = ptrType
	trait.Type = structType

	return trait
}

// RegisterNodeTrait enable node being configed by name required by nmd graph
// the name is simply got by converting node struct name to snake case, i.e. ChanSink -> chan_sink
func RegisterNodeTrait(traits ...*NodeTrait) error {
	// sanity check at start-up time to avoid runtime checking
	for _, trait := range traits {
		name := trait.NodeType
		if name == "" {
			return errors.New("empty trait name")
		}
		if _, ok := nodeTraitRegistry[name]; ok {
			return fmt.Errorf("node traiit(%v) already registered", name)
		}
		if trait.FactoryFunc == nil {
			return fmt.Errorf("node trait(%v) has no factory method", name)
		}

		logger.Infof("[NODE TRAIT]: name: %v, type: %v", name, trait.Type.String())
		nodeTraitRegistry[name] = trait
	}
	return nil
}

func NodeTraitOfType(name string) (nt *NodeTrait, exist bool) {
	var pnt *NodeTrait
	if pnt, exist = nodeTraitRegistry[name]; exist {
		nt = pnt.Clone()
	}
	return
}
