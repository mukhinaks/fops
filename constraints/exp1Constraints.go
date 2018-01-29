package constraints

import (
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type OPConstraints struct {
	TimeLimit              int
	StartID                int
	EndID                  int
	StartLocationDistances map[int]float64
	EndLocationDistance    map[int]float64
	StartEndDistance       float64
	StartLocation          points.BaseLocation
	EndLocation            points.BaseLocation
}

func (f *OPConstraints) Init(locs []generic.Point) generic.Constraints {
	start := locs[f.StartID].(points.BaseLocation)
	end := locs[f.EndID].(points.BaseLocation)

	f.StartLocationDistances = make(map[int]float64)
	f.EndLocationDistance = make(map[int]float64)

	for idx, location := range locs {
		f.StartLocationDistances[idx] = points.EuclidianDistance(start, location.(points.BaseLocation))
		f.EndLocationDistance[idx] = points.EuclidianDistance(end, location.(points.BaseLocation))
	}
	f.StartEndDistance = points.EuclidianDistance(start, end)

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

func (f *OPConstraints) Boundary(route map[int]generic.Point, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)

	if duration > f.TimeLimit {
		return false
	}
	return true
}

func (f *OPConstraints) SinglePointConstraints(location generic.Point, id int) bool {
	if id == f.StartID || id == f.EndID {
		return false
	}
	if f.StartLocationDistances[id] >= f.StartEndDistance && f.EndLocationDistance[id] >= f.StartEndDistance {
		return false
	}
	return true
}
