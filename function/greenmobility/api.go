package greenmobility

import (
	"encoding/json"
	"fmt"
	"function/geo"
	"net/http"
	"net/url"
)

func Query(endpoint string, center geo.Position, radius float64, fuel int) ([]geo.Position, error) {
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("lat", fmt.Sprintf("%f", center.Lat))
	q.Set("lng", fmt.Sprintf("%f", center.Lon))
	q.Set("rad", fmt.Sprintf("%.1f", radius))
	q.Set("excludeStationedVehicles", "true")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var cars []struct {
		SOC int `json:"stateOfCharge"`
		Pos struct {
			Coords [2]float64 `json:"coordinates"`
		} `json:"position"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cars); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var res []geo.Position
	for _, c := range cars {
		if c.SOC <= fuel {
			res = append(res, geo.Position{Lat: c.Pos.Coords[1], Lon: c.Pos.Coords[0]})
		}
	}
	return res, nil
}
