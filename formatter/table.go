package formatter

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mitchellh/colorstring"
	pad "github.com/willf/pad/utf8"
	"io"
)

type Table struct {
	columnSize []int
	headers    []string
	w          io.Writer
	colorizer  func(row []string, index int) string
}

func (t *Table) WriteHeaders() {
	var buffer bytes.Buffer

	buffer.WriteString(" ")
	for i, v := range t.headers {
		if i != 0 {
			buffer.WriteString(" │ ")
		}
		buffer.WriteString(pad.Right(v, t.columnSize[i], " "))
	}
	buffer.WriteString("\n")
	_, _ = t.w.Write(buffer.Bytes())

	buffer = bytes.Buffer{}

	buffer.WriteString(" ")
	for i, _ := range t.headers {
		if i != 0 {
			buffer.WriteString("━┿━")
		}
		buffer.WriteString(pad.Right("", t.columnSize[i], "━"))
	}
	buffer.WriteString("\n")
	_, _ = t.w.Write(buffer.Bytes())

}

func (t *Table) AnalyseRow(r ...string) error {
	if len(r) != len(t.headers) {
		return errors.New("columns count does't match")
	}

	for i, v := range r {
		l := len(v)
		if l > t.columnSize[i] {
			t.columnSize[i] = l
		}
	}

	return nil
}

func (t *Table) WriteRow(row ...string) {
	var buffer bytes.Buffer

	buffer.WriteString(" ")
	for i, v := range row {
		if i != 0 {
			buffer.WriteString(" │ ")
		}

		if t.colorizer != nil {
			color := t.colorizer(row, i)
			if color != "" {
				v = fmt.Sprintf("\033[%sm%s\033[0m", colorstring.DefaultColors[color], v)
			}
		}
		buffer.WriteString(pad.Right(v, t.columnSize[i], " "))
	}
	buffer.WriteString("\n")
	_, _ = t.w.Write(buffer.Bytes())
}

func (t *Table) SetColorizer(f func(row []string, index int) string) {
	t.colorizer = f
}

func NewTable(w io.Writer, cols ...string) *Table {
	t := Table{
		columnSize: make([]int, len(cols)),
		headers:    cols,
		w:          w,
	}

	for i, v := range cols {
		t.columnSize[i] = len(v)
	}

	return &t

}
