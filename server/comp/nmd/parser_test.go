package nmd_test

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/appcrash/media/server/comp/nmd"
	"testing"
)

func TestNodeDef_Link(t *testing.T) {
	testStr := `[a@mya:ofa aa='\'oiow' aa=''] <msg1,msg2> {[b c='http://iow.com' cc=.89] ,[c]}
		<msg3,msg4> {[d], [e@ff:io ew=23]} -> [f]`
	input := antlr.NewInputStream(testStr)
	lexer := nmd.NewnmdLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := nmd.NewnmdParser(stream)
	listener := nmd.NewListener("test_session")
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Graph())
	for _, n := range listener.NodeDefs {
		switch n.Name {
		case "a":
			if n.Type != "ofa" && n.Scope != "mya" {
				t.Fatal("parse node def failed")
			}
			if n.Deps == nil {
				t.Fatal("parse link operator failed")
			}
			preferOffer := n.Deps[0].PreferOffer
			if preferOffer[0] != "msg1" || preferOffer[1] != "msg2" {
				t.Fatal("parse msg type list failed")
			}
		case "b":
			for _, to := range n.Deps {
				if to.LinkTo.Name != "e" && to.LinkTo.Name != "d" {
					t.Fatal("parse link statement failed")
				}
			}
		case "c":
			if n.Deps == nil {
				t.Fatal("parse link operator failed")
			}
			preferOffer := n.Deps[0].PreferOffer
			if preferOffer[0] != "msg3" || preferOffer[1] != "msg4" {
				t.Fatal("parse msg type list failed")
			}
		}

	}
}

func TestCallCast(t *testing.T) {
	testStr := `[a] <-> 'call cmd';
                [b@test:type_b] <-- 'cast cmd'`
	input := antlr.NewInputStream(testStr)
	lexer := nmd.NewnmdLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := nmd.NewnmdParser(stream)
	listener := nmd.NewListener("test_session")
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Graph())

	callDef := listener.CallDefs[0]
	castDef := listener.CastDefs[0]
	if callDef.Node.Name != "a" || callDef.Cmd != "call cmd" {
		t.Fatal("parse call statement failed")
	}
	if castDef.Node.Name != "b" || castDef.Node.Type != "type_b" || castDef.Cmd != "cast cmd" {
		t.Fatal("parse cast statement failed")
	}
}

func TestSink(t *testing.T) {
	testStr := `[a] <-> 'call cmd';
                <-chan mychannel`
	input := antlr.NewInputStream(testStr)
	lexer := nmd.NewnmdLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := nmd.NewnmdParser(stream)
	listener := nmd.NewListener("test_session")
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Graph())

	sinkDef := listener.SinkDefs[0]
	if sinkDef.NodeName != "mychannel" {
		t.Fatal("parse sink statement failed")
	}
}
