package structuredheaders

import (
	"errors"
	"strings"
	"unicode"
)

type Parametrized interface {
	Parameters() map[string]Item
}

func (s *scanner) scanParameters() (map[string]*item, error) {
	p := make(map[string]*item)

	if s.rune != ';' {
		return p, nil
	}

	for s.scanRune() {
		if unicode.IsSpace(s.rune) {
			continue
		}

		key, err := s.scanKey()

		if err != nil {
			return nil, err
		}

		p[key], err = s.scanParam()
		if err != nil {
			return nil, err
		}

		if s.rune != ';' {
			break
		}
	}

	return p, nil
}

func (s *scanner) scanKey() (string, error) {
	r := s.rune

	var sb strings.Builder

	sb.WriteRune(r)

	if !unicode.Is(unicode.Latin, r) && r != '*' {
		return "", errors.New("structuredheaders.scanKey: can't scan key")
	}

	for s.scanRune() {
		r = s.rune

		if !isKeyRune(r) {
			break
		}

		sb.WriteRune(r)
	}

	return sb.String(), nil
}

func isKeyRune(r rune) bool {
	return unicode.Is(unicode.Latin, r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' || r == '*'
}

func (s *scanner) scanParam() (*item, error) {
	if s.rune != '=' {
		return &item{
			t:       itemTypeBoolean,
			boolean: true,
		}, nil
	}

	return s.scanBareItem()
}
