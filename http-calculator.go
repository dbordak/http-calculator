package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func argHandler(q url.Values) (int, int, error) {
	xs := q.Get("x")
	ys := q.Get("y")
	if xs == "" || ys == "" {
		return 0, 0, errors.New("Argument Missing")
	}

	x, err := strconv.Atoi(xs)
	if err != nil {
		return 0, 0, err
	}

	y, err := strconv.Atoi(ys)
	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}

func main() {
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		x, y, err := argHandler(r.URL.Query())
		if err != nil {
			// TODO
		}

		fmt.Fprintf(w, "x+y=%d", x+y)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
