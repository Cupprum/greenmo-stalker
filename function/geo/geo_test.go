package geo

import (
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	p1 := Position{Lat: 55.6727, Lon: 12.5645} // Copenhagen Central
	p2 := Position{Lat: 55.6814, Lon: 12.5758} // Round Tower

	got := Distance(p1, p2)
	want := 1.21 // approximately 1.21 km

	// We allow a small margin of error (0.05 km) for spherical approximation
	if math.Abs(got-want) > 0.05 {
		t.Errorf("Distance() = %f; want %f (within 0.05 margin)", got, want)
	}
}

func TestZeroDistance(t *testing.T) {
	p := Position{Lat: 55.0, Lon: 12.0}
	if d := Distance(p, p); d != 0 {
		t.Errorf("Distance to self should be 0, got %f", d)
	}
}
