package constraints

import (
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type RestarauntsConstraints struct {
	TimeLimit int
	StartID   int
	EndID     int

	StartLocation      points.CityBrandLocation
	EndLocation        points.CityBrandLocation
	ForbiddenLocations []int
	StartTime          int
	DayOfWeek          int
}

func (f RestarauntsConstraints) Init(locs []generic.Point) generic.Constraints {
	start := locs[f.StartID].(points.CityBrandLocation)
	end := locs[f.EndID].(points.CityBrandLocation)

	f.StartLocation = start
	f.EndLocation = end

	return f
}

func (f RestarauntsConstraints) routeTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		duration = f.StartLocation.Duration + f.EndLocation.Duration + points.WalkingTime(f.StartLocation, f.EndLocation)
	} else {
		loc := route[orderOfLocations[0]].(points.CityBrandLocation)
		duration = f.StartLocation.Duration + f.EndLocation.Duration + points.WalkingTime(f.StartLocation, loc)
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			if key == f.StartID || key == f.EndID {
				continue
			}
			walkTime := points.WalkingTime(route[key].(points.CityBrandLocation), route[orderOfLocations[i+1]].(points.CityBrandLocation))
			duration += route[key].(points.CityBrandLocation).Duration + int(walkTime)
		}
		duration += points.WalkingTime(f.EndLocation, route[orderOfLocations[len(orderOfLocations)-1]].(points.CityBrandLocation)) +
			route[orderOfLocations[len(orderOfLocations)-1]].(points.CityBrandLocation).Duration

	}
	return duration
}

func (f RestarauntsConstraints) FinalRouteTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		return 0
	} else {
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			walkTime := points.WalkingTime(route[key].(points.BaseLocation), route[orderOfLocations[i+1]].(points.BaseLocation))
			duration += route[key].(points.BaseLocation).Duration + int(walkTime)
		}
		duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation).Duration

	}
	return duration
}

func (f RestarauntsConstraints) TimeUpdate(currentTime int, duration int) int {
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

func (f RestarauntsConstraints) Boundary(route map[int]generic.Point, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)
	if duration > f.TimeLimit {
		return false
	}

	time := f.TimeUpdate(f.StartTime, f.StartLocation.Duration)
	time = f.TimeUpdate(time, points.WalkingTime(f.StartLocation, route[orderOfLocations[0]].(points.CityBrandLocation)))

	return true
}

func (f RestarauntsConstraints) ReducePoints(route map[int]generic.Point, orderOfLocations []int, locations map[int]generic.Point) map[int]generic.Point {
	return locations
}

func (f RestarauntsConstraints) SinglePointConstraints(location generic.Point, id int) bool {
	for _, i := range f.ForbiddenLocations {
		if i == id {
			return false
		}
	}

	if id == f.StartID || id == f.EndID {
		return false
	}

	loc := location.(points.CityBrandLocation)
	for _, category := range loc.Categories {
		if category == "Restaurant" {
			return true
		}
	}
	return false
}

func (f RestarauntsConstraints) UpdateConstraint(route map[int]generic.Point, orderOfPoints []int, locations []generic.Point) generic.Constraints {
	return f
}
