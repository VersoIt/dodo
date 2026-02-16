package domain

import "fmt"

type Coordinates struct {
	Lat float64
	Lng float64
}

func NewCoordinates(lat, lng float64) (Coordinates, error) {
	if lat < -90 || lat > 90 {
		return Coordinates{}, fmt.Errorf("invalid latitude: %f", lat)
	}
	if lng < -180 || lng > 180 {
		return Coordinates{}, fmt.Errorf("invalid longitude: %f", lng)
	}
	return Coordinates{Lat: lat, Lng: lng}, nil
}
