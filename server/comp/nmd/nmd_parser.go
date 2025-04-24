// Code generated from nmd.g4 by ANTLR 4.9.2. DO NOT EDIT.

package nmd // nmd
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 22, 131,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 7, 3,
	36, 10, 3, 12, 3, 14, 3, 39, 11, 3, 3, 3, 5, 3, 42, 10, 3, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 5, 4, 49, 10, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5,
	7, 5, 57, 10, 5, 12, 5, 14, 5, 60, 11, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7,
	3, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 7, 9,
	78, 10, 9, 12, 9, 14, 9, 81, 11, 9, 3, 9, 3, 9, 5, 9, 85, 10, 9, 3, 10,
	3, 10, 3, 10, 7, 10, 90, 10, 10, 12, 10, 14, 10, 93, 11, 10, 3, 10, 3,
	10, 3, 11, 3, 11, 3, 11, 5, 11, 100, 10, 11, 3, 11, 3, 11, 5, 11, 104,
	10, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 13, 3, 13, 3, 13, 7, 13, 113, 10,
	13, 12, 13, 14, 13, 116, 11, 13, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 5,
	14, 123, 10, 14, 3, 15, 3, 15, 3, 15, 3, 15, 5, 15, 129, 10, 15, 3, 15,
	2, 2, 16, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 2, 2, 2,
	133, 2, 30, 3, 2, 2, 2, 4, 32, 3, 2, 2, 2, 6, 48, 3, 2, 2, 2, 8, 50, 3,
	2, 2, 2, 10, 61, 3, 2, 2, 2, 12, 65, 3, 2, 2, 2, 14, 69, 3, 2, 2, 2, 16,
	84, 3, 2, 2, 2, 18, 86, 3, 2, 2, 2, 20, 96, 3, 2, 2, 2, 22, 105, 3, 2,
	2, 2, 24, 109, 3, 2, 2, 2, 26, 122, 3, 2, 2, 2, 28, 128, 3, 2, 2, 2, 30,
	31, 5, 4, 3, 2, 31, 3, 3, 2, 2, 2, 32, 37, 5, 6, 4, 2, 33, 34, 7, 3, 2,
	2, 34, 36, 5, 6, 4, 2, 35, 33, 3, 2, 2, 2, 36, 39, 3, 2, 2, 2, 37, 35,
	3, 2, 2, 2, 37, 38, 3, 2, 2, 2, 38, 41, 3, 2, 2, 2, 39, 37, 3, 2, 2, 2,
	40, 42, 7, 3, 2, 2, 41, 40, 3, 2, 2, 2, 41, 42, 3, 2, 2, 2, 42, 5, 3, 2,
	2, 2, 43, 49, 5, 18, 10, 2, 44, 49, 5, 8, 5, 2, 45, 49, 5, 10, 6, 2, 46,
	49, 5, 12, 7, 2, 47, 49, 5, 14, 8, 2, 48, 43, 3, 2, 2, 2, 48, 44, 3, 2,
	2, 2, 48, 45, 3, 2, 2, 2, 48, 46, 3, 2, 2, 2, 48, 47, 3, 2, 2, 2, 49, 7,
	3, 2, 2, 2, 50, 51, 5, 16, 9, 2, 51, 52, 5, 26, 14, 2, 52, 58, 5, 16, 9,
	2, 53, 54, 5, 26, 14, 2, 54, 55, 5, 16, 9, 2, 55, 57, 3, 2, 2, 2, 56, 53,
	3, 2, 2, 2, 57, 60, 3, 2, 2, 2, 58, 56, 3, 2, 2, 2, 58, 59, 3, 2, 2, 2,
	59, 9, 3, 2, 2, 2, 60, 58, 3, 2, 2, 2, 61, 62, 5, 18, 10, 2, 62, 63, 7,
	4, 2, 2, 63, 64, 7, 18, 2, 2, 64, 11, 3, 2, 2, 2, 65, 66, 5, 18, 10, 2,
	66, 67, 7, 5, 2, 2, 67, 68, 7, 18, 2, 2, 68, 13, 3, 2, 2, 2, 69, 70, 7,
	6, 2, 2, 70, 71, 7, 22, 2, 2, 71, 15, 3, 2, 2, 2, 72, 85, 5, 18, 10, 2,
	73, 74, 7, 7, 2, 2, 74, 79, 5, 18, 10, 2, 75, 76, 7, 8, 2, 2, 76, 78, 5,
	18, 10, 2, 77, 75, 3, 2, 2, 2, 78, 81, 3, 2, 2, 2, 79, 77, 3, 2, 2, 2,
	79, 80, 3, 2, 2, 2, 80, 82, 3, 2, 2, 2, 81, 79, 3, 2, 2, 2, 82, 83, 7,
	9, 2, 2, 83, 85, 3, 2, 2, 2, 84, 72, 3, 2, 2, 2, 84, 73, 3, 2, 2, 2, 85,
	17, 3, 2, 2, 2, 86, 87, 7, 10, 2, 2, 87, 91, 5, 20, 11, 2, 88, 90, 5, 22,
	12, 2, 89, 88, 3, 2, 2, 2, 90, 93, 3, 2, 2, 2, 91, 89, 3, 2, 2, 2, 91,
	92, 3, 2, 2, 2, 92, 94, 3, 2, 2, 2, 93, 91, 3, 2, 2, 2, 94, 95, 7, 11,
	2, 2, 95, 19, 3, 2, 2, 2, 96, 99, 7, 22, 2, 2, 97, 98, 7, 12, 2, 2, 98,
	100, 7, 22, 2, 2, 99, 97, 3, 2, 2, 2, 99, 100, 3, 2, 2, 2, 100, 103, 3,
	2, 2, 2, 101, 102, 7, 13, 2, 2, 102, 104, 7, 22, 2, 2, 103, 101, 3, 2,
	2, 2, 103, 104, 3, 2, 2, 2, 104, 21, 3, 2, 2, 2, 105, 106, 7, 22, 2, 2,
	106, 107, 7, 14, 2, 2, 107, 108, 5, 28, 15, 2, 108, 23, 3, 2, 2, 2, 109,
	114, 7, 22, 2, 2, 110, 111, 7, 8, 2, 2, 111, 113, 7, 22, 2, 2, 112, 110,
	3, 2, 2, 2, 113, 116, 3, 2, 2, 2, 114, 112, 3, 2, 2, 2, 114, 115, 3, 2,
	2, 2, 115, 25, 3, 2, 2, 2, 116, 114, 3, 2, 2, 2, 117, 118, 7, 15, 2, 2,
	118, 119, 5, 24, 13, 2, 119, 120, 7, 16, 2, 2, 120, 123, 3, 2, 2, 2, 121,
	123, 7, 17, 2, 2, 122, 117, 3, 2, 2, 2, 122, 121, 3, 2, 2, 2, 123, 27,
	3, 2, 2, 2, 124, 129, 7, 18, 2, 2, 125, 129, 7, 22, 2, 2, 126, 129, 7,
	19, 2, 2, 127, 129, 7, 20, 2, 2, 128, 124, 3, 2, 2, 2, 128, 125, 3, 2,
	2, 2, 128, 126, 3, 2, 2, 2, 128, 127, 3, 2, 2, 2, 129, 29, 3, 2, 2, 2,
	14, 37, 41, 48, 58, 79, 84, 91, 99, 103, 114, 122, 128,
}
var literalNames = []string{
	"", "';'", "'<->'", "'<--'", "'<-chan'", "'{'", "','", "'}'", "'['", "']'",
	"'@'", "':'", "'='", "'<'", "'>'", "'->'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "QUOTED_STRING",
	"INT", "FLOAT", "WS", "ID",
}

var ruleNames = []string{
	"graph", "stmt_list", "stmt", "link_stmt", "call_stmt", "cast_stmt", "sink_stmt",
	"endpoint", "node_def", "node_id", "node_prop", "msg_type_list", "link_operator",
	"property",
}

type nmdParser struct {
	*antlr.BaseParser
}

// NewnmdParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *nmdParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewnmdParser(input antlr.TokenStream) *nmdParser {
	this := new(nmdParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "nmd.g4"

	return this
}

// nmdParser tokens.
const (
	nmdParserEOF           = antlr.TokenEOF
	nmdParserT__0          = 1
	nmdParserT__1          = 2
	nmdParserT__2          = 3
	nmdParserT__3          = 4
	nmdParserT__4          = 5
	nmdParserT__5          = 6
	nmdParserT__6          = 7
	nmdParserT__7          = 8
	nmdParserT__8          = 9
	nmdParserT__9          = 10
	nmdParserT__10         = 11
	nmdParserT__11         = 12
	nmdParserT__12         = 13
	nmdParserT__13         = 14
	nmdParserT__14         = 15
	nmdParserQUOTED_STRING = 16
	nmdParserINT           = 17
	nmdParserFLOAT         = 18
	nmdParserWS            = 19
	nmdParserID            = 20
)

// nmdParser rules.
const (
	nmdParserRULE_graph         = 0
	nmdParserRULE_stmt_list     = 1
	nmdParserRULE_stmt          = 2
	nmdParserRULE_link_stmt     = 3
	nmdParserRULE_call_stmt     = 4
	nmdParserRULE_cast_stmt     = 5
	nmdParserRULE_sink_stmt     = 6
	nmdParserRULE_endpoint      = 7
	nmdParserRULE_node_def      = 8
	nmdParserRULE_node_id       = 9
	nmdParserRULE_node_prop     = 10
	nmdParserRULE_msg_type_list = 11
	nmdParserRULE_link_operator = 12
	nmdParserRULE_property      = 13
)

// IGraphContext is an interface to support dynamic dispatch.
type IGraphContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGraphContext differentiates from other interfaces.
	IsGraphContext()
}

type GraphContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGraphContext() *GraphContext {
	var p = new(GraphContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_graph
	return p
}

func (*GraphContext) IsGraphContext() {}

func NewGraphContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GraphContext {
	var p = new(GraphContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_graph

	return p
}

func (s *GraphContext) GetParser() antlr.Parser { return s.parser }

func (s *GraphContext) Stmt_list() IStmt_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmt_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmt_listContext)
}

func (s *GraphContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GraphContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GraphContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterGraph(s)
	}
}

func (s *GraphContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitGraph(s)
	}
}

func (p *nmdParser) Graph() (localctx IGraphContext) {
	localctx = NewGraphContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, nmdParserRULE_graph)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(28)
		p.Stmt_list()
	}

	return localctx
}

// IStmt_listContext is an interface to support dynamic dispatch.
type IStmt_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStmt_listContext differentiates from other interfaces.
	IsStmt_listContext()
}

type Stmt_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmt_listContext() *Stmt_listContext {
	var p = new(Stmt_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_stmt_list
	return p
}

func (*Stmt_listContext) IsStmt_listContext() {}

func NewStmt_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Stmt_listContext {
	var p = new(Stmt_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_stmt_list

	return p
}

func (s *Stmt_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Stmt_listContext) AllStmt() []IStmtContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStmtContext)(nil)).Elem())
	var tst = make([]IStmtContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStmtContext)
		}
	}

	return tst
}

func (s *Stmt_listContext) Stmt(i int) IStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStmtContext)
}

func (s *Stmt_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Stmt_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Stmt_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterStmt_list(s)
	}
}

func (s *Stmt_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitStmt_list(s)
	}
}

func (p *nmdParser) Stmt_list() (localctx IStmt_listContext) {
	localctx = NewStmt_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, nmdParserRULE_stmt_list)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(30)
		p.Stmt()
	}
	p.SetState(35)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(31)
				p.Match(nmdParserT__0)
			}
			{
				p.SetState(32)
				p.Stmt()
			}

		}
		p.SetState(37)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	p.SetState(39)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == nmdParserT__0 {
		{
			p.SetState(38)
			p.Match(nmdParserT__0)
		}

	}

	return localctx
}

// IStmtContext is an interface to support dynamic dispatch.
type IStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStmtContext differentiates from other interfaces.
	IsStmtContext()
}

type StmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtContext() *StmtContext {
	var p = new(StmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_stmt
	return p
}

func (*StmtContext) IsStmtContext() {}

func NewStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtContext {
	var p = new(StmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_stmt

	return p
}

func (s *StmtContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtContext) Node_def() INode_defContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INode_defContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INode_defContext)
}

func (s *StmtContext) Link_stmt() ILink_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILink_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILink_stmtContext)
}

func (s *StmtContext) Call_stmt() ICall_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICall_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICall_stmtContext)
}

func (s *StmtContext) Cast_stmt() ICast_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICast_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICast_stmtContext)
}

func (s *StmtContext) Sink_stmt() ISink_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISink_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISink_stmtContext)
}

func (s *StmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterStmt(s)
	}
}

func (s *StmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitStmt(s)
	}
}

func (p *nmdParser) Stmt() (localctx IStmtContext) {
	localctx = NewStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, nmdParserRULE_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(41)
			p.Node_def()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(42)
			p.Link_stmt()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(43)
			p.Call_stmt()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(44)
			p.Cast_stmt()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(45)
			p.Sink_stmt()
		}

	}

	return localctx
}

// ILink_stmtContext is an interface to support dynamic dispatch.
type ILink_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLink_stmtContext differentiates from other interfaces.
	IsLink_stmtContext()
}

type Link_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLink_stmtContext() *Link_stmtContext {
	var p = new(Link_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_link_stmt
	return p
}

func (*Link_stmtContext) IsLink_stmtContext() {}

func NewLink_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Link_stmtContext {
	var p = new(Link_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_link_stmt

	return p
}

func (s *Link_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Link_stmtContext) AllEndpoint() []IEndpointContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEndpointContext)(nil)).Elem())
	var tst = make([]IEndpointContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEndpointContext)
		}
	}

	return tst
}

func (s *Link_stmtContext) Endpoint(i int) IEndpointContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEndpointContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEndpointContext)
}

func (s *Link_stmtContext) AllLink_operator() []ILink_operatorContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILink_operatorContext)(nil)).Elem())
	var tst = make([]ILink_operatorContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILink_operatorContext)
		}
	}

	return tst
}

func (s *Link_stmtContext) Link_operator(i int) ILink_operatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILink_operatorContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ILink_operatorContext)
}

func (s *Link_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Link_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Link_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterLink_stmt(s)
	}
}

func (s *Link_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitLink_stmt(s)
	}
}

func (p *nmdParser) Link_stmt() (localctx ILink_stmtContext) {
	localctx = NewLink_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, nmdParserRULE_link_stmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(48)
		p.Endpoint()
	}
	{
		p.SetState(49)
		p.Link_operator()
	}
	{
		p.SetState(50)
		p.Endpoint()
	}
	p.SetState(56)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == nmdParserT__12 || _la == nmdParserT__14 {
		{
			p.SetState(51)
			p.Link_operator()
		}
		{
			p.SetState(52)
			p.Endpoint()
		}

		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ICall_stmtContext is an interface to support dynamic dispatch.
type ICall_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetCmd returns the cmd token.
	GetCmd() antlr.Token

	// SetCmd sets the cmd token.
	SetCmd(antlr.Token)

	// IsCall_stmtContext differentiates from other interfaces.
	IsCall_stmtContext()
}

type Call_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	cmd    antlr.Token
}

func NewEmptyCall_stmtContext() *Call_stmtContext {
	var p = new(Call_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_call_stmt
	return p
}

func (*Call_stmtContext) IsCall_stmtContext() {}

func NewCall_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Call_stmtContext {
	var p = new(Call_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_call_stmt

	return p
}

func (s *Call_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Call_stmtContext) GetCmd() antlr.Token { return s.cmd }

func (s *Call_stmtContext) SetCmd(v antlr.Token) { s.cmd = v }

func (s *Call_stmtContext) Node_def() INode_defContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INode_defContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INode_defContext)
}

func (s *Call_stmtContext) QUOTED_STRING() antlr.TerminalNode {
	return s.GetToken(nmdParserQUOTED_STRING, 0)
}

func (s *Call_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Call_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Call_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterCall_stmt(s)
	}
}

func (s *Call_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitCall_stmt(s)
	}
}

func (p *nmdParser) Call_stmt() (localctx ICall_stmtContext) {
	localctx = NewCall_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, nmdParserRULE_call_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(59)
		p.Node_def()
	}
	{
		p.SetState(60)
		p.Match(nmdParserT__1)
	}
	{
		p.SetState(61)

		var _m = p.Match(nmdParserQUOTED_STRING)

		localctx.(*Call_stmtContext).cmd = _m
	}

	return localctx
}

// ICast_stmtContext is an interface to support dynamic dispatch.
type ICast_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetCmd returns the cmd token.
	GetCmd() antlr.Token

	// SetCmd sets the cmd token.
	SetCmd(antlr.Token)

	// IsCast_stmtContext differentiates from other interfaces.
	IsCast_stmtContext()
}

type Cast_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	cmd    antlr.Token
}

func NewEmptyCast_stmtContext() *Cast_stmtContext {
	var p = new(Cast_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_cast_stmt
	return p
}

func (*Cast_stmtContext) IsCast_stmtContext() {}

func NewCast_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Cast_stmtContext {
	var p = new(Cast_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_cast_stmt

	return p
}

func (s *Cast_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Cast_stmtContext) GetCmd() antlr.Token { return s.cmd }

func (s *Cast_stmtContext) SetCmd(v antlr.Token) { s.cmd = v }

func (s *Cast_stmtContext) Node_def() INode_defContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INode_defContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INode_defContext)
}

func (s *Cast_stmtContext) QUOTED_STRING() antlr.TerminalNode {
	return s.GetToken(nmdParserQUOTED_STRING, 0)
}

func (s *Cast_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Cast_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Cast_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterCast_stmt(s)
	}
}

func (s *Cast_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitCast_stmt(s)
	}
}

func (p *nmdParser) Cast_stmt() (localctx ICast_stmtContext) {
	localctx = NewCast_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, nmdParserRULE_cast_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(63)
		p.Node_def()
	}
	{
		p.SetState(64)
		p.Match(nmdParserT__2)
	}
	{
		p.SetState(65)

		var _m = p.Match(nmdParserQUOTED_STRING)

		localctx.(*Cast_stmtContext).cmd = _m
	}

	return localctx
}

// ISink_stmtContext is an interface to support dynamic dispatch.
type ISink_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetNode returns the node token.
	GetNode() antlr.Token

	// SetNode sets the node token.
	SetNode(antlr.Token)

	// IsSink_stmtContext differentiates from other interfaces.
	IsSink_stmtContext()
}

type Sink_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	node   antlr.Token
}

func NewEmptySink_stmtContext() *Sink_stmtContext {
	var p = new(Sink_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_sink_stmt
	return p
}

func (*Sink_stmtContext) IsSink_stmtContext() {}

func NewSink_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Sink_stmtContext {
	var p = new(Sink_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_sink_stmt

	return p
}

func (s *Sink_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Sink_stmtContext) GetNode() antlr.Token { return s.node }

func (s *Sink_stmtContext) SetNode(v antlr.Token) { s.node = v }

func (s *Sink_stmtContext) ID() antlr.TerminalNode {
	return s.GetToken(nmdParserID, 0)
}

func (s *Sink_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Sink_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Sink_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterSink_stmt(s)
	}
}

func (s *Sink_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitSink_stmt(s)
	}
}

func (p *nmdParser) Sink_stmt() (localctx ISink_stmtContext) {
	localctx = NewSink_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, nmdParserRULE_sink_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(67)
		p.Match(nmdParserT__3)
	}
	{
		p.SetState(68)

		var _m = p.Match(nmdParserID)

		localctx.(*Sink_stmtContext).node = _m
	}

	return localctx
}

// IEndpointContext is an interface to support dynamic dispatch.
type IEndpointContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEndpointContext differentiates from other interfaces.
	IsEndpointContext()
}

type EndpointContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEndpointContext() *EndpointContext {
	var p = new(EndpointContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_endpoint
	return p
}

func (*EndpointContext) IsEndpointContext() {}

func NewEndpointContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EndpointContext {
	var p = new(EndpointContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_endpoint

	return p
}

func (s *EndpointContext) GetParser() antlr.Parser { return s.parser }

func (s *EndpointContext) AllNode_def() []INode_defContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INode_defContext)(nil)).Elem())
	var tst = make([]INode_defContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INode_defContext)
		}
	}

	return tst
}

func (s *EndpointContext) Node_def(i int) INode_defContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INode_defContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INode_defContext)
}

func (s *EndpointContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EndpointContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EndpointContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterEndpoint(s)
	}
}

func (s *EndpointContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitEndpoint(s)
	}
}

func (p *nmdParser) Endpoint() (localctx IEndpointContext) {
	localctx = NewEndpointContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, nmdParserRULE_endpoint)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(82)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case nmdParserT__7:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(70)
			p.Node_def()
		}

	case nmdParserT__4:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(71)
			p.Match(nmdParserT__4)
		}
		{
			p.SetState(72)
			p.Node_def()
		}
		p.SetState(77)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == nmdParserT__5 {
			{
				p.SetState(73)
				p.Match(nmdParserT__5)
			}
			{
				p.SetState(74)
				p.Node_def()
			}

			p.SetState(79)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(80)
			p.Match(nmdParserT__6)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// INode_defContext is an interface to support dynamic dispatch.
type INode_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNode_defContext differentiates from other interfaces.
	IsNode_defContext()
}

type Node_defContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNode_defContext() *Node_defContext {
	var p = new(Node_defContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_node_def
	return p
}

func (*Node_defContext) IsNode_defContext() {}

func NewNode_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Node_defContext {
	var p = new(Node_defContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_node_def

	return p
}

func (s *Node_defContext) GetParser() antlr.Parser { return s.parser }

func (s *Node_defContext) Node_id() INode_idContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INode_idContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INode_idContext)
}

func (s *Node_defContext) AllNode_prop() []INode_propContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INode_propContext)(nil)).Elem())
	var tst = make([]INode_propContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INode_propContext)
		}
	}

	return tst
}

func (s *Node_defContext) Node_prop(i int) INode_propContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INode_propContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INode_propContext)
}

func (s *Node_defContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Node_defContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Node_defContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterNode_def(s)
	}
}

func (s *Node_defContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitNode_def(s)
	}
}

func (p *nmdParser) Node_def() (localctx INode_defContext) {
	localctx = NewNode_defContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, nmdParserRULE_node_def)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(84)
		p.Match(nmdParserT__7)
	}
	{
		p.SetState(85)
		p.Node_id()
	}
	p.SetState(89)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == nmdParserID {
		{
			p.SetState(86)
			p.Node_prop()
		}

		p.SetState(91)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(92)
		p.Match(nmdParserT__8)
	}

	return localctx
}

// INode_idContext is an interface to support dynamic dispatch.
type INode_idContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name token.
	GetName() antlr.Token

	// GetScope returns the scope token.
	GetScope() antlr.Token

	// GetTyp returns the typ token.
	GetTyp() antlr.Token

	// SetName sets the name token.
	SetName(antlr.Token)

	// SetScope sets the scope token.
	SetScope(antlr.Token)

	// SetTyp sets the typ token.
	SetTyp(antlr.Token)

	// IsNode_idContext differentiates from other interfaces.
	IsNode_idContext()
}

type Node_idContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	name   antlr.Token
	scope  antlr.Token
	typ    antlr.Token
}

func NewEmptyNode_idContext() *Node_idContext {
	var p = new(Node_idContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_node_id
	return p
}

func (*Node_idContext) IsNode_idContext() {}

func NewNode_idContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Node_idContext {
	var p = new(Node_idContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_node_id

	return p
}

func (s *Node_idContext) GetParser() antlr.Parser { return s.parser }

func (s *Node_idContext) GetName() antlr.Token { return s.name }

func (s *Node_idContext) GetScope() antlr.Token { return s.scope }

func (s *Node_idContext) GetTyp() antlr.Token { return s.typ }

func (s *Node_idContext) SetName(v antlr.Token) { s.name = v }

func (s *Node_idContext) SetScope(v antlr.Token) { s.scope = v }

func (s *Node_idContext) SetTyp(v antlr.Token) { s.typ = v }

func (s *Node_idContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(nmdParserID)
}

func (s *Node_idContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(nmdParserID, i)
}

func (s *Node_idContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Node_idContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Node_idContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterNode_id(s)
	}
}

func (s *Node_idContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitNode_id(s)
	}
}

func (p *nmdParser) Node_id() (localctx INode_idContext) {
	localctx = NewNode_idContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, nmdParserRULE_node_id)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(94)

		var _m = p.Match(nmdParserID)

		localctx.(*Node_idContext).name = _m
	}
	p.SetState(97)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == nmdParserT__9 {
		{
			p.SetState(95)
			p.Match(nmdParserT__9)
		}
		{
			p.SetState(96)

			var _m = p.Match(nmdParserID)

			localctx.(*Node_idContext).scope = _m
		}

	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == nmdParserT__10 {
		{
			p.SetState(99)
			p.Match(nmdParserT__10)
		}
		{
			p.SetState(100)

			var _m = p.Match(nmdParserID)

			localctx.(*Node_idContext).typ = _m
		}

	}

	return localctx
}

// INode_propContext is an interface to support dynamic dispatch.
type INode_propContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetKey returns the key token.
	GetKey() antlr.Token

	// SetKey sets the key token.
	SetKey(antlr.Token)

	// GetValue returns the value rule contexts.
	GetValue() IPropertyContext

	// SetValue sets the value rule contexts.
	SetValue(IPropertyContext)

	// IsNode_propContext differentiates from other interfaces.
	IsNode_propContext()
}

type Node_propContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	key    antlr.Token
	value  IPropertyContext
}

func NewEmptyNode_propContext() *Node_propContext {
	var p = new(Node_propContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_node_prop
	return p
}

func (*Node_propContext) IsNode_propContext() {}

func NewNode_propContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Node_propContext {
	var p = new(Node_propContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_node_prop

	return p
}

func (s *Node_propContext) GetParser() antlr.Parser { return s.parser }

func (s *Node_propContext) GetKey() antlr.Token { return s.key }

func (s *Node_propContext) SetKey(v antlr.Token) { s.key = v }

func (s *Node_propContext) GetValue() IPropertyContext { return s.value }

func (s *Node_propContext) SetValue(v IPropertyContext) { s.value = v }

func (s *Node_propContext) ID() antlr.TerminalNode {
	return s.GetToken(nmdParserID, 0)
}

func (s *Node_propContext) Property() IPropertyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPropertyContext)
}

func (s *Node_propContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Node_propContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Node_propContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterNode_prop(s)
	}
}

func (s *Node_propContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitNode_prop(s)
	}
}

func (p *nmdParser) Node_prop() (localctx INode_propContext) {
	localctx = NewNode_propContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, nmdParserRULE_node_prop)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)

		var _m = p.Match(nmdParserID)

		localctx.(*Node_propContext).key = _m
	}
	{
		p.SetState(104)
		p.Match(nmdParserT__11)
	}
	{
		p.SetState(105)

		var _x = p.Property()

		localctx.(*Node_propContext).value = _x
	}

	return localctx
}

// IMsg_type_listContext is an interface to support dynamic dispatch.
type IMsg_type_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMsg_type_listContext differentiates from other interfaces.
	IsMsg_type_listContext()
}

type Msg_type_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMsg_type_listContext() *Msg_type_listContext {
	var p = new(Msg_type_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_msg_type_list
	return p
}

func (*Msg_type_listContext) IsMsg_type_listContext() {}

func NewMsg_type_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Msg_type_listContext {
	var p = new(Msg_type_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_msg_type_list

	return p
}

func (s *Msg_type_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Msg_type_listContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(nmdParserID)
}

func (s *Msg_type_listContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(nmdParserID, i)
}

func (s *Msg_type_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Msg_type_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Msg_type_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterMsg_type_list(s)
	}
}

func (s *Msg_type_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitMsg_type_list(s)
	}
}

func (p *nmdParser) Msg_type_list() (localctx IMsg_type_listContext) {
	localctx = NewMsg_type_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, nmdParserRULE_msg_type_list)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(107)
		p.Match(nmdParserID)
	}
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == nmdParserT__5 {
		{
			p.SetState(108)
			p.Match(nmdParserT__5)
		}
		{
			p.SetState(109)
			p.Match(nmdParserID)
		}

		p.SetState(114)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ILink_operatorContext is an interface to support dynamic dispatch.
type ILink_operatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLink_operatorContext differentiates from other interfaces.
	IsLink_operatorContext()
}

type Link_operatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLink_operatorContext() *Link_operatorContext {
	var p = new(Link_operatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_link_operator
	return p
}

func (*Link_operatorContext) IsLink_operatorContext() {}

func NewLink_operatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Link_operatorContext {
	var p = new(Link_operatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_link_operator

	return p
}

func (s *Link_operatorContext) GetParser() antlr.Parser { return s.parser }

func (s *Link_operatorContext) Msg_type_list() IMsg_type_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMsg_type_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMsg_type_listContext)
}

func (s *Link_operatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Link_operatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Link_operatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterLink_operator(s)
	}
}

func (s *Link_operatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitLink_operator(s)
	}
}

func (p *nmdParser) Link_operator() (localctx ILink_operatorContext) {
	localctx = NewLink_operatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, nmdParserRULE_link_operator)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(120)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case nmdParserT__12:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(115)
			p.Match(nmdParserT__12)
		}
		{
			p.SetState(116)
			p.Msg_type_list()
		}
		{
			p.SetState(117)
			p.Match(nmdParserT__13)
		}

	case nmdParserT__14:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(119)
			p.Match(nmdParserT__14)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IPropertyContext is an interface to support dynamic dispatch.
type IPropertyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPropertyContext differentiates from other interfaces.
	IsPropertyContext()
}

type PropertyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertyContext() *PropertyContext {
	var p = new(PropertyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = nmdParserRULE_property
	return p
}

func (*PropertyContext) IsPropertyContext() {}

func NewPropertyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyContext {
	var p = new(PropertyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = nmdParserRULE_property

	return p
}

func (s *PropertyContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyContext) CopyFrom(ctx *PropertyContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *PropertyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type PropIdContext struct {
	*PropertyContext
}

func NewPropIdContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PropIdContext {
	var p = new(PropIdContext)

	p.PropertyContext = NewEmptyPropertyContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PropertyContext))

	return p
}

func (s *PropIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropIdContext) ID() antlr.TerminalNode {
	return s.GetToken(nmdParserID, 0)
}

func (s *PropIdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterPropId(s)
	}
}

func (s *PropIdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitPropId(s)
	}
}

type PropFloatContext struct {
	*PropertyContext
}

func NewPropFloatContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PropFloatContext {
	var p = new(PropFloatContext)

	p.PropertyContext = NewEmptyPropertyContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PropertyContext))

	return p
}

func (s *PropFloatContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropFloatContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(nmdParserFLOAT, 0)
}

func (s *PropFloatContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterPropFloat(s)
	}
}

func (s *PropFloatContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitPropFloat(s)
	}
}

type PropQuoteStringContext struct {
	*PropertyContext
}

func NewPropQuoteStringContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PropQuoteStringContext {
	var p = new(PropQuoteStringContext)

	p.PropertyContext = NewEmptyPropertyContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PropertyContext))

	return p
}

func (s *PropQuoteStringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropQuoteStringContext) QUOTED_STRING() antlr.TerminalNode {
	return s.GetToken(nmdParserQUOTED_STRING, 0)
}

func (s *PropQuoteStringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterPropQuoteString(s)
	}
}

func (s *PropQuoteStringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitPropQuoteString(s)
	}
}

type PropIntContext struct {
	*PropertyContext
}

func NewPropIntContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PropIntContext {
	var p = new(PropIntContext)

	p.PropertyContext = NewEmptyPropertyContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PropertyContext))

	return p
}

func (s *PropIntContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropIntContext) INT() antlr.TerminalNode {
	return s.GetToken(nmdParserINT, 0)
}

func (s *PropIntContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.EnterPropInt(s)
	}
}

func (s *PropIntContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nmdListener); ok {
		listenerT.ExitPropInt(s)
	}
}

func (p *nmdParser) Property() (localctx IPropertyContext) {
	localctx = NewPropertyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, nmdParserRULE_property)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(126)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case nmdParserQUOTED_STRING:
		localctx = NewPropQuoteStringContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(122)
			p.Match(nmdParserQUOTED_STRING)
		}

	case nmdParserID:
		localctx = NewPropIdContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(123)
			p.Match(nmdParserID)
		}

	case nmdParserINT:
		localctx = NewPropIntContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(124)
			p.Match(nmdParserINT)
		}

	case nmdParserFLOAT:
		localctx = NewPropFloatContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(125)
			p.Match(nmdParserFLOAT)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}
