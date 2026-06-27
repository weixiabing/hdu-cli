package table

import (
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

// PrintStruct prints a table of the given struct
func PrintStruct(in interface{}, tags ...string) {
	if in == nil {
		return
	}

	rows := rows(in, tags...)
	table := tablewriter.NewWriter(os.Stdout)
	headers := append([]string{"Name", "Value"}, tags...)
	headerRow := make([]any, 0, len(headers))
	for _, header := range headers {
		headerRow = append(headerRow, header)
	}
	table.Header(headerRow...)

	for _, v := range rows {
		row := make([]any, 0, len(v))
		for _, column := range v {
			row = append(row, column)
		}
		_ = table.Append(row...)
	}
	_ = table.Render() // Send output
}

// rows returns a slice of strings representing the struct fields
func rows(v interface{}, tags ...string) [][]string {
	t := reflect.TypeOf(v).Elem()
	r := reflect.ValueOf(v)

	rows := make([][]string, 0)

	for i := 0; i < t.NumField(); i++ {
		row := make([]string, 0)
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		row = append(row,
			field.Name,
			fmt.Sprint(reflect.Indirect(r).FieldByName(field.Name)),
		)

		for _, v := range tags {
			column := field.Tag.Get(v)
			row = append(row, column)
		}
		rows = append(rows, row)
	}
	return rows
}
