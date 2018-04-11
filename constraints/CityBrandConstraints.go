package constraints

import (
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

// CityBrandConstraints implements Constraint interface for solving orienteering problem with time windows
type CityBrandConstraints struct {
	TimeLimit          int
	StartID            int
	StartLocation      points.CityBrandLocation
	StartTime          int
	DayOfWeek          int
	ForbiddenLocations []int
}

func (f *CityBrandConstraints) Init(locs []generic.Point) generic.Constraints {
	start := locs[f.StartID].(points.CityBrandLocation)

	f.StartLocation = start

	return f
}

func (f *CityBrandConstraints) routeTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		duration = f.StartLocation.Duration
	} else {
		loc := route[orderOfLocations[0]].(points.CityBrandLocation)
		duration = f.StartLocation.Duration + points.WalkingTime(f.StartLocation, loc)
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			if key == f.StartID {
				continue
			}
			walkTime := points.WalkingTime(route[key].(points.CityBrandLocation), route[orderOfLocations[i+1]].(points.CityBrandLocation))
			duration += route[key].(points.CityBrandLocation).Duration + int(walkTime)
		}
		duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.CityBrandLocation).Duration

	}

	return duration
}

func (f *CityBrandConstraints) timeUpdate(currentTime int, duration int) int {
	hours := currentTime/100 + duration/60
	minutes := currentTime - hours + (duration - duration/60)
	if minutes > 60 {
		hours++
		minutes -= 60
	}
	/*
		if hours > 24 {
			hours -= 24

			return (hours*100 + minutes)
		}
	*/
	return (hours*100 + minutes)
}

func (f *CityBrandConstraints) FinalRouteTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		return 0
	}
	for i := 0; i < len(orderOfLocations)-1; i++ {
		key := orderOfLocations[i]
		walkTime := points.WalkingTime(route[key].(points.CityBrandLocation), route[orderOfLocations[i+1]].(points.CityBrandLocation))
		duration += route[key].(points.CityBrandLocation).Duration + int(walkTime)
	}
	duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.CityBrandLocation).Duration

	return duration
}

func (f *CityBrandConstraints) Boundary(route map[int]generic.Point, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)
	if duration > f.TimeLimit {
		return false
	}

	time := f.timeUpdate(f.StartTime, f.StartLocation.Duration)
	time = f.timeUpdate(time, points.WalkingTime(f.StartLocation, route[orderOfLocations[0]].(points.CityBrandLocation)))
	return true
}

func (f *CityBrandConstraints) ReducePoints(route map[int]generic.Point, orderOfLocations []int, locations map[int]generic.Point) map[int]generic.Point {
	return locations
}

func (f *CityBrandConstraints) SinglePointConstraints(location generic.Point, id int) bool {
	for _, i := range f.ForbiddenLocations {
		if i == id {
			return false
		}
	}

	if id == f.StartID {
		return false
	}

	loc := location.(points.CityBrandLocation)
	for _, category := range loc.Categories {
		if category == "Restaurant" {
			return false
		}
	}
	return true
}

func (f *CityBrandConstraints) UpdateConstraint(route map[int]generic.Point, orderOfPoints []int, locations []generic.Point) generic.Constraints {

	return f
}
