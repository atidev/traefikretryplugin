package structuredheaders

import (
	"io"
	"strings"
)

type scanner struct {
	r    *strings.Reader
	rune rune
	err  error
	eof  bool
}

func (s *scanner) scanRune() bool {
	ch, _, err := s.r.ReadRune()

	s.rune = ch
	s.err = err

	if err == io.EOF {
		s.eof = true
	}

	return err == nil
}

func newScanner(r *strings.Reader) *scanner {
	return &scanner{
		r: r,
	}
}
