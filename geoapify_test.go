package function_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeoapifyParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("width") != "600" {
			t.Errorf("width wrong")
		}
		if r.URL.Query().Get("height") != "600" {
			t.Errorf("height wrong")
		}
		if r.URL.Query().Get("center") != "lonlat:12.519,55.787" {
			t.Errorf("center wrong")
		}
		if r.URL.Query().Get("zoom") != "14" {
			t.Errorf("zoom wrong")
		}
		if r.URL.Query().Get("apiKey") == "" {
			t.Errorf("apiKey missing")
		}
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{137, 80, 78, 71})
	}))
	defer server.Close()

	// TODO: call generateMap(server.URL, 55.787, 12.519, [], [], apiKey)
}