package structuredheaders

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type number struct {
	t       numberType
	integer int
	float   float64
}

type Number interface {
	Integer() (int, error)
	Float() (float64, error)
}

func (n *number) Integer() (int, error) {
	if n.t != integerNumberType {
		return 0, errors.New("structuredheaders.Integer: not an integer")
	}

	return n.integer, nil
}

func (n *number) Float() (float64, error) {
	if n.t == integerNumberType {
		return float64(n.integer), nil
	}

	return n.float, nil
}

type numberType int

const (
	integerNumberType numberType = iota
	floatNumberType
)

func (s *scanner) scanNumber() (*number, error) {
	numberType := integerNumberType

	var sb strings.Builder

	r := s.rune

	sb.WriteRune(r)

	for s.scanRune() && isNumberRune(s.rune) {
		if s.rune == '.' {
			if numberType == integerNumberType {
				numberType = floatNumberType
			} else {
				break
			}
		}

		sb.WriteRune(r)
	}

	ns := sb.String()
	n := &number{
		t: numberType,
	}

	switch numberType {
	case integerNumberType:
		i, err := strconv.Atoi(ns)
		if err != nil {
			return nil, err
		}

		n.integer = i

	case floatNumberType:
		f, err := strconv.ParseFloat(ns, 32)
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.scanNumber: %w", err)
		}

		n.float = f
	}

	return n, nil
}

func isNumberRune(r rune) bool {
	return unicode.IsDigit(r) || r == '.'
}
