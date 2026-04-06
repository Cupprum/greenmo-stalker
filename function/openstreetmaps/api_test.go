package openstreetmaps_test

import (
	"function/geo"
	"function/openstreetmaps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeoapifyParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("center") != "lonlat:12.519000,55.787000" {
			t.Errorf("center wrong, got %s", r.URL.Query().Get("center"))
		}
		if r.URL.Query().Get("apiKey") != "test-key" {
			t.Errorf("apiKey missing")
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte{137, 80, 78, 71})
	}))
	defer server.Close()

	_, err := openstreetmaps.GenerateMap(server.URL, geo.Position{Lat: 55.787, Lon: 12.519}, nil, nil, "test-key")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGeoapifyMarkers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		marker := r.URL.Query().Get("marker")

		if !strings.Contains(marker, "color:#3ea635;size:medium") || !strings.Contains(marker, "lonlat:12.515000,55.790000") {
			t.Errorf("green marker missing or wrong: %s", marker)
		}
		if !strings.Contains(marker, "type:material;color:#5588d0;icon:ev_station;size:medium") || !strings.Contains(marker, "lonlat:12.520000,55.787000") {
			t.Errorf("blue marker missing or wrong: %s", marker)
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte{137, 80, 78, 71})
	}))
	defer server.Close()

	cars := []geo.Position{{Lat: 55.790, Lon: 12.515}}
	chargers := []geo.Position{{Lat: 55.787, Lon: 12.520}}

	_, err := openstreetmaps.GenerateMap(server.URL, geo.Position{Lat: 55.0, Lon: 12.0}, cars, chargers, "test-key")
	if err != nil {
		t.Fatal(err)
	}
}
