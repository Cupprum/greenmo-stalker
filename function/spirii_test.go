package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Charger struct {
	Properties struct {
		AvailableConnectors int `json:"availableConnectors"`
	} `json:"properties"`
	Geometry struct {
		Coordinates [2]float64 `json:"coordinates"`
	} `json:"geometry"`
}

func TestSpiriiParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("neCoordinates") != "55.794430, 12.527933" {
			t.Errorf("neCoordinates wrong")
		}
		if r.URL.Query().Get("swCoordinates") != "55.779566, 12.511368" {
			t.Errorf("swCoordinates wrong")
		}
		if r.URL.Query().Get("zoomLevel") != "22" {
			t.Errorf("zoomLevel wrong")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"features": []Charger{}})
	}))
	defer server.Close()

	// TODO: chargers, err := QuerySpirii(server.URL, 55.794430, 12.527933, 55.779566, 12.511368)
}

func TestSpiriiFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chargers := []Charger{
			{Properties: struct {
				AvailableConnectors int `json:"availableConnectors"`
			}{AvailableConnectors: 2}},
			{Properties: struct {
				AvailableConnectors int `json:"availableConnectors"`
			}{AvailableConnectors: 0}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"features": chargers})
	}))
	defer server.Close()

	// TODO: chargers, err := QuerySpirii(server.URL, 55.794430, 12.527933, 55.779566, 12.511368)
	// TODO: assert len(chargers) == 1
}
