package nmd

import (
	"errors"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/appcrash/media/server/utils"
)

// parser analyzes session's event graph description, collect each node's properties and build the DAG between
// nodes when session created

type GraphTopology struct {
	nodeDefs []*NodeDef
	callDefs []*CallActionDefs
	castDefs []*CastActionDefs
	sinkDefs []*SinkActionDefs

	// topographical sorted, filter node with less dependency comes first
	// all loop nodes come last
	// i.e. first filter node in the list has not any receiver
	sortedNodeDefs []*NodeDef
	nbParseError   int
}

// LoopNodeFilter returns true when node can be looped
// the default behaviour of nmd node is a filter, but some nodes can accept input from a sub-graph,
// meanwhile output to the same sub-graph, which effectively make the whole graph not a DAG(loop created).
// this filter identify these nodes and topographicalSort will ignore them, so that sorting can succeed.
type LoopNodeFilter func(nodeName string) bool

func NewGraphTopology() *GraphTopology {
	return &GraphTopology{}
}

func (gt *GraphTopology) ParseGraph(sessionId, desc string, loopFilter LoopNodeFilter) error {
	input := antlr.NewInputStream(desc)
	lexer := NewnmdLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewnmdParser(stream)
	listener := NewListener(sessionId)
	parser.AddErrorListener(listener)
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Graph())

	if listener.errorString != "" {
		return errors.New(listener.errorString)
	}

	gt.nodeDefs = listener.NodeDefs
	gt.callDefs = listener.CallDefs
	gt.castDefs = listener.CastDefs
	gt.sinkDefs = listener.SinkDefs

	return gt.topographicalSort(loopFilter)
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
// this sorting helps composer inspecting the dependency of nodes with filter behaviour. non-filter nodes
// are filtered out by loopFilter. this check is just a simple sanity check(filters should not form loop,
// otherwise data stream would flow infinitely), and can not prevent from buggy graph in runtime.
func (gt *GraphTopology) topographicalSort(loopFilter LoopNodeFilter) (err error) {
	n := len(gt.nodeDefs)
	loopNodeSet := utils.NewSet[int]()
	outDegree := make([]int, n) // for each node, how many nodes it connects to (it depends on them)
	var result []int
	mat := make([]uint8, n*n) // in fact: 2d-array

	if loopFilter == nil {
		loopFilter = func(_ string) bool {
			return false // not loop node
		}
	}

	// initialize the dependency matrix, ignore all non-filter nodes
	for _, fn := range gt.nodeDefs {
		from := fn.Index
		if loopFilter(fn.Name) {
			loopNodeSet.Add(from)
			outDegree[from] = -1 // don't choose this node when sorting
			continue
		}
		for _, tn := range fn.Deps {
			to := tn.Index
			if loopFilter(tn.Name) {
				loopNodeSet.Add(to)
				continue
			}
			outDegree[from]++
			mat[from*n+to] = 1
		}
	}
	filterNodeNb := n - loopNodeSet.Size()

	// iterate until all filter nodes sorted
	for len(result) < filterNodeNb {
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
			err = errors.New("graph topographical sort finds a loop in filter nodes")
			return
		}
	}

	// record the result node list
	var i int
	gt.sortedNodeDefs = make([]*NodeDef, n)
	for i = 0; i < filterNodeNb; i++ {
		gt.sortedNodeDefs[i] = gt.nodeDefs[result[i]]
	}
	// append loop nodes at the end of sorted node list
	for loopNodeSet.Size() > 0 {
		gt.sortedNodeDefs[i] = gt.nodeDefs[loopNodeSet.GetAndRemove()]
	}
	return
}
