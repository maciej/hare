package main

import (
	"fmt"
	"io"
	"net/http"
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

var longestDesc int

func init() {
	for _, s := range asciiSpecials {
		if len(s.Description) > longestDesc {
			longestDesc = len(s.Description)
		}
	}
}

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

func printAsciiTable(w io.Writer) {
	const columnSpacing = 4
	for col := 0; col < 4; col++ {
		fmt.Fprintf(w, "HEX DEC CHAR")
		if col == 0 {
			fmt.Fprintf(w, strings.Repeat(" ", longestDesc+1+columnSpacing))
		} else if col != 3 {
			fmt.Fprintf(w, strings.Repeat(" ", columnSpacing))
		}
	}
	fmt.Fprintln(w)

	for col := 0; col < 4; col++ {
		fmt.Fprintf(w, "------------")
		if col == 0 {
			fmt.Fprintf(w, strings.Repeat(" ", longestDesc+1+columnSpacing))
		} else if col != 3 {
			fmt.Fprintf(w, strings.Repeat(" ", columnSpacing))
		}
	}
	fmt.Fprintln(w)

	for i := 0; i < 32; i++ {
		for col := 0; col < 4; col++ {
			cell := i + col*32
			ch := asciiCharString(cell)

			fmt.Fprintf(w, "%02x %3d %s", cell, cell, ch)
			if col == 0 {
				fmt.Fprintf(w, strings.Repeat(" ", longestDesc-len(ch)+6+columnSpacing))
			} else if col != 3 {
				fmt.Fprintf(w, strings.Repeat(" ", 5+columnSpacing-len(ch)))
			}

		}
		fmt.Fprintln(w)
	}
}

func asciiHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "text/plain")
	printAsciiTable(w)
}
