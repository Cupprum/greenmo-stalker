package spirii_test

import (
	"encoding/json"
	"function/geo"
	"function/spirii"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpiriiParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotNe := r.URL.Query().Get("neCoordinates")
		wantNe := "55.794430, 12.527933"
		if gotNe != wantNe {
			t.Errorf("neCoordinates wrong, got %q, want %q", gotNe, wantNe)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	ne := geo.Position{Lat: 55.794430, Lon: 12.527933}
	sw := geo.Position{Lat: 55.779566, Lon: 12.511368}
	_, err := spirii.Query(server.URL, ne, sw)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpiriiFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := []map[string]interface{}{
			{
				"properties": map[string]interface{}{"availableConnectors": 2},
				"geometry":   map[string]interface{}{"coordinates": [2]float64{12.520, 55.787}},
			},
			{
				"properties": map[string]interface{}{"availableConnectors": 0},
				"geometry":   map[string]interface{}{"coordinates": [2]float64{12.525, 55.785}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	ne := geo.Position{Lat: 55.79, Lon: 12.52}
	sw := geo.Position{Lat: 55.77, Lon: 12.51}
	chargers, err := spirii.Query(server.URL, ne, sw)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(chargers) != 1 {
		t.Fatalf("Expected 1 charger after filtering, got %d", len(chargers))
	}
	if chargers[0].Lat != 55.787 || chargers[0].Lon != 12.520 {
		t.Errorf("Coordinate mapping failed: got Lat %f Lon %f", chargers[0].Lat, chargers[0].Lon)
	}
}
