package table

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Table table structure
type Table struct {
	Header []string
	Data   [][]string
	Footer []string
}

// GetCols get table cols
func (t *Table) GetCols() int {
	return len(t.Header)
}

// PrintTable show table
func (t *Table) PrintTable(displayFooter bool) error {
	// write to stdout
	table := tablewriter.NewWriter(os.Stdout)

	// setting table style
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(false)

	// table header
	table.SetHeader(t.Header)

	// table style
	var colorList []tablewriter.Colors
	for i := 0; i < len(t.Header); i++ {
		colorList = append(colorList, tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold})
	}
	table.SetHeaderColor(colorList...)

	// table rows
	table.AppendBulk(t.Data)

	// footer
	if displayFooter {
		footer := make([]string, t.GetCols())

		for i := 0; i < t.GetCols()-1; i++ {
			footer[i] = ""
		}

		footer[t.GetCols()-1] = fmt.Sprintf(" Total: %v ", len(t.Data))
		t.Footer = footer

		// table footer
		table.SetFooter(t.Footer)
	}

	table.Render()
	return nil
}
