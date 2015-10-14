package main

import (
	"fmt"
	"github.com/garoago/runelib"
	"html"
	"log"
	"net/http"
	"strings"
)

func main() {
	//http.Handle("/foo", fooHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hellow, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query()
		words := strings.Split(search.Get("q"), " ")
		saved := make(chan bool)
		runeIndex := runelib.GetIndex(saved)
		count := 0
		format := "U+%04X  %c \t%s\n"

		for _, uchar := range runeIndex.Find(words) {
			if uchar > 0xFFFF {
				format = "U+%5X %c \t%s\n"
			}
			fmt.Fprintf(w, format, uchar, uchar, runeIndex.Name(uchar))
			count++
		}

		fmt.Fprintf(w, "%d characters found\n", count)
		if <-saved {
			fmt.Fprint(w, "Index saved.")
		}

		fmt.Fprintf(w, "busca, %q", search)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
