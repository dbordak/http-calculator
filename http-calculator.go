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
	Action string  `json:"action"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Answer float64 `json:"answer"`
	Cached bool    `json:"cached"`
	Error  string  `json:"error,omitempty"`
}

type MathFunc func(float64, float64) float64

func argHandler(q url.Values) (float64, float64, error) {
	xs := q.Get("x")
	ys := q.Get("y")
	if len(xs) == 0 || len(ys) == 0 {
		return 0, 0, errors.New("Argument Missing")
	}

	x, err := strconv.ParseFloat(xs, 64)
	if err != nil {
		return 0, 0, err
	}

	y, err := strconv.ParseFloat(ys, 64)
	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}

func mathHandler(name string, mathFunc MathFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		x, y, err := argHandler(r.URL.Query())
		if err != nil {
			m := Response{name, 0, 0, 0, false, fmt.Sprintf("%s", err)}
			json.NewEncoder(w).Encode(m)
		} else {
			m := Response{name, x, y, mathFunc(x, y), false, ""}
			json.NewEncoder(w).Encode(m)
		}
	}
}

func main() {
	http.HandleFunc("/add", mathHandler("add",
		func(x float64, y float64) float64 {
			return x + y
		}))

	http.HandleFunc("/subtract", mathHandler("subtract",
		func(x float64, y float64) float64 {
			return x - y
		}))

	http.HandleFunc("/multiply", mathHandler("multiply",
		func(x float64, y float64) float64 {
			return x * y
		}))

	http.HandleFunc("/divide", mathHandler("divide",
		func(x float64, y float64) float64 {
			return x / y
		}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
