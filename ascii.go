package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type asciiSpecial struct {
	Char        string
	Description string
}

var asciiSpecials = [32]asciiSpecial{
	{"NUL", "null"},
	{"SOH", "start of heading"},
	{"STX", "start of text"},
	{"ETX", "end of text"},
	{"EOT", "end of transmission"},
	{"ENQ", "enquiry"},
	{"ACK", "acknowledge"},
	{"BEL", "bell"},
	{"BS", "backspace"},
	{"TAB", "horizontal tab"},
	{"LF", "NL line feed, new line"},
	{"VT", "vertical tab"},
	{"FF", "NP form feed, new page"},
	{"CR", "carriage return"},
	{"SO", "shift out"},
	{"SI", "shift in"},
	{"DLE", "data link escape"},
	{"DC1", "device control 1"},
	{"DC2", "device control 2"},
	{"DC3", "device control 3"},
	{"DC4", "device control 4"},
	{"NAK", "negative acknowledge"},
	{"SYN", "synchronous idle"},
	{"ETB", "end of trans. block"},
	{"CAN", "cancel"},
	{"EM", "end of medium"},
	{"SUB", "substitute"},
	{"ESC", "escape"},
	{"FS", "file separator"},
	{"GS", "group separator"},
	{"RS", "record separator"},
	{"US", "unit separator"},
}

var allowedCols = []int{1, 2, 4}

func asciiCharString(i int) string {
	switch i {
	case 32:
		return "SPACE"
	case 127:
		return "DEL"
	default:
		if i > 127 {
			panic(fmt.Sprintf("outside ascii: %d", i))
		} else if i < 32 {
			return fmt.Sprintf("%-3s (%s)", asciiSpecials[i].Char, asciiSpecials[i].Description)
		}
		return string(rune(i))
	}
}

func printAsciiTable(w io.Writer, columns int) {
	const columnSpacing = 4

	var allowed bool
	for _, col := range allowedCols {
		if columns == col {
			allowed = true
			break
		}
	}
	if !allowed {
		panic(fmt.Sprintf("unsupported column count %d", columns))
	}

	rows := 128 / columns
	columnWidths := make([]int, columns)

	header := "HEX DEC CHAR"
	header2 := "------------"

	for col := 0; col < columns; col++ {
		for row := 0; row < rows; row++ {
			cell := col*rows + row
			cellLen := 7 + len(asciiCharString(cell))
			if cellLen > columnWidths[col] {
				columnWidths[col] = cellLen
			}
		}

		if len(header) > columnWidths[col] {
			columnWidths[col] = len(header)
		}
	}

	for col := 0; col < columns; col++ {
		fmt.Fprintf(w, header)
		padding := columnWidths[col] - len(header) + columnSpacing
		fmt.Fprintf(w, strings.Repeat(" ", padding))
	}
	fmt.Fprintln(w)

	for col := 0; col < columns; col++ {
		fmt.Fprintf(w, header2)
		padding := columnWidths[col] - len(header2) + columnSpacing
		fmt.Fprintf(w, strings.Repeat(" ", padding))
	}
	fmt.Fprintln(w)

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			cell := col*rows + row
			ch := asciiCharString(cell)
			fmt.Fprintf(w, "%02x %3d %s", cell, cell, ch)

			padding := columnWidths[col] + columnSpacing - len(ch) - 7
			fmt.Fprintf(w, strings.Repeat(" ", padding))
		}
		fmt.Fprintln(w)
	}

}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")

	columns := 4

	if r.URL.Query().Has("m") {
		mv := r.URL.Query().Get("m")
		if mv == "" {
			columns = 1
		} else if mb, err := strconv.ParseBool(mv); err == nil && mb {
			columns = 1
		}
	}

	printAsciiTable(w, columns)
}
