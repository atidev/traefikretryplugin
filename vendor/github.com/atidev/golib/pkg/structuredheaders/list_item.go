package structuredheaders

import "errors"

type ListItem interface {
	Item() (Item, error)
	InnerList() (InnerList, error)
}

type listItemType int

const (
	itemListItemType listItemType = iota
	listListItemType
)

type listItem struct {
	t listItemType
	i *item
	l *innerList
}

func (l *listItem) Item() (Item, error) {
	if l.t != itemListItemType {
		return nil, errors.New("structuredheaders.Item: list item not an item")
	}

	return l.i, nil
}

func (l *listItem) InnerList() (InnerList, error) {
	if l.t != listListItemType {
		return nil, errors.New("list item not a list")
	}

	return l.l, nil
}

func (s *scanner) scanListItem() (*listItem, error) {
	i, err := s.scanItem()

	if err != nil {
		if errors.Is(err, notAnItem) {
			l, err := s.scanInnerList()
			if err != nil {
				if errors.Is(err, notAnInnerList) {
					return nil, errors.New("list item neither item nor inner list")
				}
				return nil, err
			}

			return &listItem{
				t: listListItemType,
				l: l,
			}, nil
		} else {
			return nil, err
		}
	}

	return &listItem{
		t: itemListItemType,
		i: i,
	}, nil
}
