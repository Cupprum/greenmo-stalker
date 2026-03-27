package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeoapifyParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("width") != "600" {
			t.Errorf("width wrong")
		}
		if r.URL.Query().Get("height") != "600" {
			t.Errorf("height wrong")
		}
		if r.URL.Query().Get("center") != "lonlat:12.519,55.787" {
			t.Errorf("center wrong")
		}
		if r.URL.Query().Get("zoom") != "14" {
			t.Errorf("zoom wrong")
		}
		if r.URL.Query().Get("apiKey") == "" {
			t.Errorf("apiKey missing")
		}
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{137, 80, 78, 71})
	}))
	defer server.Close()

	// TODO: img, err := GenerateMap(server.URL, 55.787, 12.519, []Position{}, []Position{}, "test-key")
}

func TestGeoapifyMarkers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		marker := r.URL.Query().Get("marker")
		if !strings.Contains(marker, "color:%233ea635") {
			t.Errorf("green marker missing")
		}
		if !strings.Contains(marker, "color:%23f30e0e") {
			t.Errorf("red marker missing")
		}
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{137, 80, 78, 71})
	}))
	defer server.Close()

	// TODO: img, err := GenerateMap(server.URL, 55.787, 12.519, carPositions, chargerPositions, "test-key")
}

type Position struct {
	Lat float64
	Lon float64
}
