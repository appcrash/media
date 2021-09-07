package event

// dlink is directed arrow that connect two nodes, event can only flow
// in one direction, if dlink is not used anymore, call tear down to
// notify the other side releasing it
type dlink struct {
	graph    *Graph
	name     string
	fromNode *NodeDelegate
	toNode   *NodeDelegate

	// index of nodeInfo's dlink array
	fromIndex int
	toIndex   int
}

func generateLinkName(fromScope string, fromNodeName string, toScope string, toNodeName string) string {
	return fromScope + ":" + fromNodeName + "#" + toScope + ":" + toNodeName
}

func newLink(graph *Graph, fromNode *NodeDelegate, toNode *NodeDelegate) *dlink {
	name := generateLinkName(fromNode.getNodeScope(), fromNode.getNodeName(),
		toNode.getNodeScope(), toNode.getNodeName())
	return &dlink{
		graph:     graph,
		name:      name,
		fromNode:  fromNode,
		toNode:    toNode,
		fromIndex: -1,
		toIndex:   -1,
	}
}
