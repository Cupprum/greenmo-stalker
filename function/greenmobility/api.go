package greenmobility

import (
	"encoding/json"
	"fmt"
	"function/geo"
	"net/http"
	"net/url"
	"sort"
)

type Car struct {
	Pos        geo.Position
	Charge     int
	Discounted bool
}

func Query(endpoint string, nw, se geo.Position, fuel int) ([]Car, error) {
	center := geo.Position{Lat: (nw.Lat + se.Lat) / 2, Lon: (nw.Lon + se.Lon) / 2}
	radius := geo.Distance(nw, se) / 2

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

	var rawCars []struct {
		ID  int `json:"id"`
		SOC int `json:"stateOfCharge"`
		Pos struct {
			Coords [2]float64 `json:"coordinates"`
		} `json:"position"`
		Benefit string `json:"benefit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawCars); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var cars []Car
	for _, c := range rawCars {
		pos := geo.Position{Lat: c.Pos.Coords[1], Lon: c.Pos.Coords[0]}

		// Greenmo thinks in circles, we think in squares
		inBox := nw.Lon < pos.Lon && pos.Lon < se.Lon && se.Lat < pos.Lat && pos.Lat < nw.Lat
		if c.SOC <= fuel && inBox {
			cars = append(cars, Car{Pos: pos, Charge: c.SOC, Discounted: c.Benefit == "DISCOUNTED"})
		}
	}
	sort.Slice(cars, func(i, j int) bool {
		return cars[i].Charge > cars[j].Charge
	})
	return cars, nil
}
