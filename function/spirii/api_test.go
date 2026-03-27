package spirii_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Charger struct {
	Properties struct {
		ID                  string `json:"id"`
		AvailableConnectors int    `json:"availableConnectors"`
	} `json:"properties"`
	Geometry struct {
		Coordinates [2]float64 `json:"coordinates"` // [lon, lat]
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

	// TODO: Call implementation
	// _, err := QuerySpirii(server.URL, 55.794430, 12.527933, 55.779566, 12.511368)
	// if err != nil {
	// 	t.Errorf("Expected no error, got %v", err)
	// }
}

func TestSpiriiFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chargers := []Charger{
			{
				Properties: struct {
					ID                  string `json:"id"`
					AvailableConnectors int    `json:"availableConnectors"`
				}{ID: "charger_001", AvailableConnectors: 2}, // Should pass (> 0)
				Geometry: struct {
					Coordinates [2]float64 `json:"coordinates"`
				}{Coordinates: [2]float64{12.520, 55.787}},
			},
			{
				Properties: struct {
					ID                  string `json:"id"`
					AvailableConnectors int    `json:"availableConnectors"`
				}{ID: "charger_002", AvailableConnectors: 0}, // Should fail (0)
				Geometry: struct {
					Coordinates [2]float64 `json:"coordinates"`
				}{Coordinates: [2]float64{12.525, 55.785}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"features": chargers})
	}))
	defer server.Close()

	// TODO: Processed result call
	// chargers, err := QuerySpirii(server.URL, 55.794430, 12.527933, 55.779566, 12.511368)

	// TODO: Assertions
	// if err != nil {
	// 	t.Fatalf("Expected no error, got %v", err)
	// }
	// if len(chargers) != 1 {
	// 	t.Fatalf("Expected 1 charger, got %d", len(chargers))
	// }

	// TODO: Verify coordinate mapping from [lon, lat] to {lat, lon}
	// if chargers[0].Lat != 55.787 {
	// 	t.Errorf("Expected latitude 55.787, got %f", chargers[0].Lat)
	// }
	// if chargers[0].Lon != 12.520 {
	// 	t.Errorf("Expected longitude 12.520, got %f", chargers[0].Lon)
	// }
}
