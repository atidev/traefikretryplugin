package structuredheaders

import (
	"errors"
	"fmt"
	"unicode"
)

func (s *scanner) scanDictionary() (map[string]ListItem, error) {
	d := make(map[string]ListItem)

	for s.scanRune() {
		if s.rune == ' ' {
			for s.scanRune() {
				if !unicode.IsSpace(s.rune) {
					break
				}
			}
		}

		k, i, err := s.scanDictionaryItem()
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.scanDictionary: %w", err)
		}

		if !s.eof && s.rune != ',' {
			return nil, errors.New("structuredheaders.scanList: wrong dictionary delimiter")
		}

		d[k] = i
	}

	return d, nil
}

func (s *scanner) scanDictionaryItem() (string, ListItem, error) {
	k, err := s.scanKey()
	if err != nil {
		return "", nil, fmt.Errorf("structuredheaders.scanDictionaryItem: %w", err)
	}

	if s.rune == ',' || s.eof {
		return k, &listItem{
			t: itemListItemType,
			i: &item{
				t:       itemTypeBoolean,
				boolean: true,
			},
		}, nil
	} else if s.rune == '=' {
		li, err := s.scanListItem()
		if err != nil {
			return "", nil, fmt.Errorf("structuredheaders.scanDictionaryItem: %w", err)
		}

		return k, li, nil
	}

	return "", nil, errors.New("structuredheaders.scanDictionaryItem: not an item")
}
