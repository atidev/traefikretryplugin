package structuredheaders

import "errors"

func (s *scanner) scanBoolean() (bool, error) {
	s.scanRune()

	var b bool
	switch s.rune {
	case '0':
		b = false
	case '1':
		b = true
	default:
		return false, errors.New("structuredheaders.scanBoolean: not a boolean character")
	}

	s.scanRune()

	return b, nil
}
