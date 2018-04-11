package points

import (
	"fmt"
	"math"

	"github.com/mukhinaks/fops/generic"
)

type Location struct {
	X float64
	Y float64
}

func EuclidianDistance(location1 generic.Point, location2 generic.Point) float64 {
	switch v := location1.(type) {
	case Location:
		loc1 := location1.(Location)
		loc2 := location2.(Location)
		return distanceToPoint(loc1.X, loc1.Y, loc2.X, loc2.Y)
	case BaseLocation:
		loc1 := location1.(BaseLocation)
		loc2 := location2.(BaseLocation)
		return distanceToPoint(loc1.X, loc1.Y, loc2.X, loc2.Y)
	case CityBrandLocation:
		loc1 := location1.(CityBrandLocation)
		loc2 := location2.(CityBrandLocation)
		return distanceToPoint(loc1.X, loc1.Y, loc2.X, loc2.Y)
	default:
		fmt.Println("Unexpected location type", v)
		return -1
	}

}

func EuclidianDistanceToLineSegment(startLocation generic.Point, endLocation generic.Point, newLocation generic.Point) float64 {

	switch v := startLocation.(type) {
	case Location:
		start := startLocation.(Location)
		end := endLocation.(Location)
		location := newLocation.(Location)

		return distanceToLine(start.X, start.Y, end.X, end.Y, location.X, location.Y)

	case BaseLocation:
		start := startLocation.(BaseLocation)
		end := endLocation.(BaseLocation)
		location := newLocation.(BaseLocation)

		return distanceToLine(start.X, start.Y, end.X, end.Y, location.X, location.Y)

	case CityBrandLocation:
		start := startLocation.(CityBrandLocation)
		end := endLocation.(CityBrandLocation)
		location := newLocation.(CityBrandLocation)

		return distanceToLine(start.X, start.Y, end.X, end.Y, location.X, location.Y)
	default:
		fmt.Println("Unexpected location type", v)
		return -1
	}
}

func WalkingTime(location1 generic.Point, location2 generic.Point) int {

	switch v := location1.(type) {
	case Location:
		loc1 := location1.(Location)
		loc2 := location2.(Location)
		return int(distanceToPoint(loc1.X, loc1.Y, loc2.X, loc2.Y) / 66.7)
	case BaseLocation:
		loc1 := location1.(BaseLocation)
		loc2 := location2.(BaseLocation)
		return int(distanceToPoint(loc1.X, loc1.Y, loc2.X, loc2.Y) / 66.7)
	case CityBrandLocation:
		loc1 := location1.(CityBrandLocation)
		loc2 := location2.(CityBrandLocation)
		return int(distanceToPoint(loc1.X, loc1.Y, loc2.X, loc2.Y) / 66.7)
	default:
		fmt.Println("Unexpected location type", v)
		return -1
	}

}

func distanceToPoint(loc1Lat float64, loc1Lng float64, loc2Lat float64, loc2Lng float64) float64 {
	result := math.Sqrt(math.Pow((loc1Lat-loc2Lat), 2) + math.Pow((loc1Lng-loc2Lng), 2))
	return result
}

func scalarDot(loc1Lat float64, loc1Lng float64, loc2Lat float64, loc2Lng float64) float64 {
	result := (loc1Lat * loc2Lat) + (loc1Lng * loc2Lng)
	return result
}

func distanceToLine(loc1Lat float64, loc1Lng float64, loc2Lat float64, loc2Lng float64, newLat float64, newLng float64) float64 {
	vector1Lat := loc1Lat - loc2Lat
	vector1Lng := loc1Lng - loc2Lng
	vector2Lat := loc1Lat - newLat
	vector2Lng := loc1Lng - newLng

	scalarProduct := scalarDot(vector1Lat, vector1Lng, vector2Lat, vector2Lng)
	if scalarProduct <= 0 {
		return distanceToPoint(newLat, newLng, loc1Lat, loc1Lng)
	}

	length := scalarDot(vector1Lat, vector1Lng, vector1Lat, vector1Lng)
	if length <= scalarProduct {
		return distanceToPoint(newLat, newLng, loc2Lat, loc2Lng)
	}

	b := scalarProduct / length
	locLat := loc1Lat + b*vector1Lat
	locLng := loc1Lng + b*vector1Lng

	return distanceToPoint(newLat, newLng, locLat, locLng)
}
