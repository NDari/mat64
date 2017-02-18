package matrix

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
)

func printErr(s string) {
	color.Red(s)
	color.Yellow("\nStack trace for this error:\n\n")
	q := string(debug.Stack())
	w := strings.Split(q, "\n")
	fmt.Println(strings.Join(w[7:], "\n"))
	panic(s)
}

func printHelperErr(s string) {
	color.Red(s)
	color.Yellow("\nStack trace for this error:\n\n")
	q := string(debug.Stack())
	w := strings.Split(q, "\n")
	fmt.Println(strings.Join(w[9:], "\n"))
	panic(s)
}
