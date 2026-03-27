package openstreetmaps

import (
	"fmt"
	"function/geo"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GenerateMap(endpoint string, center geo.Position, cars, chargers []geo.Position, key string) ([]byte, error) {
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("style", "maptiler-3d")
	q.Set("width", "600")
	q.Set("height", "600")
	q.Set("center", fmt.Sprintf("lonlat:%f,%f", center.Lon, center.Lat))
	q.Set("zoom", "14")
	q.Set("apiKey", key)

	var m []string
	for _, p := range cars {
		m = append(m, fmt.Sprintf("lonlat:%f,%f;color:#3ea635;size:medium", p.Lon, p.Lat))
	}
	for _, p := range chargers {
		m = append(m, fmt.Sprintf("lonlat:%f,%f;color:#f30e0e;size:medium", p.Lon, p.Lat))
	}

	if len(m) > 0 {
		q.Set("marker", strings.Join(m, "|"))
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
