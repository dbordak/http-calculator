package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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

// Parses the x and y values out of the query string.
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

// Closure that generates an http response function using the given math
// function.
func mathHandler(cache map[string]Item, mathFunc MathFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		action := r.URL.Path[1:]

		// Create a unique key based on method, x, and y. Can't just use the URL
		// string from the request since x and y can be in either order.
		rkey := r.URL.Query().Get("x") + action + r.URL.Query().Get("y")
		cached, present := cache[rkey]
		if present && cached.Expiry > time.Now().UnixNano() {
			json.NewEncoder(w).Encode(cached.Response)

			// Update expiry
			cache[rkey] = Item{cached.Response, time.Now().Add(expiry).UnixNano()}
			return
		}

		// Create a new response when it's not in the cache (or expired)
		x, y, err := argHandler(r.URL.Query())
		m := Response{action, x, y, 0, false, ""}
		if err != nil {
			m.Error = fmt.Sprintf("%s", err)
			json.NewEncoder(w).Encode(m)
		} else {
			m.Answer = mathFunc(x, y)
			json.NewEncoder(w).Encode(m)
		}

		// Add response to the cache
		m.Cached = true
		item := Item{m, time.Now().Add(expiry).UnixNano()}
		cache[rkey] = item
	}
}

func main() {
	cache := make(map[string]Item)

	http.HandleFunc("/add", mathHandler(cache,
		func(x float64, y float64) float64 {
			return x + y
		}))

	http.HandleFunc("/subtract", mathHandler(cache,
		func(x float64, y float64) float64 {
			return x - y
		}))

	http.HandleFunc("/multiply", mathHandler(cache,
		func(x float64, y float64) float64 {
			return x * y
		}))

	http.HandleFunc("/divide", mathHandler(cache,
		func(x float64, y float64) float64 {
			return x / y
		}))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
