package intervals

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/atidev/golib/pkg/set"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Interval interface {
	Includes(num int) bool
	String() string
}

func NewInterval(str string) (Interval, error) {
	s := bufio.NewScanner(strings.NewReader(str))

	s.Split(splitTokens)

	return readInterval(s)
}

type intervalBound struct {
	int
	strict bool
}
type boundedInterval struct {
	from intervalBound
	to   intervalBound
}

type interval struct {
	b []*boundedInterval
	v set.Set[int]
}

func (iv *interval) String() string {
	var sb strings.Builder

	for i, b := range iv.b {
		if i > 0 {
			sb.WriteRune(' ')
		}

		sb.WriteString(b.String())
	}

	for i, value := range iv.v.Values() {
		if i > 0 || len(iv.b) > 0 {
			sb.WriteRune(' ')
		}

		sb.WriteString(strconv.Itoa(value))
	}

	return sb.String()
}

func (iv *interval) Includes(num int) bool {
	if iv.v.Includes(num) {
		return true
	}

	for _, br := range iv.b {
		if br.Includes(num) {
			return true
		}
	}

	return false
}

func (bi *boundedInterval) Includes(num int) bool {
	from := bi.from.int < num || !bi.from.strict && bi.from.int == num

	to := bi.to.int > num || !bi.to.strict && bi.to.int == num

	return from && to
}

func (bi *boundedInterval) String() string {
	var sb strings.Builder

	if bi.from.strict {
		sb.WriteRune('(')
	} else {
		sb.WriteRune('[')
	}

	sb.WriteString(strconv.Itoa(bi.from.int))
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(bi.to.int))

	if bi.to.strict {
		sb.WriteRune(')')
	} else {
		sb.WriteRune(']')
	}

	return sb.String()
}

func (bi *boundedInterval) intersect(b *boundedInterval) *boundedInterval {
	if bi.to.int > b.from.int && bi.from.int < b.to.int {
		return &boundedInterval{
			from: b.from,
			to:   bi.to,
		}
	}
	if !bi.to.strict && bi.to.int == b.from.int {
		return &boundedInterval{
			from: b.from,
			to:   bi.to,
		}
	}

	return nil
}

func intersectBounded(a *boundedInterval, b *boundedInterval) *boundedInterval {
	if mr := a.intersect(b); mr != nil {
		return mr
	}

	return b.intersect(a)
}

func readInterval(s *bufio.Scanner) (Interval, error) {
	rgs := make([]*boundedInterval, 0)

	values := set.NewSet[int]()

	for s.Scan() {
		token := s.Text()
		if isIntervalBeginString(token) {
			r, err := scanBounded(s, token)
			if err != nil {
				return nil, fmt.Errorf("intervals.readInterval: %w", err)
			}

			rgs = append(rgs, r)
		} else {
			i, err := strconv.Atoi(token)
			if err != nil {
				return nil, fmt.Errorf("intervals.readInterval: %w", err)
			}

			values.Add(i)
		}
	}

	if err := s.Err(); err != nil && err != io.EOF {
		return nil, fmt.Errorf("intervals.readInterval: %w", err)
	}

	merged := mergeBounded(rgs)

	return &interval{
		b: merged,
		v: mergeValues(values, merged),
	}, nil
}

func mergeValues(values set.Set[int], intervals []*boundedInterval) set.Set[int] {
	included := set.NewSet[int]()

	for _, i := range values.Values() {
		for _, r := range intervals {
			if r.Includes(i) {
				included.Add(i)
			}
		}
	}

	v := set.NewSet[int]()

	for _, i := range values.Values() {
		if !included.Includes(i) {
			v.Add(i)
		}
	}

	return v
}

func mergeBounded(rgs []*boundedInterval) []*boundedInterval {
	if len(rgs) <= 1 {
		return rgs
	}

	merged := make([]*boundedInterval, 0, len(rgs))

	idx := set.NewSet[int]()

	for i, r := range rgs {
		if idx.Includes(i) {
			continue
		}

		start := i + 1
		for j, rr := range rgs[start:] {
			if mr := intersectBounded(r, rr); mr != nil {
				idx.Add(i)
				idx.Add(start + j)

				merged = append(merged, mr)
			}
		}

		if !idx.Includes(i) {
			merged = append(merged, r)
		}
	}

	if idx.Len() > 0 {
		return mergeBounded(merged)
	} else {
		return merged
	}
}

func scanBounded(s *bufio.Scanner, token string) (*boundedInterval, error) {
	i, err := scanInt(s)
	if err != nil {
		return nil, fmt.Errorf("intervals.scanBounded: %w", err)
	}

	from := intervalBound{
		int:    i,
		strict: isStrictIntervalString(token),
	}

	i, err = scanInt(s)
	if err != nil {
		return nil, fmt.Errorf("intervals.scanBounded: %w", err)
	}

	if !s.Scan() {
		if err = s.Err(); err != nil {
			return nil, fmt.Errorf("intervals.scanBounded: %w", err)
		}
	}

	token = s.Text()
	if !isIntervalEndString(token) {
		return nil, errors.New("intervals.scanBounded: unclosed interval")
	}

	to := intervalBound{
		int:    i,
		strict: isStrictIntervalString(token),
	}

	bi := &boundedInterval{
		from: from,
		to:   to,
	}

	if from.int > to.int || from.strict && to.strict && from.int == to.int {
		return nil, fmt.Errorf("scanBounded: interval `%s` is empty", bi.String())
	}

	return bi, nil
}

func scanInt(s *bufio.Scanner) (int, error) {
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return 0, fmt.Errorf("intervals.scanInt: %w", err)
		} else {
			return 0, errors.New("intervals.scanInt: no tokens")
		}
	}

	t := s.Text()

	return strconv.Atoi(t)
}

func isIntervalBound(r rune) bool {
	return isIntervalBegin(r) || isIntervalEnd(r)
}

func isIntervalBeginString(str string) bool {
	r, _ := utf8.DecodeRuneInString(str)

	return isIntervalBegin(r)
}

func isIntervalBegin(r rune) bool {
	return r == '[' || r == '('
}

func isStrictInterval(r rune) bool {
	return r == '(' || r == ')'
}

func isStrictIntervalString(str string) bool {
	r, _ := utf8.DecodeRuneInString(str)

	return isStrictInterval(r)
}

func isIntervalEndString(str string) bool {
	r, _ := utf8.DecodeRuneInString(str)

	return isIntervalEnd(r)
}

func isIntervalEnd(r rune) bool {
	return r == ']' || r == ')'
}

func splitTokens(data []byte, atEof bool) (advance int, token []byte, err error) {
	start := 0

	var (
		r     rune
		width int
	)

	for width = 0; start < len(data); start += width {
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			break
		}
	}

	r, width = utf8.DecodeRune(data[start:])
	if isIntervalBound(r) {
		return start + width, data[start : start+width], nil
	}

	for i := start; start+i < len(data); i += width {
		r, width = utf8.DecodeRune(data[i:])

		if unicode.IsSpace(r) || isIntervalBound(r) {
			return i, data[start:i], nil
		}

		if !unicode.IsDigit(r) {
			return 0, nil, errors.New("intervals.splitTokens: unknown symbol")
		}
	}

	if atEof && len(data) > start {
		return len(data), data[start:], nil
	}

	return start, nil, nil
}
