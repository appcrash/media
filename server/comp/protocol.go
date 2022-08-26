package comp

import (
	"fmt"
	"strings"
)

// CommandInitiator send CtrlMessage to nodes and get replies from them, and simply use string as rpc protocol.
// this file provides utility functions to build and parse strings based on the protocol.

//var regRemoveSpace = regexp.MustCompile("\\s+")
const (
	stateNormal = iota
	stateInQuoteBackslash
	stateInQuote
	stateNone
)

func With(args ...string) []string {
	return args
}

// WithString parse a string into substrings separated by SPACE or TAB, double quote is required if substring contains
// space, and backslash is requited for escaped double quote, for example:
// cmd arg_a "arg with space", "arg with \" escaped \" quote"
func WithString(args string) (result []string, err error) {
	i, start := 0, 0
	n := len(args)
	s := stateNone
	for i < n {
		switch args[i] {
		case ' ', '\t':
			if s == stateNormal {
				slice := args[start:i]
				result = append(result, slice)
				s = stateNone
			} else if s == stateInQuoteBackslash {

			}
		case '\\':
			if s == stateNone {
				s = stateNormal
				start = i
			} else if s == stateInQuote {
				s = stateInQuoteBackslash
			}
		case '"':
			if s == stateNone {
				// start quote
				s = stateInQuote
				start = i
			} else if s == stateInQuote {
				// end quote
				slice := args[start+1 : i] // remove start quote char
				slice = strings.Replace(slice, "\\\"", "\"", -1)
				result = append(result, slice)
				s = stateNone
			} else if s == stateInQuoteBackslash {
				// in form of "...\"...", just forward
				s = stateInQuote
			} else {
				err = fmt.Errorf("illegal double quote at index %v", i)
				return
			}
		default:
			// any character other than above
			if s == stateNone {
				start = i
				s = stateNormal
			}
		}
		i++
	}

	switch s {
	case stateNormal:
		slice := args[start:i]
		result = append(result, slice)
	case stateInQuote, stateInQuoteBackslash:
		err = fmt.Errorf("quote is not closed at the end")
	}
	return
}

// WithConnect dynamically ask the callee to set data pipe endpoint
func WithConnect(toSession, toName string) []string {
	return With("conn", toSession, toName)
}

// WithOk return normal reply
func WithOk(args ...string) (r []string) {
	r = append(r, "ok")
	r = append(r, args...)
	return
}

// WithError return abnormal reply
func WithError(args ...string) (r []string) {
	r = append(r, "err")
	r = append(r, args...)
	return
}
