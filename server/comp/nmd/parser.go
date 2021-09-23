package nmd

import (
	"errors"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"regexp"
	"strings"
)

// parser analyzes session's event graph description, collect each node's properties and build the DAG between
// nodes when session created

//------------------------- NodeProp -------------------------

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
	return fmt.Sprintf("%v@%v:%v  prop: %v", nd.Name, nd.Scope, nd.Type, prop)
}

type Endpoint struct {
	Nodes []*NodeDef
}

type LinkDef struct {
	From, To *NodeDef
}

type GraphTopology struct {
	nodeDefs       []*NodeDef
	sortedNodeDefs []*NodeDef // topographical sorted, node with less dependency comes first
	nbParseError   int
}

func NewGraphTopology() *GraphTopology {
	return &GraphTopology{}
}

func (gt *GraphTopology) ParseGraph(sessionId, desc string) error {
	input := antlr.NewInputStream(desc)
	lexer := NewnmdLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewnmdParser(stream)
	listener := NewListener(sessionId)
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Graph())

	gt.nodeDefs = listener.NodeDefs
	return gt.topographicalSort()
}

func (gt *GraphTopology) GetSortedNodeDefs() []*NodeDef {
	return gt.sortedNodeDefs
}

// O(n*n) sort algorithm, ok when n is small
func (gt *GraphTopology) topographicalSort() (err error) {
	n := len(gt.nodeDefs)
	outDegree := make([]int, n) // for each node, how many nodes it connects to (it depends on them)
	var result []int
	mat := make([]uint8, n*n) // in fact: 2d-array

	// initialize the dependency matrix
	for _, fn := range gt.nodeDefs {
		from := fn.Index
		for _, tn := range fn.Deps {
			to := tn.Index
			outDegree[from]++
			mat[from*n+to] = 1
		}
	}

	// iterate until all nodes sorted
	for len(result) < n {
		// search all nodes that have no dependency
		found := false
		for i := 0; i < n; i++ {
			if outDegree[i] == 0 {
				found = true
				outDegree[i] = -1 // choose it only once
				result = append(result, i)
				// the node removed, then delete all inward connections to it
				for j := 0; j < n; j++ {
					// check all bits in column(i)
					if mat[j*n+i] > 0 {
						outDegree[j]--
					}
				}
			}
		}
		if !found {
			// can not find a candidate that has no dependency, means loop detected
			err = errors.New("graph topographical sort finds a loop")
			return
		}
	}

	// record the result node list
	gt.sortedNodeDefs = make([]*NodeDef, n)
	for i := 0; i < n; i++ {
		gt.sortedNodeDefs[i] = gt.nodeDefs[result[i]]
	}
	return
}
