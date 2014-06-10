package termtable

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

type Table struct {
	Rows    [][]string
	Columns [][]string
	Options *TableOptions

	HasHeader bool

	numColumns   int
	columnsWidth []int

	hasRendered     bool
	numRenderedRows int
}

type TableOptions struct {
	Padding      int
	UseSeparator bool
}

var defaultTableOptions = &TableOptions{
	Padding:      1,
	UseSeparator: false,
}

func NewTable(rows [][]string, options *TableOptions) *Table {
	t := &Table{
		Options: options,
	}
	if t.Options == nil {
		t.Options = defaultTableOptions
	}
	if rows != nil {
		t.Rows = rows
		t.computeProperties()
	}
	return t
}

func (t *Table) SetHeader(header []string) {
	t.HasHeader = true
	// There is a better way to do this
	t.Rows = append([][]string{header}, t.Rows...)
	t.computeProperties()
}

func (t *Table) AddRow(row []string) {
	t.Rows = append(t.Rows, row)
	if t.hasRendered {
		// dynamicly format row
	} else {
		t.computeProperties()
	}
	t.numRenderedRows++
}

func (t *Table) computeProperties() {
	if len(t.Rows) > 0 {
		t.numColumns = len(t.Rows[0])
		t.columnsWidth = make([]int, t.numColumns)
		t.recalculate()
	}
}

func (t *Table) recalculate() {
	t.Columns = [][]string{}
	for i := 0; i < t.numColumns; i++ {
		t.Columns = append(t.Columns, []string{})
	}
	for _, row := range t.Rows {
		for j, cellContent := range row {
			t.Columns[j] = append(t.Columns[j], cellContent)
			t.columnsWidth[j] = int(math.Max(float64(len(cellContent)), float64(t.columnsWidth[j])))
		}
	}
}

func (t *Table) Render() string {
	// allocate a 1k byte buffer
	bb := make([]byte, 0, 1024)
	buf := bytes.NewBuffer(bb)

	i := 0

	if !t.hasRendered {
		if t.HasHeader {
			if t.Options.UseSeparator {
				buf.WriteString(t.separatorLine())
				buf.WriteRune('\n')
			}
			for j := range t.Rows[0] {
				buf.WriteString(t.getCell(i, j))
			}
			i = 1
			buf.WriteRune('\n')
		}

		if t.Options.UseSeparator {
			buf.WriteString(t.separatorLine())
			buf.WriteRune('\n')
		}
	} else {
		i = t.numRenderedRows
		fmt.Printf("\033[1A")
	}

	for i < len(t.Rows) {
		row := t.Rows[i]
		for j := range row {
			buf.WriteString(t.getCell(i, j))
		}
		if i < len(t.Rows)-1 {
			buf.WriteRune('\n')
		}
		i++
	}

	if t.Options.UseSeparator {
		buf.WriteRune('\n')
		buf.WriteString(t.separatorLine())
	}

	t.hasRendered = true

	return buf.String()
}

func (t *Table) separatorLine() string {
	sep := "+"
	for _, w := range t.columnsWidth {
		sep += strings.Repeat("-", w+2*t.Options.Padding)
		sep += "+"
	}
	return sep
}

func (t *Table) getCell(row, col int) string {
	cellContent := t.Rows[row][col]
	spacePadding := strings.Repeat(" ", t.Options.Padding)

	var cellStr string

	if t.Options.UseSeparator {
		cellStr += "|"
		cellStr += spacePadding
	}

	cellStr += cellContent
	cellStr += strings.Repeat(" ", t.columnsWidth[col]-len(cellContent))
	cellStr += spacePadding

	if t.Options.UseSeparator {
		if col == t.numColumns-1 {
			cellStr += "|"
		}
	}

	return cellStr
}
