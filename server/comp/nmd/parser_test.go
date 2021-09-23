package nmd_test

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/appcrash/media/server/comp/nmd"
	"testing"
)

func TestNmd(t *testing.T) {
	testStr := `[a@mya:ofa aa='sfa\\'oiow' aa=''] -> {[b c='http://iow.com' cc=.89] ,[c]} -> {[d], [e@ff:io ew=23]}`
	input := antlr.NewInputStream(testStr)
	lexer := nmd.NewnmdLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := nmd.NewnmdParser(stream)
	listener := nmd.NewListener("test_session")
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Graph())
	for _, n := range listener.NodeDefs {
		for _, to := range n.Deps {
			t.Logf("%v -> %v", n, to)
		}
	}
}
