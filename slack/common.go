package slack

import (
	"fmt"
	"os"
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, s)
		panic(e)
	}
}
