package constraints

import (
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/locations"
)

type OPConstraints struct {
	TimeLimit              int
	StartID                int
	EndID                  int
	StartLocationDistances map[int]float64
	EndLocationDistance    map[int]float64
	StartEndDistance       float64
	StartLocation          locations.BaseLocation
	EndLocation            locations.BaseLocation
}

func (f *OPConstraints) Init(locs []generic.Location) generic.Constraints {
	start := locs[f.StartID].(locations.BaseLocation)
	end := locs[f.EndID].(locations.BaseLocation)

	f.StartLocationDistances = make(map[int]float64)
	f.EndLocationDistance = make(map[int]float64)

	for idx, location := range locs {
		f.StartLocationDistances[idx] = locations.EuclidianDistance(start, location.(locations.BaseLocation))
		f.EndLocationDistance[idx] = locations.EuclidianDistance(end, location.(locations.BaseLocation))
	}
	f.StartEndDistance = locations.EuclidianDistance(start, end)

	f.StartLocation = start
	f.EndLocation = end

	return f
}

func (f *OPConstraints) routeTime(route map[int]generic.Location, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		duration = f.StartLocation.Duration + f.EndLocation.Duration + locations.WalkingTime(f.StartLocation, f.EndLocation)
	} else {
		loc := route[orderOfLocations[0]].(locations.BaseLocation)
		duration = f.StartLocation.Duration + f.EndLocation.Duration + locations.WalkingTime(f.StartLocation, loc)
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			walkTime := locations.WalkingTime(route[key].(locations.BaseLocation), route[orderOfLocations[i+1]].(locations.BaseLocation))
			duration += route[key].(locations.BaseLocation).Duration + int(walkTime)
		}
		duration += locations.WalkingTime(f.EndLocation, route[orderOfLocations[len(orderOfLocations)-1]].(locations.BaseLocation)) + route[orderOfLocations[len(orderOfLocations)-1]].(locations.BaseLocation).Duration

	}
	return duration
}

func (f *OPConstraints) FinalRouteTime(route map[int]generic.Location, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		return 0
	} else {
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			walkTime := locations.WalkingTime(route[key].(locations.BaseLocation), route[orderOfLocations[i+1]].(locations.BaseLocation))
			duration += route[key].(locations.BaseLocation).Duration + int(walkTime)
		}
		duration += route[orderOfLocations[len(orderOfLocations)-1]].(locations.BaseLocation).Duration

	}
	return duration
}

func (f *OPConstraints) Boundary(route map[int]generic.Location, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)

	if duration > f.TimeLimit {
		return false
	}
	return true
}

func (f *OPConstraints) LocationConstraints(location generic.Location, id int) bool {
	if id == f.StartID || id == f.EndID {
		return false
	}
	if f.StartLocationDistances[id] >= f.StartEndDistance && f.EndLocationDistance[id] >= f.StartEndDistance {
		return false
	}
	return true
}
