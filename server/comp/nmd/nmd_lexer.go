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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 17, 117,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5,
	3, 5, 3, 6, 3, 6, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 10, 3, 10, 3,
	11, 3, 11, 3, 12, 3, 12, 3, 12, 7, 12, 64, 10, 12, 12, 12, 14, 12, 67,
	11, 12, 3, 12, 3, 12, 3, 13, 3, 13, 3, 13, 3, 14, 6, 14, 75, 10, 14, 13,
	14, 14, 14, 76, 3, 15, 6, 15, 80, 10, 15, 13, 15, 14, 15, 81, 3, 15, 3,
	15, 7, 15, 86, 10, 15, 12, 15, 14, 15, 89, 11, 15, 3, 15, 3, 15, 6, 15,
	93, 10, 15, 13, 15, 14, 15, 94, 5, 15, 97, 10, 15, 3, 16, 6, 16, 100, 10,
	16, 13, 16, 14, 16, 101, 3, 16, 3, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3, 19,
	3, 19, 3, 19, 7, 19, 113, 10, 19, 12, 19, 14, 19, 116, 11, 19, 3, 65, 2,
	20, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 10, 19, 11, 21, 12,
	23, 13, 25, 2, 27, 14, 29, 15, 31, 16, 33, 2, 35, 2, 37, 17, 3, 2, 5, 5,
	2, 11, 12, 15, 15, 34, 34, 3, 2, 50, 59, 5, 2, 67, 92, 97, 97, 99, 124,
	2, 123, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3,
	2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 17,
	3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2, 2, 2,
	27, 3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2, 2, 2, 2, 37, 3, 2, 2, 2,
	3, 39, 3, 2, 2, 2, 5, 41, 3, 2, 2, 2, 7, 44, 3, 2, 2, 2, 9, 46, 3, 2, 2,
	2, 11, 48, 3, 2, 2, 2, 13, 50, 3, 2, 2, 2, 15, 52, 3, 2, 2, 2, 17, 54,
	3, 2, 2, 2, 19, 56, 3, 2, 2, 2, 21, 58, 3, 2, 2, 2, 23, 60, 3, 2, 2, 2,
	25, 70, 3, 2, 2, 2, 27, 74, 3, 2, 2, 2, 29, 96, 3, 2, 2, 2, 31, 99, 3,
	2, 2, 2, 33, 105, 3, 2, 2, 2, 35, 107, 3, 2, 2, 2, 37, 109, 3, 2, 2, 2,
	39, 40, 7, 61, 2, 2, 40, 4, 3, 2, 2, 2, 41, 42, 7, 47, 2, 2, 42, 43, 7,
	64, 2, 2, 43, 6, 3, 2, 2, 2, 44, 45, 7, 125, 2, 2, 45, 8, 3, 2, 2, 2, 46,
	47, 7, 46, 2, 2, 47, 10, 3, 2, 2, 2, 48, 49, 7, 127, 2, 2, 49, 12, 3, 2,
	2, 2, 50, 51, 7, 93, 2, 2, 51, 14, 3, 2, 2, 2, 52, 53, 7, 95, 2, 2, 53,
	16, 3, 2, 2, 2, 54, 55, 7, 66, 2, 2, 55, 18, 3, 2, 2, 2, 56, 57, 7, 60,
	2, 2, 57, 20, 3, 2, 2, 2, 58, 59, 7, 63, 2, 2, 59, 22, 3, 2, 2, 2, 60,
	65, 7, 41, 2, 2, 61, 64, 5, 25, 13, 2, 62, 64, 11, 2, 2, 2, 63, 61, 3,
	2, 2, 2, 63, 62, 3, 2, 2, 2, 64, 67, 3, 2, 2, 2, 65, 66, 3, 2, 2, 2, 65,
	63, 3, 2, 2, 2, 66, 68, 3, 2, 2, 2, 67, 65, 3, 2, 2, 2, 68, 69, 7, 41,
	2, 2, 69, 24, 3, 2, 2, 2, 70, 71, 7, 94, 2, 2, 71, 72, 7, 41, 2, 2, 72,
	26, 3, 2, 2, 2, 73, 75, 5, 33, 17, 2, 74, 73, 3, 2, 2, 2, 75, 76, 3, 2,
	2, 2, 76, 74, 3, 2, 2, 2, 76, 77, 3, 2, 2, 2, 77, 28, 3, 2, 2, 2, 78, 80,
	5, 33, 17, 2, 79, 78, 3, 2, 2, 2, 80, 81, 3, 2, 2, 2, 81, 79, 3, 2, 2,
	2, 81, 82, 3, 2, 2, 2, 82, 83, 3, 2, 2, 2, 83, 87, 7, 48, 2, 2, 84, 86,
	5, 33, 17, 2, 85, 84, 3, 2, 2, 2, 86, 89, 3, 2, 2, 2, 87, 85, 3, 2, 2,
	2, 87, 88, 3, 2, 2, 2, 88, 97, 3, 2, 2, 2, 89, 87, 3, 2, 2, 2, 90, 92,
	7, 48, 2, 2, 91, 93, 5, 33, 17, 2, 92, 91, 3, 2, 2, 2, 93, 94, 3, 2, 2,
	2, 94, 92, 3, 2, 2, 2, 94, 95, 3, 2, 2, 2, 95, 97, 3, 2, 2, 2, 96, 79,
	3, 2, 2, 2, 96, 90, 3, 2, 2, 2, 97, 30, 3, 2, 2, 2, 98, 100, 9, 2, 2, 2,
	99, 98, 3, 2, 2, 2, 100, 101, 3, 2, 2, 2, 101, 99, 3, 2, 2, 2, 101, 102,
	3, 2, 2, 2, 102, 103, 3, 2, 2, 2, 103, 104, 8, 16, 2, 2, 104, 32, 3, 2,
	2, 2, 105, 106, 9, 3, 2, 2, 106, 34, 3, 2, 2, 2, 107, 108, 9, 4, 2, 2,
	108, 36, 3, 2, 2, 2, 109, 114, 5, 35, 18, 2, 110, 113, 5, 35, 18, 2, 111,
	113, 5, 33, 17, 2, 112, 110, 3, 2, 2, 2, 112, 111, 3, 2, 2, 2, 113, 116,
	3, 2, 2, 2, 114, 112, 3, 2, 2, 2, 114, 115, 3, 2, 2, 2, 115, 38, 3, 2,
	2, 2, 116, 114, 3, 2, 2, 2, 13, 2, 63, 65, 76, 81, 87, 94, 96, 101, 112,
	114, 3, 8, 2, 2,
}

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "';'", "'->'", "'{'", "','", "'}'", "'['", "']'", "'@'", "':'", "'='",
}

var lexerSymbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "QUOTED_STRING", "INT", "FLOAT",
	"WS", "ID",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
	"T__9", "QUOTED_STRING", "ESC", "INT", "FLOAT", "WS", "DIGIT", "LETTER",
	"ID",
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
	nmdLexerQUOTED_STRING = 11
	nmdLexerINT           = 12
	nmdLexerFLOAT         = 13
	nmdLexerWS            = 14
	nmdLexerID            = 15
)
