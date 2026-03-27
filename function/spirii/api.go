package spirii

import (
	"encoding/json"
	"fmt"
	"function/geo"
	"net/http"
	"net/url"
)

func Query(endpoint string, ne, sw geo.Position) ([]geo.Position, error) {
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("neCoordinates", fmt.Sprintf("%f, %f", ne.Lat, ne.Lon))
	q.Set("swCoordinates", fmt.Sprintf("%f, %f", sw.Lat, sw.Lon))
	q.Set("zoomLevel", "22")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	var data struct {
		Features []struct {
			Props struct {
				Avail int `json:"availableConnectors"`
			} `json:"properties"`
			Geom struct {
				Coords [2]float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var res []geo.Position
	for _, f := range data.Features {
		if f.Props.Avail > 0 {
			res = append(res, geo.Position{Lat: f.Geom.Coords[1], Lon: f.Geom.Coords[0]})
		}
	}
	return res, nil
}
