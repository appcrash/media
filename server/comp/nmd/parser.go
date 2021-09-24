package nmd

import (
	"errors"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// parser analyzes session's event graph description, collect each node's properties and build the DAG between
// nodes when session created

type GraphTopology struct {
	nodeDefs       []*NodeDef
	callDefs       []*CallActionDefs
	castDefs       []*CastActionDefs
	sinkDefs       []*SinkActionDefs
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
	gt.callDefs = listener.CallDefs
	gt.castDefs = listener.CastDefs
	gt.sinkDefs = listener.SinkDefs
	return gt.topographicalSort()
}

func (gt *GraphTopology) GetSortedNodeDefs() []*NodeDef {
	return gt.sortedNodeDefs
}

func (gt *GraphTopology) GetCallActions() []*CallActionDefs {
	return gt.callDefs
}

func (gt *GraphTopology) GetCastActions() []*CastActionDefs {
	return gt.castDefs
}

func (gt *GraphTopology) GetSinkActions() []*SinkActionDefs {
	return gt.sinkDefs
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
