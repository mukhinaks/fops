package constraints

import (
	"fmt"
	"math/rand"

	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type EROPFPConstraints struct {
	TimeLimit              int
	StartID                int
	EndID                  int
	StartLocationDistances map[int]float64
	EndLocationDistance    map[int]float64
	StartEndDistance       float64
	StartLocation          points.BaseLocation
	EndLocation            points.BaseLocation
}

func (f *EROPFPConstraints) Init(locs []generic.Point) generic.Constraints {
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

func (f *EROPFPConstraints) routeTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0

	if route == nil {
		duration = f.StartLocation.Duration + f.EndLocation.Duration + points.WalkingTime(f.StartLocation, f.EndLocation)
	} else {
		loc := route[orderOfLocations[0]].(points.BaseLocation)
		duration = f.StartLocation.Duration + points.WalkingTime(f.StartLocation, loc)

		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			if key == f.StartID || key == f.EndID {
				continue
			}
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(route)
					fmt.Println(orderOfLocations)
				}
			}()

			walkTime := points.WalkingTime(route[key].(points.BaseLocation), route[orderOfLocations[i+1]].(points.BaseLocation))
			duration += route[key].(points.BaseLocation).Duration + int(walkTime)

		}

		duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation).Duration

		duration += points.WalkingTime(f.EndLocation, route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation)) +
			f.EndLocation.Duration

	}
	return duration
}

func (f *EROPFPConstraints) FinalRouteTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0

	if route == nil {
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

func (f *EROPFPConstraints) Boundary(route map[int]generic.Point, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)

	if duration > f.TimeLimit {
		return false
	}
	return true
}

func (f *EROPFPConstraints) ReducePoints(route map[int]generic.Point, orderOfLocations []int, locations map[int]generic.Point) map[int]generic.Point {
	filteredLocations := make(map[int]generic.Point)

	if route == nil || len(orderOfLocations) < 2 {
		for key, location := range locations {
			if f.SinglePointConstraints(location, key) {
				filteredLocations[key] = location
			}
		}
		return locations
	}

	latestLocation := route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation)

	distance := points.EuclidianDistance(latestLocation, f.EndLocation)

	for key, location := range locations {

		if points.EuclidianDistance(location, latestLocation) <= distance && f.EndLocationDistance[key] <= distance {
			filteredLocations[key] = location

			candidate := make(map[int]generic.Point)
			for index, value := range route {
				candidate[index] = value
			}
			candidate[key] = location
			if f.Boundary(candidate, append(orderOfLocations, key)) {
				filteredLocations[key] = location
			}

		}
	}

	idx := len(locations)
	sizeFiltered := len(filteredLocations)
	numberOfEvents := rand.Int()
	for len(filteredLocations) < sizeFiltered+numberOfEvents {
		event := latestLocation
		if f.SinglePointConstraints(event, idx) {
			filteredLocations[idx] = event
			idx++
		}
	}

	return filteredLocations
}

func (f *EROPFPConstraints) SinglePointConstraints(location generic.Point, id int) bool {

	if id == f.StartID || id == f.EndID {
		return false
	}

	if f.StartLocationDistances[id] >= f.StartEndDistance && f.EndLocationDistance[id] >= f.StartEndDistance {
		return false
	}

	time := location.(points.BaseLocation).Duration + points.WalkingTime(f.EndLocation, location.(points.BaseLocation)) +
		f.EndLocation.Duration + f.StartLocation.Duration + points.WalkingTime(f.StartLocation, location)

	if time > f.TimeLimit {
		return false
	}

	return true
}

func (f *EROPFPConstraints) UpdateConstraint(route map[int]generic.Point, orderOfPoints []int, locations []generic.Point) generic.Constraints {
	f.StartID = orderOfPoints[0]
	f.EndID = orderOfPoints[len(orderOfPoints)-1]

	return f.Init(locations)
}
