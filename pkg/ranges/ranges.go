package ranges

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/niki-timofe/traefikretryplugin/pkg/set"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Range interface {
	Includes(num int) bool
	String() string
}

func NewRange(str string) (Range, error) {
	reader := strings.NewReader(str)

	scanner := bufio.NewScanner(reader)

	scanner.Split(splitTokens)

	return readRange(scanner)
}

type rangeBound struct {
	int
	strict bool
}
type boundedRange struct {
	from rangeBound
	to   rangeBound
}

type rangesDefinition struct {
	b []*boundedRange
	v set.Set[int]
}

func (r *rangesDefinition) String() string {
	var sb strings.Builder

	for i, b := range r.b {
		if i > 0 {
			sb.WriteRune(' ')
		}

		sb.WriteString(b.String())
	}

	for i, value := range r.v.Values() {
		if i > 0 || len(r.b) > 0 {
			sb.WriteRune(' ')
		}

		sb.WriteString(strconv.Itoa(value))
	}

	return sb.String()
}

func (r *rangesDefinition) Includes(num int) bool {
	if r.v.Includes(num) {
		return true
	}

	for _, br := range r.b {
		if br.Includes(num) {
			return true
		}
	}

	return false
}

func (r *boundedRange) Includes(num int) bool {
	from := r.from.int < num || !r.from.strict && r.from.int == num

	to := r.to.int > num || !r.to.strict && r.to.int == num

	return from && to
}

func (r *boundedRange) String() string {
	var sb strings.Builder

	if r.from.strict {
		sb.WriteRune('(')
	} else {
		sb.WriteRune('[')
	}

	sb.WriteString(strconv.Itoa(r.from.int))
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(r.to.int))

	if r.to.strict {
		sb.WriteRune(')')
	} else {
		sb.WriteRune(']')
	}

	return sb.String()
}

func (r *boundedRange) intersect(b *boundedRange) *boundedRange {
	if r.to.int > b.from.int && r.from.int < b.to.int {
		return &boundedRange{
			from: b.from,
			to:   r.to,
		}
	}
	if !r.to.strict && r.to.int == b.from.int {
		return &boundedRange{
			from: b.from,
			to:   r.to,
		}
	}

	return nil
}

func intersectRanges(a *boundedRange, b *boundedRange) *boundedRange {
	if mr := a.intersect(b); mr != nil {
		return mr
	}

	return b.intersect(a)
}

func readRange(s *bufio.Scanner) (Range, error) {
	rgs := make([]*boundedRange, 0)

	values := set.NewSet[int]()

	for s.Scan() {
		token := s.Text()
		if isRngBeginString(token) {
			r, err := scanBounded(s, token)
			if err != nil {
				return nil, fmt.Errorf("ranges.readRange: %w", err)
			}

			rgs = append(rgs, r)
		} else {
			i, err := strconv.Atoi(token)
			if err != nil {
				return nil, fmt.Errorf("ranges.readRange: %w", err)
			}

			values.Add(i)
		}
	}

	if err := s.Err(); err != nil && err != io.EOF {
		return nil, fmt.Errorf("ranges.readRange: %w", err)
	}

	merged := mergeRanges(rgs)

	return &rangesDefinition{
		b: merged,
		v: mergeValues(values, merged),
	}, nil
}

func mergeValues(values set.Set[int], ranges []*boundedRange) set.Set[int] {
	included := set.NewSet[int]()

	for _, i := range values.Values() {
		for _, r := range ranges {
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

func mergeRanges(rgs []*boundedRange) []*boundedRange {
	if len(rgs) <= 1 {
		return rgs
	}

	merged := make([]*boundedRange, 0, len(rgs))

	idx := set.NewSet[int]()

	for i, r := range rgs {
		if idx.Includes(i) {
			continue
		}

		start := i + 1
		for j, rr := range rgs[start:] {
			if mr := intersectRanges(r, rr); mr != nil {
				idx.Add(i)
				idx.Add(start + j)

				merged = append(merged, mr)
			}
		}

		if idx.Includes(i) {
			merged = append(merged, r)
		}
	}

	if idx.Len() > 0 {
		return mergeRanges(merged)
	} else {
		return merged
	}
}

func scanBounded(s *bufio.Scanner, token string) (*boundedRange, error) {
	strict := isStrictRngString(token)
	i, err := scanInt(s)
	if err != nil {
		return nil, fmt.Errorf("ranges.scanBounded: %w", err)
	}

	from := rangeBound{
		int:    i,
		strict: strict,
	}

	i, err = scanInt(s)
	if err != nil {
		return nil, fmt.Errorf("ranges.scanBounded: %w", err)
	}

	if !s.Scan() {
		if err = s.Err(); err != nil {
			return nil, fmt.Errorf("ranges.scanBounded: %w", err)
		}
	}

	token = s.Text()
	if !isRngEndString(token) {
		return nil, errors.New("ranges.scanBounded: unclosed range")
	}

	to := rangeBound{
		int:    i,
		strict: isStrictRngString(token),
	}

	r := &boundedRange{
		from: from,
		to:   to,
	}

	if from.int > to.int || from.strict && to.strict && from.int == to.int {
		return nil, fmt.Errorf("scanBounded: range `%s` is empty", r.String())
	}

	return r, nil
}

func scanInt(s *bufio.Scanner) (int, error) {
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return 0, fmt.Errorf("ranges.scanInt: %w", err)
		} else {
			return 0, errors.New("ranges.scanInt: no tokens")
		}
	}

	t := s.Text()

	return strconv.Atoi(t)
}

func isRngBound(r rune) bool {
	return isRngBegin(r) || isRngEnd(r)
}

func isRngBeginString(str string) bool {
	r, _ := utf8.DecodeRuneInString(str)

	return isRngBegin(r)
}

func isRngBegin(r rune) bool {
	return r == '[' || r == '('
}

func isStrictRng(r rune) bool {
	return r == '(' || r == ')'
}

func isStrictRngString(str string) bool {
	r, _ := utf8.DecodeRuneInString(str)

	return isStrictRng(r)
}

func isRngEndString(str string) bool {
	r, _ := utf8.DecodeRuneInString(str)

	return isRngEnd(r)
}

func isRngEnd(r rune) bool {
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
	if isRngBound(r) {
		return start + width, data[start : start+width], nil
	}

	for i := start; start+i < len(data); i += width {
		r, width = utf8.DecodeRune(data[i:])

		if unicode.IsSpace(r) || isRngBound(r) {
			return i, data[start:i], nil
		}

		if !unicode.IsDigit(r) {
			return 0, nil, errors.New("ranges.splitTokens: unknown symbol")
		}
	}

	if atEof && len(data) > start {
		return len(data), data[start:], nil
	}

	return start, nil, nil
}
