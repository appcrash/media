package utils

import (
	"regexp"
	"strings"
	"unicode"
)

const regSnakeToCamelCasePattern = `_+[a-z]`

var regSnakeToCamelCase = regexp.MustCompile(regSnakeToCamelCasePattern)

// SnakeToCamelCase converts string in form of foo_bar or foo__bar ... into fooBar
func SnakeToCamelCase(str string) string {
	return regSnakeToCamelCase.ReplaceAllStringFunc(str, func(match string) string {
		last := string(match[len(match)-1])
		return strings.ToUpper(last)
	})
}

// CamelCaseToSnake converts string in form of fooBar into foo_bar
func CamelCaseToSnake(str string) string {
	var sb strings.Builder
	var hasSnake = true // avoid inserting '_' at head
	for _, c := range str {
		if c == '_' {
			if !hasSnake {
				sb.WriteByte('_')
			}
			hasSnake = true
		} else {
			if unicode.IsUpper(c) {
				if !hasSnake {
					sb.WriteByte('_')
				}
				sb.WriteRune(unicode.ToLower(c))
			} else {
				sb.WriteRune(c)
			}
			hasSnake = false
		}
	}
	return sb.String()
}
