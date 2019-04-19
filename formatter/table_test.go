package formatter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestTableHeaders(t *testing.T) {
	table := NewTable(nil, "a", "bb", "ccc", "dddd")
	assert.Equal(t, []int{1, 2, 3, 4}, table.columnSize)
}

func TestTableOutput(t *testing.T) {

	var b bytes.Buffer
	w := io.Writer(&b)

	table := NewTable(w, "a", "bb", "ccc", "dddd")
	table.WriteHeaders()

	expected := ` a │ bb │ ccc │ dddd
 ━━┿━━━━┿━━━━━┿━━━━━
`
	assert.Equal(t, expected, b.String())

}

func TestTableOutputWithRows(t *testing.T) {

	var b bytes.Buffer
	w := io.Writer(&b)

	table := NewTable(w, "beautiful", "header", "isn't it")
	row := []string{"In sed", "accumsan", "Lorem ipsum dolor sit amet, consectetur"}
	err := table.AnalyseRow(row...)
	assert.NoError(t, err)

	table.WriteHeaders()
	table.WriteRow(row...)
	expected := ` beautiful │ header   │ isn't it                               
 ━━━━━━━━━━┿━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 In sed    │ accumsan │ Lorem ipsum dolor sit amet, consectetur
`
	assert.Equal(t, expected, b.String())

}
