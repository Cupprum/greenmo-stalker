package spirii

import (
	"encoding/json"
	"fmt"
	"function/geo"
	"net/http"
	"net/url"
	"strings"
)

func Query(endpoint string, nw, se geo.Position) ([]geo.Position, error) {
	// Spirii uses NE/SW
	ne, sw := geo.Position{Lat: nw.Lat, Lon: se.Lon}, geo.Position{Lat: se.Lat, Lon: nw.Lon}

	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("neCoordinates", fmt.Sprintf("%f, %f", ne.Lat, ne.Lon))
	q.Set("swCoordinates", fmt.Sprintf("%f, %f", sw.Lat, sw.Lon))
	q.Set("includeOccupied", "true")
	q.Set("includeOutOfService", "true")
	q.Set("includeRoaming", "true")
	q.Set("onlyIncludeFavourite", "false")
	q.Set("zoomLevel", "22")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	type Charger struct {
		Props struct {
			Id    string `json:"id"`
			Avail int    `json:"availableConnectors"`
		} `json:"properties"`
		Geom struct {
			Coords [2]float64 `json:"coordinates"`
		} `json:"geometry"`
	}
	var data []Charger
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var res []geo.Position
	for _, c := range data {
		// Only include chargers which are available and do not belong to Clever
		if c.Props.Avail > 0 && !strings.Contains(c.Props.Id, "CLE") {
			res = append(res, geo.Position{Lat: c.Geom.Coords[1], Lon: c.Geom.Coords[0]})
		}
	}
	return res, nil
}
