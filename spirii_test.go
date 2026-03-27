package function_test

import (
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"testing"
)

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

	// TODO: call querySpirii(server.URL, 55.794430, 12.527933, 55.779566, 12.511368)
}