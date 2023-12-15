package structuredheaders

import (
	"errors"
	"fmt"
)

func (s *scanner) scanList() ([]ListItem, error) {
	l := make([]ListItem, 0)

	for s.scanRune() {
		i, err := s.scanListItem()
		if err != nil {
			return nil, fmt.Errorf("structuredheaders.scanList: %w", err)
		}

		if !s.eof && s.rune != ',' {
			return nil, errors.New("structuredheaders.scanList: wrong list delimiter")
		}

		l = append(l, i)
	}

	return l, nil
}
