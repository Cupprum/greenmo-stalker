package osm_test

import (
	"function/geo"
	osm "function/openstreetmaps"
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

	_, err := osm.GenerateMap(server.URL, geo.Position{Lat: 55.787, Lon: 12.519}, nil, "test-key")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGeoapifyMarkers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		marker := r.URL.Query().Get("marker")

		for _, v := range []string{
			"color:#3ea635",
			"size:medium",
			"lonlat:12.515000,55.790000",
			"text:22",
		} {
			if !strings.Contains(marker, v) {
				t.Errorf("green marker missing or wrong: %s, value: %s", marker, v)
				break
			}
		}

		for _, v := range []string{
			"color:#5588d0",
			"size:medium",
			"lonlat:12.520000,55.787000",
			"icon:ev_station",
		} {
			if !strings.Contains(marker, v) {
				t.Errorf("blue marker missing or wrong: %s, value: %s", marker, v)
				break
			}
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte{137, 80, 78, 71})
	}))
	defer server.Close()

	ms := []osm.Marker{
		{Pos: geo.Position{Lat: 55.790, Lon: 12.515}, Color: "#3ea635", Text: "22"},
		{Pos: geo.Position{Lat: 55.787, Lon: 12.520}, Color: "#5588d0", Icon: "ev_station"},
	}

	_, err := osm.GenerateMap(server.URL, geo.Position{Lat: 55.0, Lon: 12.0}, ms, "test-key")
	if err != nil {
		t.Fatal(err)
	}
}
