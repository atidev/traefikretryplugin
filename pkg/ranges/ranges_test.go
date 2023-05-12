package ranges

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

const test = "[400 500 ] 505   [111 450] (505 600] 505 400 200"

const testErr = "[300 500] abc [300 400]"

const testClean = "[111 500] (505 600] 505"

var testTokens = []string{"[", "400", "500", "]", "505", "[", "111", "450", "]", "(", "505", "600", "]", "505", "400", "200"}

func TestSplitTokens(t *testing.T) {
	reader := strings.NewReader(test)

	scanner := bufio.NewScanner(reader)

	scanner.Split(splitTokens)

	assert := func(i int, exp string, act string) {
		if exp != act {
			t.Errorf("expected %s got %s", exp, act)
		}
	}

	for i := 0; scanner.Scan(); i++ {
		assert(i, testTokens[i], scanner.Text())
	}
}

func TestStingRepresentation(t *testing.T) {
	r, err := NewRange(test)

	if err != nil {
		t.Error(err)
	}

	assert := func(exp string, act string) {
		if exp != act {
			t.Errorf("expected %s got %s", exp, act)
		}
	}

	assert(testClean, fmt.Sprint(r))
}

func TestErr(t *testing.T) {
	_, err := NewRange(testErr)

	t.Logf("%v", err)

	if err == nil {
		t.Errorf("no err")
	}
}

func BenchmarkRangeCreate(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := NewRange(test)
		if err != nil {
			b.Fatalf("%v", err)
		}
	}
}
