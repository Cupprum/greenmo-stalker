package osm

import (
	"fmt"
	"function/geo"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Marker struct {
	Pos   geo.Position
	Color string
	Text  string
	Icon  string
}

func GenerateMap(endpoint string, center geo.Position, markers []Marker, key string) ([]byte, error) {
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("style", "maptiler-3d")
	q.Set("width", "600")
	q.Set("height", "600")
	q.Set("center", fmt.Sprintf("lonlat:%f,%f", center.Lon, center.Lat))
	q.Set("zoom", "14")
	q.Set("apiKey", key)

	var ms []string
	for _, m := range markers {
		ll := fmt.Sprintf("lonlat:%f,%f", m.Pos.Lon, m.Pos.Lat)
		c := fmt.Sprintf("color:%v", m.Color)
		s := fmt.Sprintf("size:%v", "medium")
		add := fmt.Sprintf("text:%v", m.Text)
		if m.Icon != "" {
			add = fmt.Sprintf("icon:%v", m.Icon)
		}
		ms = append(ms, strings.Join([]string{ll, c, s, add}, ";"))
	}

	if len(ms) > 0 {
		q.Set("marker", strings.Join(ms, "|"))
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
