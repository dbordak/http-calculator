package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Response struct {
	Action string `json:"action"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Answer int    `json:"answer"`
	Cached bool   `json:"cached"`
	Error  string `json:"error,omitempty"`
}

func argHandler(q url.Values) (int, int, error) {
	xs := q.Get("x")
	ys := q.Get("y")
	if len(xs) == 0 || len(ys) == 0 {
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
		w.Header().Set("Content-Type", "application/json")

		x, y, err := argHandler(r.URL.Query())
		if err != nil {
			m := Response{"add", 0, 0, 0, false, fmt.Sprintf("%s", err)}
			json.NewEncoder(w).Encode(m)
		} else {
			m := Response{"add", x, y, x + y, false, ""}
			json.NewEncoder(w).Encode(m)
		}
	})

	// subtract, multiply, divide

	log.Fatal(http.ListenAndServe(":8080", nil))
}
