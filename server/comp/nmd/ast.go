package nmd

import (
	"fmt"
	"github.com/appcrash/media/server/utils"
)

// TODO: every statement should have sequence id

type NodeProp struct {
	Key, Type string
	Value     interface{}
}

func (np *NodeProp) FormalizeKey() {
	// normal keys with camel case remain intact
	np.Key = utils.SnakeToCamelCase(np.Key)
}

type LinkOperator struct {
	LinkTo      *NodeDef // which node link to
	PreferOffer []string // preferred offer connecting node suggests
}

type NodeDef struct {
	Index             int
	Name, Scope, Type string
	Props             []*NodeProp
	Deps              []*LinkOperator // record receivers of this node
}

func (nd *NodeDef) String() string {
	var prop string
	for _, p := range nd.Props {
		prop += fmt.Sprintf("%v=%v ", p.Key, p.Value)
	}
	return fmt.Sprintf("[%v@%v:%v  prop: %v] ", nd.Name, nd.Scope, nd.Type, prop)
}

type EndpointDefs struct {
	Nodes       []*NodeDef
	PreferOffer []string
}

type CallActionDefs struct {
	Node *NodeDef
	Cmd  string
}

type CastActionDefs struct {
	CallActionDefs
}

type SinkActionDefs struct {
	NodeName string
}
