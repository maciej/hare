package main

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
