package main

import (
	"fmt"
	"github.com/garoago/runelib"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:  runefinder <word>\texample: runefinder cat")
		os.Exit(1)
	}

	words := os.Args[1:]
	saved := make(chan bool)
	runeIndex := runelib.GetIndex(saved)
	count := 0
	format := "U+%04X  %c \t%s\n"

	for _, uchar := range runeIndex.Find(words) {
		if uchar > 0xFFFF {
			format = "U+%5X %c \t%s\n"
		}
		fmt.Printf(format, uchar, uchar, runeIndex.Name(uchar))
		count++
	}
	fmt.Printf("%d characters found\n", count)
	if <-saved {
		fmt.Println("Index saved.")
	}
}
