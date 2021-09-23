// Code generated from nmd.g4 by ANTLR 4.9.2. DO NOT EDIT.

package nmd

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 19, 129,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 3, 2, 3, 2, 3, 3,
	3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6,
	3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 10, 3, 10, 3, 11, 3, 11, 3, 12,
	3, 12, 3, 13, 3, 13, 3, 14, 3, 14, 3, 14, 7, 14, 76, 10, 14, 12, 14, 14,
	14, 79, 11, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15, 3, 16, 6, 16, 87, 10,
	16, 13, 16, 14, 16, 88, 3, 17, 6, 17, 92, 10, 17, 13, 17, 14, 17, 93, 3,
	17, 3, 17, 7, 17, 98, 10, 17, 12, 17, 14, 17, 101, 11, 17, 3, 17, 3, 17,
	6, 17, 105, 10, 17, 13, 17, 14, 17, 106, 5, 17, 109, 10, 17, 3, 18, 6,
	18, 112, 10, 18, 13, 18, 14, 18, 113, 3, 18, 3, 18, 3, 19, 3, 19, 3, 20,
	3, 20, 3, 21, 3, 21, 3, 21, 7, 21, 125, 10, 21, 12, 21, 14, 21, 128, 11,
	21, 3, 77, 2, 22, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 10,
	19, 11, 21, 12, 23, 13, 25, 14, 27, 15, 29, 2, 31, 16, 33, 17, 35, 18,
	37, 2, 39, 2, 41, 19, 3, 2, 5, 5, 2, 11, 12, 15, 15, 34, 34, 3, 2, 50,
	59, 5, 2, 67, 92, 97, 97, 99, 124, 2, 135, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2,
	2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3,
	2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21,
	3, 2, 2, 2, 2, 23, 3, 2, 2, 2, 2, 25, 3, 2, 2, 2, 2, 27, 3, 2, 2, 2, 2,
	31, 3, 2, 2, 2, 2, 33, 3, 2, 2, 2, 2, 35, 3, 2, 2, 2, 2, 41, 3, 2, 2, 2,
	3, 43, 3, 2, 2, 2, 5, 45, 3, 2, 2, 2, 7, 48, 3, 2, 2, 2, 9, 52, 3, 2, 2,
	2, 11, 56, 3, 2, 2, 2, 13, 58, 3, 2, 2, 2, 15, 60, 3, 2, 2, 2, 17, 62,
	3, 2, 2, 2, 19, 64, 3, 2, 2, 2, 21, 66, 3, 2, 2, 2, 23, 68, 3, 2, 2, 2,
	25, 70, 3, 2, 2, 2, 27, 72, 3, 2, 2, 2, 29, 82, 3, 2, 2, 2, 31, 86, 3,
	2, 2, 2, 33, 108, 3, 2, 2, 2, 35, 111, 3, 2, 2, 2, 37, 117, 3, 2, 2, 2,
	39, 119, 3, 2, 2, 2, 41, 121, 3, 2, 2, 2, 43, 44, 7, 61, 2, 2, 44, 4, 3,
	2, 2, 2, 45, 46, 7, 47, 2, 2, 46, 47, 7, 64, 2, 2, 47, 6, 3, 2, 2, 2, 48,
	49, 7, 62, 2, 2, 49, 50, 7, 47, 2, 2, 50, 51, 7, 64, 2, 2, 51, 8, 3, 2,
	2, 2, 52, 53, 7, 62, 2, 2, 53, 54, 7, 47, 2, 2, 54, 55, 7, 47, 2, 2, 55,
	10, 3, 2, 2, 2, 56, 57, 7, 125, 2, 2, 57, 12, 3, 2, 2, 2, 58, 59, 7, 46,
	2, 2, 59, 14, 3, 2, 2, 2, 60, 61, 7, 127, 2, 2, 61, 16, 3, 2, 2, 2, 62,
	63, 7, 93, 2, 2, 63, 18, 3, 2, 2, 2, 64, 65, 7, 95, 2, 2, 65, 20, 3, 2,
	2, 2, 66, 67, 7, 66, 2, 2, 67, 22, 3, 2, 2, 2, 68, 69, 7, 60, 2, 2, 69,
	24, 3, 2, 2, 2, 70, 71, 7, 63, 2, 2, 71, 26, 3, 2, 2, 2, 72, 77, 7, 41,
	2, 2, 73, 76, 5, 29, 15, 2, 74, 76, 11, 2, 2, 2, 75, 73, 3, 2, 2, 2, 75,
	74, 3, 2, 2, 2, 76, 79, 3, 2, 2, 2, 77, 78, 3, 2, 2, 2, 77, 75, 3, 2, 2,
	2, 78, 80, 3, 2, 2, 2, 79, 77, 3, 2, 2, 2, 80, 81, 7, 41, 2, 2, 81, 28,
	3, 2, 2, 2, 82, 83, 7, 94, 2, 2, 83, 84, 7, 41, 2, 2, 84, 30, 3, 2, 2,
	2, 85, 87, 5, 37, 19, 2, 86, 85, 3, 2, 2, 2, 87, 88, 3, 2, 2, 2, 88, 86,
	3, 2, 2, 2, 88, 89, 3, 2, 2, 2, 89, 32, 3, 2, 2, 2, 90, 92, 5, 37, 19,
	2, 91, 90, 3, 2, 2, 2, 92, 93, 3, 2, 2, 2, 93, 91, 3, 2, 2, 2, 93, 94,
	3, 2, 2, 2, 94, 95, 3, 2, 2, 2, 95, 99, 7, 48, 2, 2, 96, 98, 5, 37, 19,
	2, 97, 96, 3, 2, 2, 2, 98, 101, 3, 2, 2, 2, 99, 97, 3, 2, 2, 2, 99, 100,
	3, 2, 2, 2, 100, 109, 3, 2, 2, 2, 101, 99, 3, 2, 2, 2, 102, 104, 7, 48,
	2, 2, 103, 105, 5, 37, 19, 2, 104, 103, 3, 2, 2, 2, 105, 106, 3, 2, 2,
	2, 106, 104, 3, 2, 2, 2, 106, 107, 3, 2, 2, 2, 107, 109, 3, 2, 2, 2, 108,
	91, 3, 2, 2, 2, 108, 102, 3, 2, 2, 2, 109, 34, 3, 2, 2, 2, 110, 112, 9,
	2, 2, 2, 111, 110, 3, 2, 2, 2, 112, 113, 3, 2, 2, 2, 113, 111, 3, 2, 2,
	2, 113, 114, 3, 2, 2, 2, 114, 115, 3, 2, 2, 2, 115, 116, 8, 18, 2, 2, 116,
	36, 3, 2, 2, 2, 117, 118, 9, 3, 2, 2, 118, 38, 3, 2, 2, 2, 119, 120, 9,
	4, 2, 2, 120, 40, 3, 2, 2, 2, 121, 126, 5, 39, 20, 2, 122, 125, 5, 39,
	20, 2, 123, 125, 5, 37, 19, 2, 124, 122, 3, 2, 2, 2, 124, 123, 3, 2, 2,
	2, 125, 128, 3, 2, 2, 2, 126, 124, 3, 2, 2, 2, 126, 127, 3, 2, 2, 2, 127,
	42, 3, 2, 2, 2, 128, 126, 3, 2, 2, 2, 13, 2, 75, 77, 88, 93, 99, 106, 108,
	113, 124, 126, 3, 8, 2, 2,
}

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "';'", "'->'", "'<->'", "'<--'", "'{'", "','", "'}'", "'['", "']'",
	"'@'", "':'", "'='",
}

var lexerSymbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "QUOTED_STRING", "INT",
	"FLOAT", "WS", "ID",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
	"T__9", "T__10", "T__11", "QUOTED_STRING", "ESC", "INT", "FLOAT", "WS",
	"DIGIT", "LETTER", "ID",
}

type nmdLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

// NewnmdLexer produces a new lexer instance for the optional input antlr.CharStream.
//
// The *nmdLexer instance produced may be reused by calling the SetInputStream method.
// The initial lexer configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewnmdLexer(input antlr.CharStream) *nmdLexer {
	l := new(nmdLexer)
	lexerDeserializer := antlr.NewATNDeserializer(nil)
	lexerAtn := lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)
	lexerDecisionToDFA := make([]*antlr.DFA, len(lexerAtn.DecisionToState))
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "nmd.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// nmdLexer tokens.
const (
	nmdLexerT__0          = 1
	nmdLexerT__1          = 2
	nmdLexerT__2          = 3
	nmdLexerT__3          = 4
	nmdLexerT__4          = 5
	nmdLexerT__5          = 6
	nmdLexerT__6          = 7
	nmdLexerT__7          = 8
	nmdLexerT__8          = 9
	nmdLexerT__9          = 10
	nmdLexerT__10         = 11
	nmdLexerT__11         = 12
	nmdLexerQUOTED_STRING = 13
	nmdLexerINT           = 14
	nmdLexerFLOAT         = 15
	nmdLexerWS            = 16
	nmdLexerID            = 17
)
