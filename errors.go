package matrix

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

var (
	wrongArity    = "In %s: expected %s arguments, got %d"
	wrongArgType  = "In %s: expected input type %s, got %T"
	sizeMismatch  = "In %s: size mismatch: %dx%d vs %dx%d"
	colOutOfBound = "In %s the column %d is outside of bounds [%d, %d)"
	rowOutOfBound = "In %s the column %d is outside of bounds [%d, %d)"
)

func printErr(s string) {
	fmt.Println(s)
	q := string(debug.Stack())
	w := strings.Split(q, "\n")
	fmt.Println(strings.Join(w[7:], "\n"))
	os.Exit(1)
}

func printHelperErr(s string) {
	fmt.Println(s)
	q := string(debug.Stack())
	w := strings.Split(q, "\n")
	fmt.Println(strings.Join(w[9:], "\n"))
	os.Exit(1)
}
