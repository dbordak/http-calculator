package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

type Item struct {
	Response Response
	Expiry   int64
}

const expiry = time.Minute

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

func mathHandler(name string, cache map[string]Item, mathFunc MathFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		rkey := r.URL.Query().Get("x") + name + r.URL.Query().Get("y")
		cached, present := cache[rkey]
		if present && cached.Expiry > time.Now().UnixNano() {
			json.NewEncoder(w).Encode(cached.Response)
			return
		}

		x, y, err := argHandler(r.URL.Query())
		m := Response{name, x, y, 0, false, ""}
		if err != nil {
			m.Error = fmt.Sprintf("%s", err)
			json.NewEncoder(w).Encode(m)
		} else {
			m.Answer = mathFunc(x, y)
			json.NewEncoder(w).Encode(m)
		}

		m.Cached = true
		item := Item{m, time.Now().Add(expiry).UnixNano()}
		cache[rkey] = item
	}
}

func main() {
	cache := make(map[string]Item)

	http.HandleFunc("/add", mathHandler("add", cache,
		func(x float64, y float64) float64 {
			return x + y
		}))

	http.HandleFunc("/subtract", mathHandler("subtract", cache,
		func(x float64, y float64) float64 {
			return x - y
		}))

	http.HandleFunc("/multiply", mathHandler("multiply", cache,
		func(x float64, y float64) float64 {
			return x * y
		}))

	http.HandleFunc("/divide", mathHandler("divide", cache,
		func(x float64, y float64) float64 {
			return x / y
		}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
