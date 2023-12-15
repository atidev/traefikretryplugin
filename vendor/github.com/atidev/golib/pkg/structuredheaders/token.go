package structuredheaders

import (
	"strings"
	"unicode"
)

func (s *scanner) scanToken() string {
	var sb strings.Builder

	sb.WriteRune(s.rune)

	for s.scanRune() && isTokenRune(s.rune) {
		sb.WriteRune(s.rune)
	}

	return sb.String()
}

func isTokenRune(r rune) bool {
	return r == '!' ||
		r == '#' ||
		r == '$' ||
		r == '%' ||
		r == '&' ||
		r == '\'' ||
		r == '*' ||
		r == '+' ||
		r == '-' ||
		r == '.' ||
		r == '^' ||
		r == '_' ||
		r == '`' ||
		r == '|' ||
		r == '~' ||
		r == ':' ||
		r == '/' ||
		unicode.IsDigit(r) ||
		unicode.Is(unicode.Latin, r)
}
