package main

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
		Coordinates [2]float64 `json:"coordinates"`
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
}

func TestGreenMobilityFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cars := []Car{
			{ID: 1, StateOfCharge: 35},
			{ID: 2, StateOfCharge: 50},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cars)
	}))
	defer server.Close()

	// TODO: cars, err := QueryGreenMobility(server.URL, 55.787, 12.519, 2, 40)
	// TODO: assert len(cars) == 1
}
