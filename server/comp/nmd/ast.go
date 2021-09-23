package nmd

import (
	"fmt"
	"regexp"
	"strings"
)

type NodeProp struct {
	Key, Type string
	Value     interface{}
}

const regCamelCasePattern = `_+[a-z]`

var regCamelCase = regexp.MustCompile(regCamelCasePattern)

// Set converts key in form of foo_bar or foo__bar ... into fooBar if possible
// normal keys with camel case remain intact
func (np *NodeProp) Set(key, typ string, val interface{}) {
	newKey := regCamelCase.ReplaceAllStringFunc(key, func(match string) string {
		last := string(match[len(match)-1])
		return strings.ToUpper(last)
	})
	np.Key = newKey
	np.Type = typ
	np.Value = val
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
