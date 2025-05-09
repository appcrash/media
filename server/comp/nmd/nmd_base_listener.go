// Code generated from nmd.g4 by ANTLR 4.9.2. DO NOT EDIT.

package nmd // nmd
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasenmdListener is a complete listener for a parse tree produced by nmdParser.
type BasenmdListener struct{}

var _ nmdListener = &BasenmdListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasenmdListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasenmdListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasenmdListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasenmdListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterGraph is called when production graph is entered.
func (s *BasenmdListener) EnterGraph(ctx *GraphContext) {}

// ExitGraph is called when production graph is exited.
func (s *BasenmdListener) ExitGraph(ctx *GraphContext) {}

// EnterStmt_list is called when production stmt_list is entered.
func (s *BasenmdListener) EnterStmt_list(ctx *Stmt_listContext) {}

// ExitStmt_list is called when production stmt_list is exited.
func (s *BasenmdListener) ExitStmt_list(ctx *Stmt_listContext) {}

// EnterStmt is called when production stmt is entered.
func (s *BasenmdListener) EnterStmt(ctx *StmtContext) {}

// ExitStmt is called when production stmt is exited.
func (s *BasenmdListener) ExitStmt(ctx *StmtContext) {}

// EnterLink_stmt is called when production link_stmt is entered.
func (s *BasenmdListener) EnterLink_stmt(ctx *Link_stmtContext) {}

// ExitLink_stmt is called when production link_stmt is exited.
func (s *BasenmdListener) ExitLink_stmt(ctx *Link_stmtContext) {}

// EnterCall_stmt is called when production call_stmt is entered.
func (s *BasenmdListener) EnterCall_stmt(ctx *Call_stmtContext) {}

// ExitCall_stmt is called when production call_stmt is exited.
func (s *BasenmdListener) ExitCall_stmt(ctx *Call_stmtContext) {}

// EnterCast_stmt is called when production cast_stmt is entered.
func (s *BasenmdListener) EnterCast_stmt(ctx *Cast_stmtContext) {}

// ExitCast_stmt is called when production cast_stmt is exited.
func (s *BasenmdListener) ExitCast_stmt(ctx *Cast_stmtContext) {}

// EnterSink_stmt is called when production sink_stmt is entered.
func (s *BasenmdListener) EnterSink_stmt(ctx *Sink_stmtContext) {}

// ExitSink_stmt is called when production sink_stmt is exited.
func (s *BasenmdListener) ExitSink_stmt(ctx *Sink_stmtContext) {}

// EnterEndpoint is called when production endpoint is entered.
func (s *BasenmdListener) EnterEndpoint(ctx *EndpointContext) {}

// ExitEndpoint is called when production endpoint is exited.
func (s *BasenmdListener) ExitEndpoint(ctx *EndpointContext) {}

// EnterNode_def is called when production node_def is entered.
func (s *BasenmdListener) EnterNode_def(ctx *Node_defContext) {}

// ExitNode_def is called when production node_def is exited.
func (s *BasenmdListener) ExitNode_def(ctx *Node_defContext) {}

// EnterNode_id is called when production node_id is entered.
func (s *BasenmdListener) EnterNode_id(ctx *Node_idContext) {}

// ExitNode_id is called when production node_id is exited.
func (s *BasenmdListener) ExitNode_id(ctx *Node_idContext) {}

// EnterNode_prop is called when production node_prop is entered.
func (s *BasenmdListener) EnterNode_prop(ctx *Node_propContext) {}

// ExitNode_prop is called when production node_prop is exited.
func (s *BasenmdListener) ExitNode_prop(ctx *Node_propContext) {}

// EnterMsg_type_list is called when production msg_type_list is entered.
func (s *BasenmdListener) EnterMsg_type_list(ctx *Msg_type_listContext) {}

// ExitMsg_type_list is called when production msg_type_list is exited.
func (s *BasenmdListener) ExitMsg_type_list(ctx *Msg_type_listContext) {}

// EnterLink_operator is called when production link_operator is entered.
func (s *BasenmdListener) EnterLink_operator(ctx *Link_operatorContext) {}

// ExitLink_operator is called when production link_operator is exited.
func (s *BasenmdListener) ExitLink_operator(ctx *Link_operatorContext) {}

// EnterPropQuoteString is called when production PropQuoteString is entered.
func (s *BasenmdListener) EnterPropQuoteString(ctx *PropQuoteStringContext) {}

// ExitPropQuoteString is called when production PropQuoteString is exited.
func (s *BasenmdListener) ExitPropQuoteString(ctx *PropQuoteStringContext) {}

// EnterPropId is called when production PropId is entered.
func (s *BasenmdListener) EnterPropId(ctx *PropIdContext) {}

// ExitPropId is called when production PropId is exited.
func (s *BasenmdListener) ExitPropId(ctx *PropIdContext) {}

// EnterPropInt is called when production PropInt is entered.
func (s *BasenmdListener) EnterPropInt(ctx *PropIntContext) {}

// ExitPropInt is called when production PropInt is exited.
func (s *BasenmdListener) ExitPropInt(ctx *PropIntContext) {}

// EnterPropFloat is called when production PropFloat is entered.
func (s *BasenmdListener) EnterPropFloat(ctx *PropFloatContext) {}

// ExitPropFloat is called when production PropFloat is exited.
func (s *BasenmdListener) ExitPropFloat(ctx *PropFloatContext) {}
