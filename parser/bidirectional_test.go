package parser

import (
	"bytes"
	"github.com/3onyc/3do/util"
	"github.com/kr/pretty"
	"testing"
)

var ParserTestData = []string{`### Group 1

* Item 1

   Some content
   _which can be markdown_

* Item 2

### Group 2

* Item 1`, `### Group 1

* Item 1

   # Foo

* Item 2

   bar`}

func TestBoth(t *testing.T) {
	for _, d := range ParserTestData {
		l, err := NewParser().Parse("Title", bytes.NewBufferString(d))
		if err != nil {
			t.Error(err)
		}

		buf := bytes.NewBufferString("")
		NewWriter().Write(l, buf)

		if buf.String() != d {
			t.Error("Buffer does not match ParserTestData")
			t.Error("==== Object ====")
			t.Errorf("% #v", pretty.Formatter(l))
			t.Error("==== Buffer ====")
			t.Error(util.ShowNewLines(buf.String()))
			t.Error("=== ParserTestData ===")
			t.Error(util.ShowNewLines(d))
		}
	}
}
