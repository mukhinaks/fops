package constraints

import (
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type OPConstraints struct {
	TimeLimit     int
	StartID       int
	EndID         int
	StartLocation points.BaseLocation
	EndLocation   points.BaseLocation
}

func (f *OPConstraints) Init(locs []generic.Point) generic.Constraints {
	start := locs[f.StartID].(points.BaseLocation)
	end := locs[f.EndID].(points.BaseLocation)

	f.StartLocation = start
	f.EndLocation = end

	return f
}

func (f *OPConstraints) routeTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		duration = f.StartLocation.Duration + f.EndLocation.Duration + points.WalkingTime(f.StartLocation, f.EndLocation)
	} else {
		loc := route[orderOfLocations[0]].(points.BaseLocation)
		duration = f.StartLocation.Duration + f.EndLocation.Duration + points.WalkingTime(f.StartLocation, loc)
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			if key == f.StartID || key == f.EndID {
				continue
			}
			walkTime := points.WalkingTime(route[key].(points.BaseLocation), route[orderOfLocations[i+1]].(points.BaseLocation))
			duration += route[key].(points.BaseLocation).Duration + int(walkTime)
		}
		duration += points.WalkingTime(f.EndLocation, route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation)) +
			route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation).Duration

	}
	return duration
}

func (f *OPConstraints) FinalRouteTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil || len(orderOfLocations) == 0 {
		return 0
	}
	if len(orderOfLocations) == 1 {
		return route[orderOfLocations[0]].(points.BaseLocation).Duration
	}
	for i := 0; i < len(orderOfLocations)-1; i++ {
		key := orderOfLocations[i]
		walkTime := points.WalkingTime(route[key].(points.BaseLocation), route[orderOfLocations[i+1]].(points.BaseLocation))
		duration += route[key].(points.BaseLocation).Duration + int(walkTime)
	}

	duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation).Duration

	return duration
}

func (f *OPConstraints) Boundary(route map[int]generic.Point, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)

	if duration > f.TimeLimit {
		return false
	}
	return true
}

func (f *OPConstraints) ReducePoints(route map[int]generic.Point, orderOfLocations []int, locations map[int]generic.Point) map[int]generic.Point {
	/*
		currentRouteTime := f.FinalRouteTime(route, orderOfLocations) + f.StartLocation.Duration + f.EndLocation.Duration
		if len(orderOfLocations) != 0 {
			currentRouteTime += points.WalkingTime(f.StartLocation, route[orderOfLocations[0]].(points.BaseLocation))
		}
		filteredLocations := make(map[int]generic.Point)
		for i, location := range locations {
			time := currentRouteTime + location.(points.BaseLocation).Duration + points.WalkingTime(f.EndLocation, location.(points.BaseLocation))
			if len(orderOfLocations) != 0 {
				time += points.WalkingTime(location.(points.BaseLocation), route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation))
			}

			if time <= f.TimeLimit {
				filteredLocations[i] = location
			}

		}
		return filteredLocations
	*/
	return locations
}

func (f *OPConstraints) SinglePointConstraints(location generic.Point, id int) bool {
	if id == f.StartID || id == f.EndID {
		return false
	}

	time := location.(points.BaseLocation).Duration + points.WalkingTime(f.EndLocation, location.(points.BaseLocation)) +
		f.EndLocation.Duration + f.StartLocation.Duration + points.WalkingTime(f.StartLocation, location)

	if time > f.TimeLimit {
		return false
	}

	return true
}

func (f *OPConstraints) ComputeRouteTimeFromSample(locationsID []int, allLocations []generic.Point) int {
	lastLocation := allLocations[len(locationsID)-1].(points.BaseLocation)
	time := lastLocation.Duration

	for i := 0; i < len(locationsID)-1; i++ {
		loc1 := allLocations[locationsID[i]].(points.BaseLocation)
		loc2 := allLocations[locationsID[i+1]].(points.BaseLocation)
		time += loc1.Duration + points.WalkingTime(loc1, loc2)
	}

	return time
}

func (f *OPConstraints) UpdateConstraint(route map[int]generic.Point, orderOfPoints []int, locations []generic.Point) generic.Constraints {
	return f
}
