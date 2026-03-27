package geo

import "math"

type Position struct{ Lat, Lon float64 }

// Distance returns the distance between two points in km (Haversine approximation)
func Distance(p1, p2 Position) float64 {
	const R = 6371.0 // Earth radius in km
	dLat := (p2.Lat - p1.Lat) * math.Pi / 180
	dLon := (p2.Lon - p1.Lon) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(p1.Lat*math.Pi/180)*math.Cos(p2.Lat*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
