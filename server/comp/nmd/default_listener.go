package nmd

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strconv"
	"strings"
)

type Listener struct {
	*BasenmdListener
	*antlr.DefaultErrorListener

	sessionId string

	NodeDefs []*NodeDef
	CallDefs []*CallActionDefs
	CastDefs []*CastActionDefs
	SinkDefs []*SinkActionDefs

	nodeMap       map[string]*NodeDef
	nbNode        int
	nodeDefStack  []*NodeDef
	endpointStack []*EndpointDefs

	currentNodeProp *NodeProp
	currentNodeDef  *NodeDef
	currentEndpoint *EndpointDefs

	errorString string
}

func NewListener(sessionId string) *Listener {
	l := new(Listener)
	l.sessionId = sessionId
	l.nodeMap = make(map[string]*NodeDef)
	return l
}

func (l *Listener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.errorString += msg + "\n"
}

func unquoteString(quotedString string) string {
	ql := len(quotedString)
	if ql == 2 {
		// empty string ''
		return ""
	} else {
		// modify escaped single quote
		return strings.ReplaceAll(quotedString[1:ql-1], "\\'", "'")
	}
}

func (l *Listener) pushNodeDef() {
	l.nodeDefStack = append(l.nodeDefStack, l.currentNodeDef)
}

func (l *Listener) pushEndpoint() {
	l.endpointStack = append(l.endpointStack, l.currentEndpoint)
}

func (l *Listener) getNodeDef(name, scope, typ string) (ni *NodeDef) {
	var ok bool
	queryId := name + "_" + scope
	if ni, ok = l.nodeMap[queryId]; !ok {
		ni = &NodeDef{
			Name:  name,
			Scope: scope,
			Type:  typ,
			Index: l.nbNode,
		}
		l.nodeMap[queryId] = ni
		l.NodeDefs = append(l.NodeDefs, ni)
		l.nbNode++
	}
	//if typ != "" && typ != name && ni.TypeId != typ {
	//	// if this is not abbreviated naming, and explicitly provide type that conflicts with previous defined
	//	l.nbParseError++
	//}
	return
}

func (l *Listener) EnterLink_stmt(c *Link_stmtContext) {
	l.endpointStack = nil
}

func (l *Listener) EnterEndpoint(c *EndpointContext) {
	l.currentEndpoint = &EndpointDefs{}
	l.nodeDefStack = nil
}

func (l *Listener) EnterMsg_type_list(c *Msg_type_listContext) {
	for _, id := range c.AllID() {
		l.currentEndpoint.PreferOffer = append(l.currentEndpoint.PreferOffer, id.GetText())
	}
}

func (l *Listener) EnterNode_id(c *Node_idContext) {
	var name, scope, typ string
	name = c.GetName().GetText()
	if c.scope != nil {
		scope = c.GetScope().GetText()
	} else {
		scope = l.sessionId
	}
	if c.typ != nil {
		typ = c.GetTyp().GetText()
	} else {
		// type is omitted, use the name as type
		typ = name
	}
	l.currentNodeDef = l.getNodeDef(name, scope, typ)
}

func (l *Listener) EnterNode_prop(c *Node_propContext) {
	l.currentNodeProp = &NodeProp{Key: c.GetKey().GetText()}
}

func (l *Listener) EnterPropQuoteString(ctx *PropQuoteStringContext) {
	l.currentNodeProp.Type = "str"
	l.currentNodeProp.Value = unquoteString(ctx.GetText())
}

func (l *Listener) EnterPropId(ctx *PropIdContext) {
	l.currentNodeProp.Type = "str"
	l.currentNodeProp.Value = ctx.GetText()
}

func (l *Listener) EnterPropInt(ctx *PropIntContext) {
	l.currentNodeProp.Type = "int"
	text := ctx.GetText()
	if strings.HasPrefix(text, "0x") || strings.HasPrefix(text, "0X") {
		l.currentNodeProp.Value, _ = strconv.ParseInt(text[2:], 16, 64)
	} else {
		l.currentNodeProp.Value, _ = strconv.Atoi(text)
	}
}

func (l *Listener) EnterPropFloat(ctx *PropFloatContext) {
	l.currentNodeProp.Type = "float"
	l.currentNodeProp.Value, _ = strconv.ParseFloat(ctx.GetText(), 64)
}

func (l *Listener) ExitGraph(c *GraphContext) {
	l.currentEndpoint = nil
	l.currentNodeProp = nil
	l.currentNodeDef = nil
}

func (l *Listener) ExitEndpoint(c *EndpointContext) {
	l.currentEndpoint.Nodes = l.nodeDefStack
	l.pushEndpoint()

	if len(l.endpointStack) == 2 {
		// a link is found
		from, to := l.endpointStack[0], l.endpointStack[1]
		preferOffer := from.PreferOffer
		for _, f := range from.Nodes {
			for _, t := range to.Nodes {
				linkOperator := &LinkOperator{
					LinkTo:      t,
					PreferOffer: preferOffer,
				}
				f.Deps = append(f.Deps, linkOperator)
			}
		}
		// discard the first endpoint
		l.endpointStack = l.endpointStack[1:]
	}
}

func (l *Listener) ExitNode_def(c *Node_defContext) {
	l.pushNodeDef()
}

func (l *Listener) ExitNode_prop(c *Node_propContext) {
	l.currentNodeProp.FormalizeKey()
	l.currentNodeDef.Props = append(l.currentNodeDef.Props, l.currentNodeProp)
}

func (l *Listener) ExitCall_stmt(ctx *Call_stmtContext) {
	node := l.nodeDefStack[0]
	if ctx.cmd == nil {
		return
	}
	cmd := unquoteString(ctx.cmd.GetText())
	call := &CallActionDefs{
		Node: node,
		Cmd:  cmd,
	}
	l.CallDefs = append(l.CallDefs, call)
	l.nodeDefStack = nil
}

func (l *Listener) ExitCast_stmt(ctx *Cast_stmtContext) {
	node := l.nodeDefStack[0]
	if ctx.cmd == nil {
		return
	}
	cmd := unquoteString(ctx.cmd.GetText())
	cast := &CastActionDefs{}
	cast.Node = node
	cast.Cmd = cmd
	l.CastDefs = append(l.CastDefs, cast)
	l.nodeDefStack = nil
}

func (l *Listener) ExitSink_stmt(ctx *Sink_stmtContext) {
	l.SinkDefs = append(l.SinkDefs, &SinkActionDefs{ctx.node.GetText()})
}
