package structuredheaders

import (
	"errors"
	"strings"
	"unicode"
)

func (s *scanner) scanString() (string, error) {
	var sb strings.Builder

	for s.scanRune() && s.rune != '"' {
		if s.rune == '\\' {
			if !s.scanRune() || s.rune != '"' && s.rune != '\\' {
				return "", errors.New("structuredheaders.scanString: wrong escape seq")
			}

			sb.WriteRune(s.rune)
			continue
		}

		if !unicode.IsPrint(s.rune) {
			return "", errors.New("structuredheaders.scanString: string contains non-printable characters")
		}

		sb.WriteRune(s.rune)
	}

	if s.rune != '"' {
		return "", errors.New("structuredheaders.scanString: unterminated string")
	}

	s.scanRune()

	return sb.String(), nil
}
