package structuredheaders

import (
	"fmt"
	"net/http"
	"strings"
)

type StructuredHeader interface {
	Item(key string) (Item, error)
	List(key string) ([]ListItem, error)
	Dictionary(key string) (map[string]ListItem, error)
}

func NewStructuredHeader(header http.Header) StructuredHeader {
	return &structuredHeader{h: header}
}

type structuredHeader struct {
	h http.Header
}

func (s *structuredHeader) Item(key string) (Item, error) {
	v := s.h.Get(key)

	return newScanner(strings.NewReader(v)).scanItem()
}

func (s *structuredHeader) List(key string) ([]ListItem, error) {
	values := s.h.Values(key)

	tl := make([]ListItem, 0)

	for _, v := range values {
		vl, err := newScanner(strings.NewReader(v)).scanList()
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.List: %w", err)
		}

		tl = append(tl, vl...)
	}

	return tl, nil
}

func (s *structuredHeader) Dictionary(key string) (map[string]ListItem, error) {
	values := s.h.Values(key)

	td := make(map[string]ListItem)

	for _, v := range values {
		vd, err := newScanner(strings.NewReader(v)).scanDictionary()
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.Dictionary: %w", err)
		}

		for k, vdv := range vd {
			td[k] = vdv
		}
	}

	return td, nil
}
