package formatter

type Table struct {
	columnSize []int
	headers    []string
}

func NewTable(cols ...string) *Table {
	t := Table{
		columnSize: make([]int, len(cols)),
		headers:    cols,
	}

	for i, v := range cols {
		t.columnSize[i] = len(v)
	}

	return &t

}
