package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAddition(t *testing.T) {
	req, err := http.NewRequest("GET", "/add?x=2&y=3", nil)
	if err != nil {
		t.Fatal(err)
	}

	cache := make(map[string]Item)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MathHandler(cache,
		func(x float64, y float64) float64 {
			return x + y
		}))
	handler.ServeHTTP(rr, req)

	expected := `{"action":"add","x":2,"y":3,"answer":5,"cached":false}
`
	if rr.Body.String() != expected {
		t.Errorf("Incorrect Result! Got %s, expected %s", rr.Body.String(), expected)
	}
}

func TestSubtraction(t *testing.T) {
	req, err := http.NewRequest("GET", "/subtract?x=2&y=3", nil)
	if err != nil {
		t.Fatal(err)
	}

	cache := make(map[string]Item)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MathHandler(cache,
		func(x float64, y float64) float64 {
			return x - y
		}))
	handler.ServeHTTP(rr, req)

	expected := `{"action":"subtract","x":2,"y":3,"answer":-1,"cached":false}
`
	if rr.Body.String() != expected {
		t.Errorf("Incorrect Result! Got %s, expected %s", rr.Body.String(), expected)
	}
}

func TestMultiplication(t *testing.T) {
	req, err := http.NewRequest("GET", "/multiply?x=2&y=3", nil)
	if err != nil {
		t.Fatal(err)
	}

	cache := make(map[string]Item)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MathHandler(cache,
		func(x float64, y float64) float64 {
			return x * y
		}))
	handler.ServeHTTP(rr, req)

	expected := `{"action":"multiply","x":2,"y":3,"answer":6,"cached":false}
`
	if rr.Body.String() != expected {
		t.Errorf("Incorrect Result! Got %s, expected %s", rr.Body.String(), expected)
	}
}

func TestDivision(t *testing.T) {
	req, err := http.NewRequest("GET", "/multiply?x=3&y=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	cache := make(map[string]Item)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MathHandler(cache,
		func(x float64, y float64) float64 {
			return x / y
		}))
	handler.ServeHTTP(rr, req)

	expected := `{"action":"multiply","x":3,"y":2,"answer":1.5,"cached":false}
`
	if rr.Body.String() != expected {
		t.Errorf("Incorrect Result! Got %s, expected %s", rr.Body.String(), expected)
	}
}

func TestCacheInsertion(t *testing.T) {
	req, err := http.NewRequest("GET", "/multiply?x=3&y=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	cache := make(map[string]Item)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MathHandler(cache,
		func(x float64, y float64) float64 {
			return x / y
		}))
	handler.ServeHTTP(rr, req)
	handler.ServeHTTP(rr, req)

	expected := `{"action":"multiply","x":3,"y":2,"answer":1.5,"cached":false}
{"action":"multiply","x":3,"y":2,"answer":1.5,"cached":true}
`
	if rr.Body.String() != expected {
		t.Errorf("Incorrect Result! Got %s, expected %s", rr.Body.String(), expected)
	}
}

func TestCacheExpiry(t *testing.T) {
	req, err := http.NewRequest("GET", "/multiply?x=3&y=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	cache := make(map[string]Item)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MathHandler(cache,
		func(x float64, y float64) float64 {
			return x / y
		}))
	handler.ServeHTTP(rr, req)
	time.Sleep(time.Minute)
	handler.ServeHTTP(rr, req)

	expected := `{"action":"multiply","x":3,"y":2,"answer":1.5,"cached":false}
{"action":"multiply","x":3,"y":2,"answer":1.5,"cached":false}
`
	if rr.Body.String() != expected {
		t.Errorf("Incorrect Result! Got %s, expected %s", rr.Body.String(), expected)
	}
}

func TestCacheRefresh(t *testing.T) {
	req, err := http.NewRequest("GET", "/multiply?x=3&y=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	cache := make(map[string]Item)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MathHandler(cache,
		func(x float64, y float64) float64 {
			return x / y
		}))
	handler.ServeHTTP(rr, req)
	time.Sleep(30 * time.Second)
	handler.ServeHTTP(rr, req)
	time.Sleep(30 * time.Second)
	handler.ServeHTTP(rr, req)

	expected := `{"action":"multiply","x":3,"y":2,"answer":1.5,"cached":false}
{"action":"multiply","x":3,"y":2,"answer":1.5,"cached":true}
{"action":"multiply","x":3,"y":2,"answer":1.5,"cached":true}
`
	if rr.Body.String() != expected {
		t.Errorf("Incorrect Result! Got %s, expected %s", rr.Body.String(), expected)
	}
}
