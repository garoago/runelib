package main

import (
	"encoding/json"
	"fmt"
	"github.com/garoago/runelib"
	"html"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hellow, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query()
		format := search.Get("f")
		words := strings.Split(search.Get("q"), " ")

		if format == "" {
			format = "json"
		}

		if format == "txt" {
			txtSearch(w, words)
		} else if format == "json" {
			jsonSearch(w, words)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

type SearchResponse struct {
	Count  int
	Result runelib.RuneSlice
}

func jsonSearch(w http.ResponseWriter, words []string) {
	saved := make(chan bool)
	runeIndex := runelib.GetIndex(saved)

	result := runeIndex.Find(words)
	count := len(result)

	output := SearchResponse{
		Count:  count,
		Result: result,
	}

	serialized, _ := json.Marshal(output)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(serialized))
}

func txtSearch(w http.ResponseWriter, words []string) {
	saved := make(chan bool)
	runeIndex := runelib.GetIndex(saved)
	count := 0
	template := "U+%04X  %c \t%s\n"

	for _, uchar := range runeIndex.Find(words) {
		if uchar > 0xFFFF {
			template = "U+%5X %c \t%s\n"
		}
		fmt.Fprintf(w, template, uchar, uchar, runeIndex.Name(uchar))
		count++
	}

	fmt.Fprintf(w, "%d characters found\n", count)
	if <-saved {
		fmt.Fprint(w, "Index saved.")
	}
}
