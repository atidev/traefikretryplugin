package structuredheaders

import (
	"fmt"
	"net/http"
	"testing"
)

var header = http.Header{}

func TestItem(t *testing.T) {
	header.Add("Integer-Item", "10; parameter=10")
	header.Add("Integer-Item-Fail", "10000000000000000000000000000000001; parameter=10")
	header.Add("Decimal-Item", "10.10; parameter=10.10")
	header.Add("Decimal-Item-Fail", "100000000000000000000000000000000000000000000000000000000000.10.0; parameter=10.10")
	header.Add("String-Item", "\"string\\\"\"; parameter=\"str\"")
	header.Add("String-Item-Fail", "\"string\\a\"; parameter=\"str\"")
	header.Add("String-Item-Fail-1", "\"string\a\"; parameter=\"str\"")
	header.Add("String-Item-Fail-2", "\"string")
	header.Add("Token-Item", "Token; parameter=Token")
	header.Add("Binary-Item", ":QmluYXJ5:; parameter=:QmluYXJ5:")
	header.Add("Binary-Item-Fail", ":Qmlu^YXJ5:; parameter=:QmluYXJ5:")
	header.Add("Boolean-Item", "?1; p; param=?0")
	header.Add("Boolean-Item-1", "?0; parameter=?1")
	header.Add("Boolean-Item-Fail", "?3; parameter")
	header.Add("Boolean-Item-Params-Fail", "?0; parameter; -p=?1")
	header.Add("Boolean-Item-Params-Fail-1", "?0; parameter; p=\"st")
	header.Add("Boolean-Item-Empty-Fail", "")
	header.Add("Unknown-Item-Fail", "\a")

	sh := NewStructuredHeader(header)

	intI, err := sh.Item("Integer-Item")
	if err != nil {
		t.Error(err)
	}

	n, err := intI.Number()
	if err != nil {
		t.Error(err)
	}

	i, err := n.Integer()
	if err != nil {
		t.Error(err)
	}
	t.Log(i)

	_, err = intI.Boolean()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	_, err = intI.Str()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	_, err = intI.Token()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}
	_, err = intI.Binary()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	intI, err = sh.Item("Integer-Item-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	f, err := n.Float()
	if err != nil {
		t.Error(err)
	}
	t.Log(f)

	fltI, err := sh.Item("Decimal-Item")
	if err != nil {
		t.Error(err)
	}

	n, err = fltI.Number()
	if err != nil {
		t.Error(err)
	}

	f, err = n.Float()
	t.Log(f)

	i, err = n.Integer()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	fltI, err = sh.Item("Decimal-Item-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	strI, err := sh.Item("String-Item")
	if err != nil {
		t.Error(err)
	}

	s, err := strI.Str()
	if err != nil {
		t.Error(err)
	}
	t.Log(s)

	_, err = strI.Number()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	strI, err = sh.Item("String-Item-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	strI, err = sh.Item("String-Item-Fail-1")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	strI, err = sh.Item("String-Item-Fail-2")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	tknI, err := sh.Item("Token-Item")
	if err != nil {
		t.Error(err)
	}

	tk, err := tknI.Token()
	t.Log(tk)

	tknI, err = sh.Item("Token-Item-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	binI, err := sh.Item("Binary-Item")
	if err != nil {
		t.Error(err)
	}

	b, err := binI.Binary()
	t.Log(b)

	binI, err = sh.Item("Binary-Item-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	bI, err := sh.Item("Boolean-Item")
	if err != nil {
		t.Error(err)
	}

	bo, err := bI.Boolean()
	t.Log(bo)
	t.Log(bI.Parameters())

	bI, err = sh.Item("Boolean-Item-1")
	if err != nil {
		t.Error(err)
	}

	bo, err = bI.Boolean()
	t.Log(bo)

	bI, err = sh.Item("Boolean-Item-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	bI, err = sh.Item("Boolean-Item-Params-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	bI, err = sh.Item("Boolean-Item-Params-Fail-1")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	bI, err = sh.Item("Boolean-Item-Empty-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}

	bI, err = sh.Item("Unknown-Item-Fail")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("expected err")
	}
}

func TestList(t *testing.T) {
	header.Add("List", "10; a, 11, (1 2); b")
	header.Add("List-Fail", "10-1; a, 11, (1 2)")
	header.Add("List-Fail-1", "10; a, 11, ^1 2")
	header.Add("List-Fail-2", "10; a, 11, (1 2")
	header.Add("List-Fail-3", "10; a, 11, (1 2); b=\"unt-str")
	header.Add("List-Fail-4", "10; a, 11, (1 \"unt-str)")

	sh := NewStructuredHeader(header)

	l, err := sh.List("List")

	if err != nil {
		t.Error(err)
	}

	for _, li := range l {
		itm, err := li.Item()
		lItm, lErr := li.InnerList()

		if err != nil && lErr != nil {
			t.Error(err)
		}

		if lErr == nil {
			t.Log(lItm.Parameters())

			for _, it := range lItm.Items() {
				t.Log(fmt.Sprintf("inner: %#v", it))
			}
		} else {
			t.Log(fmt.Sprintf("%#v", itm))
		}
	}

	l, err = sh.List("List-Fail")
	if err == nil {
		t.Error("expected error")
	}

	t.Log(err)

	l, err = sh.List("List-Fail-1")
	if err == nil {
		t.Error("expected error")
	}

	t.Log(err)

	l, err = sh.List("List-Fail-2")
	if err == nil {
		t.Error("expected error")
	}

	t.Log(err)

	l, err = sh.List("List-Fail-3")
	if err == nil {
		t.Error("expected error")
	}

	t.Log(err)

	l, err = sh.List("List-Fail-4")
	if err == nil {
		t.Error("expected error")
	}

	t.Log(err)
}

func TestDictionary(t *testing.T) {
	header.Add("Dictionary", "a=\"(500 100)\"; pa=ram, b, c=Token, d, e=(?1 2), f")
	header.Add("Dictionary-Fail", "a=\"(500 100); pa=ram, b=Token")
	header.Add("Dictionary-Fail-1", "'=Token")
	header.Add("Dictionary-Fail-2", "a~")
	header.Add("Dictionary-Fail-3", "a=\"(500 100)\"; pa=ram& b=Token")

	sh := NewStructuredHeader(header)

	d, err := sh.Dictionary("Dictionary")

	if err != nil {
		t.Error(err)
	}

	for k, li := range d {
		itm, err := li.Item()
		lItm, lErr := li.InnerList()

		if err != nil && lErr != nil {
			t.Error(err)
		}

		t.Log(fmt.Sprintf("[%s]: ", k))

		if lErr == nil {
			t.Log(lItm.Parameters())

			for _, it := range lItm.Items() {
				t.Log(fmt.Sprintf("inner: %#v", it))
			}
		} else {
			t.Log(fmt.Sprintf("%#v", itm))
		}
	}

	d, err = sh.Dictionary("Dictionary-Fail")
	if err == nil {
		t.Error("expected err")
	}

	t.Log(err)

	d, err = sh.Dictionary("Dictionary-Fail-1")
	if err == nil {
		t.Error("expected err")
	}

	t.Log(err)

	d, err = sh.Dictionary("Dictionary-Fail-2")
	if err == nil {
		t.Error("expected err")
	}

	t.Log(err)

	d, err = sh.Dictionary("Dictionary-Fail-3")
	if err == nil {
		t.Error("expected err")
	}

	t.Log(err)
}
