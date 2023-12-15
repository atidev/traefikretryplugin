package structuredheaders

import (
	"errors"
	"fmt"
	"unicode"
)

type Item interface {
	Parametrized
	Number() (Number, error)
	Str() (string, error)
	Token() (string, error)
	Binary() ([]byte, error)
	Boolean() (bool, error)
}

type itemType int

const (
	itemTypeNumber itemType = iota
	itemTypeString
	itemTypeToken
	itemTypeBinary
	itemTypeBoolean
)

type item struct {
	t          itemType
	number     *number
	boolean    bool
	string     string
	token      string
	binary     []byte
	parameters map[string]*item
}

func (i *item) Number() (Number, error) {
	if i.t != itemTypeNumber {
		return nil, errors.New("structuredheaders.Number: not a number")
	}

	return i.number, nil
}

func (i *item) Str() (string, error) {
	if i.t != itemTypeString {
		return "", errors.New("structuredheaders.Str: not a string")
	}

	return i.string, nil
}

func (i *item) Token() (string, error) {
	if i.t != itemTypeToken {
		return "", errors.New("structuredheaders.Token: not a token")
	}

	return i.token, nil
}

func (i *item) Binary() ([]byte, error) {
	if i.t != itemTypeBinary {
		return nil, errors.New("structuredheaders.Binary: not a binary")
	}

	return i.binary, nil
}

func (i *item) Boolean() (bool, error) {
	if i.t != itemTypeBoolean {
		return false, errors.New("structuredheaders.Boolean: not a boolean")
	}

	return i.boolean, nil
}

func (i *item) Parameters() map[string]Item {
	p := make(map[string]Item)

	for k, v := range i.parameters {
		p[k] = v
	}

	return p
}

func (s *scanner) scanItem() (*item, error) {
	i, err := s.scanBareItem()
	if err != nil {
		return nil, fmt.Errorf("structuredheaders.scanItem: %w", err)
	}

	p, err := s.scanParameters()
	if err != nil {
		return nil, fmt.Errorf("structuredheaders.scanItem: %w", err)
	}

	i.parameters = p

	return i, nil
}

var notAnItem = errors.New("not an item")

func (s *scanner) scanBareItem() (*item, error) {
	var itemType itemType

	for s.scanRune() {
		if !unicode.IsSpace(s.rune) {
			break
		}
	}

	if err := s.err; err != nil {
		return nil, fmt.Errorf("structuredheaders.scanBareItem: %w", err)
	}

	if unicode.IsDigit(s.rune) || s.rune == '-' {
		itemType = itemTypeNumber
	} else if s.rune == '"' {
		itemType = itemTypeString
	} else if unicode.Is(unicode.Latin, s.rune) || s.rune == '*' {
		itemType = itemTypeToken
	} else if s.rune == ':' {
		itemType = itemTypeBinary
	} else if s.rune == '?' {
		itemType = itemTypeBoolean
	} else {
		return nil, fmt.Errorf("structuredheaders.scanBareItem: %w", notAnItem)
	}

	i := &item{
		t: itemType,
	}

	switch itemType {
	case itemTypeNumber:
		n, err := s.scanNumber()
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.scanBareItem: %w", err)
		}

		i.number = n
	case itemTypeString:
		str, err := s.scanString()
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.scanBareItem: %w", err)
		}

		i.string = str
	case itemTypeToken:
		i.token = s.scanToken()
	case itemTypeBinary:
		b, err := s.scanBinary()

		if err != nil {
			return nil, fmt.Errorf("structuredheaders.scanBareItem: %w", err)
		}

		i.binary = b
	case itemTypeBoolean:
		b, err := s.scanBoolean()
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.scanBareItem: %w", err)
		}

		i.boolean = b
	}

	return i, nil
}
