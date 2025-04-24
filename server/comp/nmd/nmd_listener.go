// Code generated from nmd.g4 by ANTLR 4.9.2. DO NOT EDIT.

package nmd // nmd
import "github.com/antlr/antlr4/runtime/Go/antlr"

// nmdListener is a complete listener for a parse tree produced by nmdParser.
type nmdListener interface {
	antlr.ParseTreeListener

	// EnterGraph is called when entering the graph production.
	EnterGraph(c *GraphContext)

	// EnterStmt_list is called when entering the stmt_list production.
	EnterStmt_list(c *Stmt_listContext)

	// EnterStmt is called when entering the stmt production.
	EnterStmt(c *StmtContext)

	// EnterLink_stmt is called when entering the link_stmt production.
	EnterLink_stmt(c *Link_stmtContext)

	// EnterCall_stmt is called when entering the call_stmt production.
	EnterCall_stmt(c *Call_stmtContext)

	// EnterCast_stmt is called when entering the cast_stmt production.
	EnterCast_stmt(c *Cast_stmtContext)

	// EnterSink_stmt is called when entering the sink_stmt production.
	EnterSink_stmt(c *Sink_stmtContext)

	// EnterEndpoint is called when entering the endpoint production.
	EnterEndpoint(c *EndpointContext)

	// EnterNode_def is called when entering the node_def production.
	EnterNode_def(c *Node_defContext)

	// EnterNode_id is called when entering the node_id production.
	EnterNode_id(c *Node_idContext)

	// EnterNode_prop is called when entering the node_prop production.
	EnterNode_prop(c *Node_propContext)

	// EnterMsg_type_list is called when entering the msg_type_list production.
	EnterMsg_type_list(c *Msg_type_listContext)

	// EnterLink_operator is called when entering the link_operator production.
	EnterLink_operator(c *Link_operatorContext)

	// EnterPropQuoteString is called when entering the PropQuoteString production.
	EnterPropQuoteString(c *PropQuoteStringContext)

	// EnterPropId is called when entering the PropId production.
	EnterPropId(c *PropIdContext)

	// EnterPropInt is called when entering the PropInt production.
	EnterPropInt(c *PropIntContext)

	// EnterPropFloat is called when entering the PropFloat production.
	EnterPropFloat(c *PropFloatContext)

	// ExitGraph is called when exiting the graph production.
	ExitGraph(c *GraphContext)

	// ExitStmt_list is called when exiting the stmt_list production.
	ExitStmt_list(c *Stmt_listContext)

	// ExitStmt is called when exiting the stmt production.
	ExitStmt(c *StmtContext)

	// ExitLink_stmt is called when exiting the link_stmt production.
	ExitLink_stmt(c *Link_stmtContext)

	// ExitCall_stmt is called when exiting the call_stmt production.
	ExitCall_stmt(c *Call_stmtContext)

	// ExitCast_stmt is called when exiting the cast_stmt production.
	ExitCast_stmt(c *Cast_stmtContext)

	// ExitSink_stmt is called when exiting the sink_stmt production.
	ExitSink_stmt(c *Sink_stmtContext)

	// ExitEndpoint is called when exiting the endpoint production.
	ExitEndpoint(c *EndpointContext)

	// ExitNode_def is called when exiting the node_def production.
	ExitNode_def(c *Node_defContext)

	// ExitNode_id is called when exiting the node_id production.
	ExitNode_id(c *Node_idContext)

	// ExitNode_prop is called when exiting the node_prop production.
	ExitNode_prop(c *Node_propContext)

	// ExitMsg_type_list is called when exiting the msg_type_list production.
	ExitMsg_type_list(c *Msg_type_listContext)

	// ExitLink_operator is called when exiting the link_operator production.
	ExitLink_operator(c *Link_operatorContext)

	// ExitPropQuoteString is called when exiting the PropQuoteString production.
	ExitPropQuoteString(c *PropQuoteStringContext)

	// ExitPropId is called when exiting the PropId production.
	ExitPropId(c *PropIdContext)

	// ExitPropInt is called when exiting the PropInt production.
	ExitPropInt(c *PropIntContext)

	// ExitPropFloat is called when exiting the PropFloat production.
	ExitPropFloat(c *PropFloatContext)
}
