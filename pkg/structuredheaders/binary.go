package structuredheaders

import (
	"encoding/base64"
	"errors"
	"strings"
	"unicode"
)

func (s *scanner) scanBinary() ([]byte, error) {
	var sb strings.Builder

	for s.scanRune() && isBase64Rune(s.rune) {
		sb.WriteRune(s.rune)
	}

	if s.rune != ':' {
		return nil, errors.New("structuredheaders.scanBinary: unterminated binary")
	}

	str := []byte(sb.String())

	l := base64.StdEncoding.EncodedLen(len(str))
	b := make([]byte, l)

	base64.StdEncoding.Encode(b, str)

	return b, nil
}

func isBase64Rune(r rune) bool {
	return unicode.Is(unicode.Latin, r) || unicode.IsDigit(r) || r == '+' || r == '/' || r == '='
}
