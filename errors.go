package matrix

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
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
