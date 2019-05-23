package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joseluisq/redel"
)

func main() {
	r, err := os.Open("case/src/test.txt")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer r.Close()

	w, err := os.Create("case/src/test.mod.txt")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer w.Close()

	var writer = bufio.NewWriter(w)

	// rep := redel.NewRedel(r, "require('", "')", "+++++")
	rep := redel.NewRedel(r, "START", "END", "+++++")

	replaceFunc := func(data []byte, atEOF bool) {
		_, err := writer.Write(data)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	filterFunc := func(matchValue []byte) []byte {
		value := string(matchValue)

		fmt.Println("MATCH VALUE::", value)
		fmt.Println("------")

		// if value == "~/222C" || value == "~/111B" {
		// 	return false
		// }

		// if value == "~/222C" || value == "~/111B" {
		// 	return false
		// }

		if value == " slice, " {
			// TODO: Fix Redel to support extra chars
			return append(matchValue, []byte("====")...)
		}

		// if value == "~/222C" {
		// 	return append(matchValue, []byte("AAA")...)
		// }

		// if value == " slice, " {
		// 	return false
		// }

		return matchValue
	}

	rep.FilterReplaceWith(replaceFunc, filterFunc, true)

	writer.Flush()
}
