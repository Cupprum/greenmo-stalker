package openstreetmaps_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Position helper for testing
type Position struct {
	Lat float64
	Lon float64
}

func TestGeoapifyParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("width") != "600" {
			t.Errorf("width wrong")
		}
		if r.URL.Query().Get("height") != "600" {
			t.Errorf("height wrong")
		}
		// Spec requires: lonlat:longitude,latitude
		if r.URL.Query().Get("center") != "lonlat:12.519,55.787" {
			t.Errorf("center wrong: expected lonlat:12.519,55.787, got %s", r.URL.Query().Get("center"))
		}
		if r.URL.Query().Get("zoom") != "14" {
			t.Errorf("zoom wrong")
		}
		if r.URL.Query().Get("apiKey") != "test-key" {
			t.Errorf("apiKey missing or incorrect")
		}
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{137, 80, 78, 71}) // PNG Magic Number
	}))
	defer server.Close()

	// TODO: img, err := GenerateMap(server.URL, 55.787, 12.519, []Position{}, []Position{}, "test-key")
}

func TestGeoapifyMarkers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		marker := r.URL.Query().Get("marker")

		// Verify Car Marker (Green: #3ea635)
		// Expected snippet: lonlat:12.515,55.790;color:%233ea635;size:medium
		if !strings.Contains(marker, "lonlat:12.515,55.790") || !strings.Contains(marker, "color:%233ea635") {
			t.Errorf("green car marker improperly formatted or missing: %s", marker)
		}

		// Verify Charger Marker (Red: #f30e0e)
		// Expected snippet: lonlat:12.520,55.787;color:%23f30e0e;size:medium
		if !strings.Contains(marker, "lonlat:12.520,55.787") || !strings.Contains(marker, "color:%23f30e0e") {
			t.Errorf("red charger marker improperly formatted or missing: %s", marker)
		}

		// Verify pipe separator if both exist
		if !strings.Contains(marker, "|") {
			t.Errorf("markers should be separated by a pipe character")
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{137, 80, 78, 71})
	}))
	defer server.Close()

	carPositions := []Position{{Lat: 55.790, Lon: 12.515}}
	chargerPositions := []Position{{Lat: 55.787, Lon: 12.520}}

	// TODO: img, err := GenerateMap(server.URL, 55.787, 12.519, carPositions, chargerPositions, "test-key")
}
