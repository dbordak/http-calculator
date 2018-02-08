package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		a := r.URL.Query().Get("a")
		b := r.URL.Query().Get("b")
		if a == "" || b == "" {
			return
		}
		fmt.Fprintf(w, "a=%q, b=%q", a, b)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
