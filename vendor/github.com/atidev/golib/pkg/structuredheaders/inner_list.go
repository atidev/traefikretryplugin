package structuredheaders

import (
	"errors"
	"fmt"
	"unicode"
)

type InnerList interface {
	Parametrized
	Items() []Item
}

type innerList struct {
	l          []*item
	parameters map[string]*item
}

func (i *innerList) Items() []Item {
	il := make([]Item, 0, len(i.l))

	for _, i := range i.l {
		il = append(il, i)
	}

	return il
}

func (i *innerList) Parameters() map[string]Item {
	p := make(map[string]Item)

	for k, v := range i.parameters {
		p[k] = v
	}

	return p
}

var notAnInnerList = errors.New("not an inner list")

func (s *scanner) scanInnerList() (*innerList, error) {
	if s.rune != '(' {
		return nil, fmt.Errorf("structuredheaders.scanInnerList: %w", notAnInnerList)
	}

	l := make([]*item, 0)

	var err error
	var i *item
	for i, err = s.scanItem(); unicode.IsSpace(s.rune) && err == nil; i, err = s.scanItem() {
		l = append(l, i)
	}
	l = append(l, i)

	if err != nil {
		return nil, fmt.Errorf("structuredheaders.scanInnerList: %w", err)
	}

	if s.rune != ')' {
		return nil, errors.New("structuredheaders.scanInnerList: unterminated inner list")
	}

	var p map[string]*item

	if s.scanRune() {
		p, err = s.scanParameters()
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.scanInnerList: %w", err)
		}
	}

	return &innerList{
		l:          l,
		parameters: p,
	}, nil
}
