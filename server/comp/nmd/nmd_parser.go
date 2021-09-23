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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 17, 91, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7, 4,
	8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 7, 3,
	26, 10, 3, 12, 3, 14, 3, 29, 11, 3, 3, 3, 5, 3, 32, 10, 3, 3, 4, 3, 4,
	5, 4, 36, 10, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 7, 5, 43, 10, 5, 12, 5,
	14, 5, 46, 11, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 7, 6, 53, 10, 6, 12, 6,
	14, 6, 56, 11, 6, 3, 6, 3, 6, 5, 6, 60, 10, 6, 3, 7, 3, 7, 3, 7, 7, 7,
	65, 10, 7, 12, 7, 14, 7, 68, 11, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 5, 8,
	75, 10, 8, 3, 8, 3, 8, 5, 8, 79, 10, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10,
	3, 10, 3, 10, 3, 10, 5, 10, 89, 10, 10, 3, 10, 2, 2, 11, 2, 4, 6, 8, 10,
	12, 14, 16, 18, 2, 2, 2, 93, 2, 20, 3, 2, 2, 2, 4, 22, 3, 2, 2, 2, 6, 35,
	3, 2, 2, 2, 8, 37, 3, 2, 2, 2, 10, 59, 3, 2, 2, 2, 12, 61, 3, 2, 2, 2,
	14, 71, 3, 2, 2, 2, 16, 80, 3, 2, 2, 2, 18, 88, 3, 2, 2, 2, 20, 21, 5,
	4, 3, 2, 21, 3, 3, 2, 2, 2, 22, 27, 5, 6, 4, 2, 23, 24, 7, 3, 2, 2, 24,
	26, 5, 6, 4, 2, 25, 23, 3, 2, 2, 2, 26, 29, 3, 2, 2, 2, 27, 25, 3, 2, 2,
	2, 27, 28, 3, 2, 2, 2, 28, 31, 3, 2, 2, 2, 29, 27, 3, 2, 2, 2, 30, 32,
	7, 3, 2, 2, 31, 30, 3, 2, 2, 2, 31, 32, 3, 2, 2, 2, 32, 5, 3, 2, 2, 2,
	33, 36, 5, 12, 7, 2, 34, 36, 5, 8, 5, 2, 35, 33, 3, 2, 2, 2, 35, 34, 3,
	2, 2, 2, 36, 7, 3, 2, 2, 2, 37, 38, 5, 10, 6, 2, 38, 39, 7, 4, 2, 2, 39,
	44, 5, 10, 6, 2, 40, 41, 7, 4, 2, 2, 41, 43, 5, 10, 6, 2, 42, 40, 3, 2,
	2, 2, 43, 46, 3, 2, 2, 2, 44, 42, 3, 2, 2, 2, 44, 45, 3, 2, 2, 2, 45, 9,
	3, 2, 2, 2, 46, 44, 3, 2, 2, 2, 47, 60, 5, 12, 7, 2, 48, 49, 7, 5, 2, 2,
	49, 54, 5, 12, 7, 2, 50, 51, 7, 6, 2, 2, 51, 53, 5, 12, 7, 2, 52, 50, 3,
	2, 2, 2, 53, 56, 3, 2, 2, 2, 54, 52, 3, 2, 2, 2, 54, 55, 3, 2, 2, 2, 55,
	57, 3, 2, 2, 2, 56, 54, 3, 2, 2, 2, 57, 58, 7, 7, 2, 2, 58, 60, 3, 2, 2,
	2, 59, 47, 3, 2, 2, 2, 59, 48, 3, 2, 2, 2, 60, 11, 3, 2, 2, 2, 61, 62,
	7, 8, 2, 2, 62, 66, 5, 14, 8, 2, 63, 65, 5, 16, 9, 2, 64, 63, 3, 2, 2,
	2, 65, 68, 3, 2, 2, 2, 66, 64, 3, 2, 2, 2, 66, 67, 3, 2, 2, 2, 67, 69,
	3, 2, 2, 2, 68, 66, 3, 2, 2, 2, 69, 70, 7, 9, 2, 2, 70, 13, 3, 2, 2, 2,
	71, 74, 7, 17, 2, 2, 72, 73, 7, 10, 2, 2, 73, 75, 7, 17, 2, 2, 74, 72,
	3, 2, 2, 2, 74, 75, 3, 2, 2, 2, 75, 78, 3, 2, 2, 2, 76, 77, 7, 11, 2, 2,
	77, 79, 7, 17, 2, 2, 78, 76, 3, 2, 2, 2, 78, 79, 3, 2, 2, 2, 79, 15, 3,
	2, 2, 2, 80, 81, 7, 17, 2, 2, 81, 82, 7, 12, 2, 2, 82, 83, 5, 18, 10, 2,
	83, 17, 3, 2, 2, 2, 84, 89, 7, 13, 2, 2, 85, 89, 7, 17, 2, 2, 86, 89, 7,
	14, 2, 2, 87, 89, 7, 15, 2, 2, 88, 84, 3, 2, 2, 2, 88, 85, 3, 2, 2, 2,
	88, 86, 3, 2, 2, 2, 88, 87, 3, 2, 2, 2, 89, 19, 3, 2, 2, 2, 12, 27, 31,
	35, 44, 54, 59, 66, 74, 78, 88,
}
var literalNames = []string{
	"", "';'", "'->'", "'{'", "','", "'}'", "'['", "']'", "'@'", "':'", "'='",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "QUOTED_STRING", "INT", "FLOAT",
	"WS", "ID",
}

var ruleNames = []string{
	"graph", "stmt_list", "stmt", "link_stmt", "endpoint", "node_def", "node_id",
	"node_prop", "property",
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
	nmdParserQUOTED_STRING = 11
	nmdParserINT           = 12
	nmdParserFLOAT         = 13
	nmdParserWS            = 14
	nmdParserID            = 15
)

// nmdParser rules.
const (
	nmdParserRULE_graph     = 0
	nmdParserRULE_stmt_list = 1
	nmdParserRULE_stmt      = 2
	nmdParserRULE_link_stmt = 3
	nmdParserRULE_endpoint  = 4
	nmdParserRULE_node_def  = 5
	nmdParserRULE_node_id   = 6
	nmdParserRULE_node_prop = 7
	nmdParserRULE_property  = 8
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
		p.SetState(18)
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
		p.SetState(20)
		p.Stmt()
	}
	p.SetState(25)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(21)
				p.Match(nmdParserT__0)
			}
			{
				p.SetState(22)
				p.Stmt()
			}

		}
		p.SetState(27)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	p.SetState(29)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == nmdParserT__0 {
		{
			p.SetState(28)
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

	p.SetState(33)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(31)
			p.Node_def()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(32)
			p.Link_stmt()
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
		p.SetState(35)
		p.Endpoint()
	}
	{
		p.SetState(36)
		p.Match(nmdParserT__1)
	}
	{
		p.SetState(37)
		p.Endpoint()
	}
	p.SetState(42)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == nmdParserT__1 {
		{
			p.SetState(38)
			p.Match(nmdParserT__1)
		}
		{
			p.SetState(39)
			p.Endpoint()
		}

		p.SetState(44)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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
	p.EnterRule(localctx, 8, nmdParserRULE_endpoint)
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

	p.SetState(57)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case nmdParserT__5:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(45)
			p.Node_def()
		}

	case nmdParserT__2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(46)
			p.Match(nmdParserT__2)
		}
		{
			p.SetState(47)
			p.Node_def()
		}
		p.SetState(52)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == nmdParserT__3 {
			{
				p.SetState(48)
				p.Match(nmdParserT__3)
			}
			{
				p.SetState(49)
				p.Node_def()
			}

			p.SetState(54)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(55)
			p.Match(nmdParserT__4)
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
	p.EnterRule(localctx, 10, nmdParserRULE_node_def)
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
		p.SetState(59)
		p.Match(nmdParserT__5)
	}
	{
		p.SetState(60)
		p.Node_id()
	}
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == nmdParserID {
		{
			p.SetState(61)
			p.Node_prop()
		}

		p.SetState(66)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(67)
		p.Match(nmdParserT__6)
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
	p.EnterRule(localctx, 12, nmdParserRULE_node_id)
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
		p.SetState(69)

		var _m = p.Match(nmdParserID)

		localctx.(*Node_idContext).name = _m
	}
	p.SetState(72)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == nmdParserT__7 {
		{
			p.SetState(70)
			p.Match(nmdParserT__7)
		}
		{
			p.SetState(71)

			var _m = p.Match(nmdParserID)

			localctx.(*Node_idContext).scope = _m
		}

	}
	p.SetState(76)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == nmdParserT__8 {
		{
			p.SetState(74)
			p.Match(nmdParserT__8)
		}
		{
			p.SetState(75)

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
	p.EnterRule(localctx, 14, nmdParserRULE_node_prop)

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
		p.SetState(78)

		var _m = p.Match(nmdParserID)

		localctx.(*Node_propContext).key = _m
	}
	{
		p.SetState(79)
		p.Match(nmdParserT__9)
	}
	{
		p.SetState(80)

		var _x = p.Property()

		localctx.(*Node_propContext).value = _x
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
	p.EnterRule(localctx, 16, nmdParserRULE_property)

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

	p.SetState(86)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case nmdParserQUOTED_STRING:
		localctx = NewPropQuoteStringContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(82)
			p.Match(nmdParserQUOTED_STRING)
		}

	case nmdParserID:
		localctx = NewPropIdContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(83)
			p.Match(nmdParserID)
		}

	case nmdParserINT:
		localctx = NewPropIntContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(84)
			p.Match(nmdParserINT)
		}

	case nmdParserFLOAT:
		localctx = NewPropFloatContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(85)
			p.Match(nmdParserFLOAT)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}
