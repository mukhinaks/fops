package constraints

import (
	"fmt"

	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type EnrichmentConstraints struct {
	TimeLimit              []int
	NumberOfInterval       int
	StartID                int
	EndID                  int
	StartLocationDistances map[int]float64
	EndLocationDistance    map[int]float64
	StartEndDistance       float64
	StartLocation          points.BaseLocation
	EndLocation            points.BaseLocation
	ForbiddenLocations     []int
	CompulsoryLocations    []int
	RouteTimeLimit         int
}

func (f EnrichmentConstraints) Init(locs []generic.Point) generic.Constraints {
	f.StartID = f.CompulsoryLocations[f.NumberOfInterval]
	f.EndID = f.CompulsoryLocations[f.NumberOfInterval+1]
	start := locs[f.CompulsoryLocations[f.NumberOfInterval]].(points.BaseLocation)
	end := locs[f.CompulsoryLocations[f.NumberOfInterval+1]].(points.BaseLocation)

	f.StartLocationDistances = make(map[int]float64)
	f.EndLocationDistance = make(map[int]float64)

	for idx, loc := range locs {
		location := loc.(points.BaseLocation)
		f.StartLocationDistances[idx] = points.EuclidianDistance(start, location)
		f.EndLocationDistance[idx] = points.EuclidianDistance(end, location)
	}
	f.StartEndDistance = points.EuclidianDistance(start, end)
	f.StartLocation = start
	f.EndLocation = end

	f.TimeLimit = f.computeTimeLimits(f.RouteTimeLimit, locs)
	return f
}

func (f EnrichmentConstraints) computeTimeLimits(routeTimeLimit int, locs []generic.Point) []int {
	locationsCount := make([]int, 0)
	minimumTime := make([]int, 0)
	sumLocationsCount := 0
	sumTime := 0

	for i := 0; i < len(f.CompulsoryLocations)-1; i++ {
		keyStart := f.CompulsoryLocations[i]
		keyEnd := f.CompulsoryLocations[i+1]
		distance := points.EuclidianDistance(locs[keyStart].(points.BaseLocation), locs[keyEnd].(points.BaseLocation))

		value := 0
		for id, loc := range locs {
			for _, j := range f.CompulsoryLocations {
				if j == id {
					continue
				}
			}

			if points.EuclidianDistance(loc.(points.BaseLocation), locs[keyEnd].(points.BaseLocation)) <= distance ||
				points.EuclidianDistance(loc.(points.BaseLocation), locs[keyStart].(points.BaseLocation)) <= distance {
				value++
			}
		}
		locationsCount = append(locationsCount, value)
		sumLocationsCount += value

		time :=
			points.WalkingTime(locs[keyStart].(points.BaseLocation), locs[keyEnd].(points.BaseLocation))
		if i == 0 {
			time += locs[keyStart].(points.BaseLocation).Duration + locs[keyEnd].(points.BaseLocation).Duration
		} else {
			time += locs[keyEnd].(points.BaseLocation).Duration

		}
		sumTime += time
		minimumTime = append(minimumTime, time)
	}
	timeLimits := make([]int, 0)
	timeSum := 0
	reserve := routeTimeLimit - sumTime
	for i, count := range locationsCount {
		timeInterval := minimumTime[i]
		if reserve >= 0 {
			timeInterval += int(float64(reserve) * (float64(count) / float64(sumLocationsCount)))
		} else {
			fmt.Println("Warning! Too much locations for this day limit. Time limit for day should be expanded.")
			reserve = 0
		}

		timeLimits = append(timeLimits, timeInterval)
		timeSum += timeInterval
	}

	return timeLimits
}

func (f EnrichmentConstraints) routeTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		duration = f.StartLocation.Duration + f.EndLocation.Duration + points.WalkingTime(f.StartLocation, f.EndLocation)
	} else {
		loc := route[orderOfLocations[0]].(points.BaseLocation)
		duration = f.EndLocation.Duration + points.WalkingTime(f.StartLocation, loc)
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
		if f.NumberOfInterval == 0 {
			duration += f.StartLocation.Duration
		}

	}
	return duration
}

func (f EnrichmentConstraints) FinalRouteTime(route map[int]generic.Point, orderOfLocations []int) int {
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

func (f EnrichmentConstraints) Boundary(route map[int]generic.Point, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)

	if duration > f.TimeLimit[f.NumberOfInterval] {
		return false
	}
	return true
}

func (f EnrichmentConstraints) ReducePoints(route map[int]generic.Point, orderOfLocations []int, locations map[int]generic.Point) map[int]generic.Point {
	return locations
}

func (f EnrichmentConstraints) SinglePointConstraints(location generic.Point, id int) bool {
	for _, i := range f.ForbiddenLocations {
		if i == id {
			return false
		}
	}

	for _, i := range f.CompulsoryLocations {
		if i == id {
			return false
		}
	}
	/*
		if f.StartLocationDistances[id] > f.StartEndDistance && f.EndLocationDistance[id] > f.StartEndDistance {
			return false
		}
	*/

	return true
}

func (f EnrichmentConstraints) UpdateConstraint(route map[int]generic.Point, orderOfPoints []int, locations []generic.Point) generic.Constraints {
	return f
}
