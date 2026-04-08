package greenmobility

import (
	"encoding/json"
	"fmt"
	"function/geo"
	"net/http"
	"net/url"
)

type Car struct {
	Pos      geo.Position
	Charge   int
	Discount int
}

func getCar(endpoint string, id int) (Car, error) {
	u, _ := url.Parse(fmt.Sprintf("%v/%v", endpoint, id))

	resp, err := http.Get(u.String())
	if err != nil {
		return Car{}, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Car{}, fmt.Errorf("status %d", resp.StatusCode)
	}

	var c struct {
		Discount struct {
			Percengate float64 `json:"discountPercentage"`
		} `json:"discount"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return Car{}, fmt.Errorf("decode: %w", err)
	}
	return Car{Discount: int(c.Discount.Percengate)}, nil
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

	var cars []struct {
		ID  int `json:"id"`
		SOC int `json:"stateOfCharge"`
		Pos struct {
			Coords [2]float64 `json:"coordinates"`
		} `json:"position"`
		Benefit string `json:"benefit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cars); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var res []Car
	for _, c := range cars {
		pos := geo.Position{Lat: c.Pos.Coords[1], Lon: c.Pos.Coords[0]}
		// Greenmo thinks in circles, we think in squares
		inBox := nw.Lon < pos.Lon && pos.Lon < se.Lon && se.Lat < pos.Lat && pos.Lat < nw.Lat
		if c.SOC <= fuel && inBox {
			discount := 0
			if c.Benefit == "DISCOUNTED" {
				details, err := getCar(endpoint, c.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to get car: %w", err)
				}
				discount = details.Discount
			}
			res = append(res, Car{Pos: pos, Charge: c.SOC, Discount: discount})
		}
	}
	return res, nil
}
