package util

import "fmt"

const (
	RESET = "\x1b[0m"
	GREEN = "\x1b[32;1m"
	RED   = "\x1b[31;1m"
)

func GreenPrintln(txt string) {
	fmt.Printf("%s%v%s\n", GREEN, txt, RESET)
}

func RedPrintln(txt string) {
	fmt.Printf("%s%v%s\n", RED, txt, RESET)
}

// GreenError dyes the error green.
func GreenError(errTxt string) error {
	return fmt.Errorf("%s%v%s", GREEN, errTxt, RESET)
}

// RedError dyes the error red.
func RedError(errTxt string) error {
	return fmt.Errorf("%s%v%s", RED, errTxt, RESET)
}
