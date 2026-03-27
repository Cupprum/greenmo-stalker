package greenmobility_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Car struct {
	ID            int `json:"id"`
	StateOfCharge int `json:"stateOfCharge"`
	Position      struct {
		Coordinates [2]float64 `json:"coordinates"` // [lon, lat]
	} `json:"position"`
}

func TestGreenMobilityParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("lat") != "55.787" {
			t.Errorf("lat param wrong")
		}
		if r.URL.Query().Get("lng") != "12.519" {
			t.Errorf("lng param wrong")
		}
		if r.URL.Query().Get("rad") != "2" {
			t.Errorf("rad param wrong")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Car{})
	}))
	defer server.Close()

	// TODO: call QueryGreenMobility(server.URL, 55.787, 12.519, 2, 40)
	// if err != nil {
	// 	t.Fatalf("Expected no error, got %v", err)
	// }
}

func TestGreenMobilityFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cars := []Car{
			{
				ID:            1,
				StateOfCharge: 35, // Should pass filter (<= 40)
				Position: struct {
					Coordinates [2]float64 `json:"coordinates"`
				}{Coordinates: [2]float64{12.515, 55.790}},
			},
			{
				ID:            2,
				StateOfCharge: 50, // Should fail filter (> 40)
				Position: struct {
					Coordinates [2]float64 `json:"coordinates"`
				}{Coordinates: [2]float64{12.525, 55.785}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cars)
	}))
	defer server.Close()

	// TODO: Processed result call
	// cars, err := QueryGreenMobility(server.URL, 55.787, 12.519, 2, 40)

	// TODO: Assertions
	// if err != nil {
	// 	t.Fatalf("Expected no error, got %v", err)
	// }
	// if len(cars) != 1 {
	// 	t.Fatalf("Expected 1 car, got %d", len(cars))
	// }

	// TODO: Verify coordinate mapping from [lon, lat] to {lat, lon}
	// if cars[0].Lat != 55.790 {
	// 	t.Errorf("Expected latitude 55.790, got %f", cars[0].Lat)
	// }
	// if cars[0].Lon != 12.515 {
	// 	t.Errorf("Expected longitude 12.515, got %f", cars[0].Lon)
	// }
}
