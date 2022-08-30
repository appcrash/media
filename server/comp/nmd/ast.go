package nmd

import (
	"fmt"
	"github.com/appcrash/media/server/utils"
	"regexp"
)

// TODO: every statement should have sequence id

type NodeProp struct {
	Key, Type string
	Value     interface{}
}

const regCamelCasePattern = `_+[a-z]`

var regCamelCase = regexp.MustCompile(regCamelCasePattern)

func (np *NodeProp) FormalizeKey() {
	// normal keys with camel case remain intact
	np.Key = utils.SnakeToCamelCase(np.Key)
}

type NodeDef struct {
	Index             int
	Name, Scope, Type string
	Props             []*NodeProp
	Deps              []*NodeDef // record receivers of this node
}

func (nd *NodeDef) String() string {
	var prop string
	for _, p := range nd.Props {
		prop += fmt.Sprintf("%v=%v ", p.Key, p.Value)
	}
	return fmt.Sprintf("[%v@%v:%v  prop: %v] ", nd.Name, nd.Scope, nd.Type, prop)
}

type EndpointDefs struct {
	Nodes []*NodeDef
}

type CallActionDefs struct {
	Node *NodeDef
	Cmd  string
}

type CastActionDefs struct {
	CallActionDefs
}

type SinkActionDefs struct {
	ChannelName string
}
