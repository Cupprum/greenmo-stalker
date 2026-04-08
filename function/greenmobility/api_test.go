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
		if r.URL.Query().Get("rad") != "1.3" {
			t.Errorf("rad param wrong, got %s", r.URL.Query().Get("rad"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	nw, se := geo.Position{Lat: 55.777000, Lon: 12.509000}, geo.Position{Lat: 55.797000, Lon: 12.529000}
	_, err := greenmobility.Query(server.URL, nw, se, 40)
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

	nw, se := geo.Position{Lat: 55.800000, Lon: 12.500000}, geo.Position{Lat: 55.700000, Lon: 12.600000}
	cars, err := greenmobility.Query(server.URL, nw, se, 40)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(cars) != 1 {
		t.Fatalf("Expected 1 car after filtering, got %d", len(cars))
	}

	if cars[0].Pos.Lat != 55.790 || cars[0].Pos.Lon != 12.515 {
		t.Errorf("Coordinate mapping failed, got %+v", cars[0])
	}
}

func TestGreenMobilityDiscounts(t *testing.T) {
	mux := http.NewServeMux()

	// Mock the main endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		resp := []map[string]interface{}{
			{
				"id":            101,
				"stateOfCharge": 30,
				"benefit":       "DISCOUNTED", // Should trigger discount fetch
				"position":      map[string]interface{}{"coordinates": [2]float64{12.515, 55.790}},
			},
			{
				"id":            102,
				"stateOfCharge": 30,
				"position":      map[string]interface{}{"coordinates": [2]float64{12.516, 55.791}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Mock the specific car endpoint
	mux.HandleFunc("/101", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"discount": {"discountPercentage": 0.15}}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	nw, se := geo.Position{Lat: 55.800000, Lon: 12.500000}, geo.Position{Lat: 55.700000, Lon: 12.600000}

	cars, err := greenmobility.Query(server.URL, nw, se, 40)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(cars) != 2 {
		t.Fatalf("Expected 2 cars, got %d", len(cars))
	}

	if cars[0].Discount != 15 {
		t.Errorf("Expected Car 101 to have 15%% discount, got %d", cars[0].Discount)
	}

	if cars[1].Discount != 0 {
		t.Errorf("Expected Car 102 to have 0%% discount, got %d", cars[1].Discount)
	}
}
