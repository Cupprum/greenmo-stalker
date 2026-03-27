package greenmobility_test

import (
	"encoding/json"
	"function/geo"
	"function/greenmobility"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreenMobilityParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("lat") != "55.787000" {
			t.Errorf("lat param wrong, got %s", r.URL.Query().Get("lat"))
		}
		if r.URL.Query().Get("lng") != "12.519000" {
			t.Errorf("lng param wrong")
		}
		if r.URL.Query().Get("rad") != "2.0" {
			t.Errorf("rad param wrong, got %s", r.URL.Query().Get("rad"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	_, err := greenmobility.Query(server.URL, geo.Position{Lat: 55.787, Lon: 12.519}, 2.0, 40)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGreenMobilityFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := []map[string]interface{}{
			{
				"stateOfCharge": 35,
				"position":      map[string]interface{}{"coordinates": [2]float64{12.515, 55.790}},
			},
			{
				"stateOfCharge": 50,
				"position":      map[string]interface{}{"coordinates": [2]float64{12.525, 55.785}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cars, err := greenmobility.Query(server.URL, geo.Position{Lat: 55.787, Lon: 12.519}, 2.0, 40)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(cars) != 1 {
		t.Fatalf("Expected 1 car after filtering, got %d", len(cars))
	}

	if cars[0].Lat != 55.790 || cars[0].Lon != 12.515 {
		t.Errorf("Coordinate mapping failed, got %+v", cars[0])
	}
}
