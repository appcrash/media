package comp

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Composer analyzes session's event graph description, collect each node's properties and build the DAG between
// nodes when session created, we create every node and config their properties based on the collected info, then
// link them and send initial events as described before putting them into working state.

// graph description, line by line, separated by '\n'
// each line is either a node property or node connect description:
// examples:
// node property:  [nodeName:{nodeType}]:  key1=value1;key2=value2;...
// [asr]: sampleRate=16k       a node of both name and type are asr
// [ps1:pubsub]: channel=source   a node of name(ps1) and type(pubsub)
// node connect:   [nodeNameA] -> [nodeNameB]
// [transcode] -> [asr]        establish a link from transcode to asr when initialization
const regNodePropertyPattern = `^\s*\[([\w:]+)\]:\s*(\S+)\s*$`
const regNodeConnectPattern = `^\s*\[(\w+)\]\s*->\s*\[(\w+)\]\s*$`

var regNodeProperty = regexp.MustCompile(regNodePropertyPattern)
var regNodeConnect = regexp.MustCompile(regNodeConnectPattern)

type NodeInfo struct {
	Index  int
	Type   string
	Name   string
	Props  ConfigItems
	NbDeps int
	Deps   []*NodeInfo // record receivers of this node
}

type LinkInfo struct {
	From, To *NodeInfo
}

type GraphTopology struct {
	nodeMap        map[string]*NodeInfo
	nodeList       []*NodeInfo
	sortedNodeList []*NodeInfo // topographical sorted, node with less dependency comes first
	nbNode         int
	nbParseError   int
}

func newGraphTopology() *GraphTopology {
	return &GraphTopology{
		nodeMap: make(map[string]*NodeInfo),
	}
}

func (gi *GraphTopology) parseLine(line string) {
	var match []string
	match = regNodeProperty.FindStringSubmatch(line)
	if match != nil {
		gi.parseNodeProperty(match[1], match[2])
	} else if match = regNodeConnect.FindStringSubmatch(line); match != nil {
		gi.parseNodeConnect(match[1], match[2])
	} else {
		logger.Errorf("wrong line when parse graph desc: %v\n", line)
	}
}

func (gi *GraphTopology) getNodeInfo(name string, typ string) (ni *NodeInfo) {
	var ok bool
	if ni, ok = gi.nodeMap[name]; !ok {
		if typ == "" {
			// node connect comes before node definition, so type is same as name
			typ = name
		}
		ni = &NodeInfo{
			Name:  name,
			Type:  typ,
			Index: gi.nbNode,
			Props: make(ConfigItems),
		}
		gi.nodeMap[name] = ni
		gi.nodeList = append(gi.nodeList, ni)
		gi.nbNode++
	}
	if typ != "" && ni.Type != typ {
		logger.Errorf("conflict type of node with the same name:%v", name)
		gi.nbParseError++
	}
	return
}

func (gi *GraphTopology) parseNodeProperty(nodeNameType string, props string) {
	// check if name:type pair is specified
	var nodeName, nodeType string
	pair := strings.Split(nodeNameType, ":")
	if len(pair) == 1 {
		nodeName, nodeType = pair[0], pair[0]
	} else {
		nodeName, nodeType = pair[0], pair[1]
	}
	ni := gi.getNodeInfo(nodeName, nodeType)
	for _, kvPair := range strings.Split(props, ";") {
		kv := strings.Split(kvPair, "=")
		if len(kv) != 2 {
			logger.Errorf("wrong node property: %v\n", kv)
			gi.nbParseError++
			continue
		}
		// check if value is an integer
		if v, err := strconv.Atoi(kv[1]); err == nil {
			ni.Props.Set(kv[0], v)
		} else {
			// value is string
			ni.Props.Set(kv[0], kv[1])
		}
	}
}

// O(nxn) sort algorithm, ok when n is small
func (gi *GraphTopology) topographicalSort() (err error) {
	n := gi.nbNode
	outDegree := make([]int, n) // for each node, how many nodes it connects to (it depends on them)
	var result []int
	mat := make([]uint8, n*n) // in fact: 2d-array

	// initialize the dependency matrix
	for _, fn := range gi.nodeList {
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
	gi.sortedNodeList = make([]*NodeInfo, n)
	for i := 0; i < n; i++ {
		gi.sortedNodeList[i] = gi.nodeList[result[i]]
	}
	return
}

func (gi *GraphTopology) parseNodeConnect(nodeFrom string, nodeTo string) {
	from := gi.getNodeInfo(nodeFrom, "")
	to := gi.getNodeInfo(nodeTo, "")
	if from == nil || to == nil {
		logger.Errorf("wrong link from %v to %v\n", nodeFrom, nodeTo)
		gi.nbParseError++
		return
	}
	from.Deps = append(from.Deps, to)
	from.NbDeps++
}

func (gi *GraphTopology) getLinkInfo() (li []LinkInfo) {
	for _, from := range gi.nodeList {
		for _, to := range from.Deps {
			li = append(li, LinkInfo{from, to})
		}
	}
	return
}
