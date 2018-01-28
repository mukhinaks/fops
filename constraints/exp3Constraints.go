package constraints

import (
	"fmt"

	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/locations"
)

type MultidaysConstraints struct {
	DayTimeLimit int
	DaysNumber   int

	StartID                int
	EndID                  int
	StartLocationDistances map[int]float64
	EndLocationDistance    map[int]float64
	StartEndDistance       float64
	StartLocation          locations.BaseLocation
	EndLocation            locations.BaseLocation

	TimeLimit           []int
	CurrentDay          int
	CurrentInterval     int
	ForbiddenLocations  []int
	CompulsoryLocations []int
	NumberOfInterval    int
}

func (f MultidaysConstraints) Init(locs []generic.Location) generic.Constraints {
	f.StartID = f.CompulsoryLocations[f.NumberOfInterval]
	f.EndID = f.CompulsoryLocations[f.NumberOfInterval+1]
	start := locs[f.CompulsoryLocations[f.NumberOfInterval]].(locations.BaseLocation)
	end := locs[f.CompulsoryLocations[f.NumberOfInterval+1]].(locations.BaseLocation)

	f.StartLocationDistances = make(map[int]float64)
	f.EndLocationDistance = make(map[int]float64)

	for idx, loc := range locs {
		location := loc.(locations.BaseLocation)
		f.StartLocationDistances[idx] = locations.EuclidianDistance(start, location)
		f.EndLocationDistance[idx] = locations.EuclidianDistance(end, location)
	}
	f.StartEndDistance = locations.EuclidianDistance(start, end)
	f.StartLocation = start
	f.EndLocation = end
	return f
}

func (f MultidaysConstraints) routeTime(route map[int]generic.Location, orderOfLocations []int) int {
	duration := 0
	if route[orderOfLocations[0]] == nil || len(orderOfLocations) == 0 {
		duration = f.StartLocation.Duration + f.EndLocation.Duration + locations.WalkingTime(f.StartLocation, f.EndLocation)
	} else {
		loc := route[orderOfLocations[0]].(locations.BaseLocation)
		duration = locations.WalkingTime(f.StartLocation, loc) + f.EndLocation.Duration
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			walkTime := locations.WalkingTime(route[key].(locations.BaseLocation), route[orderOfLocations[i+1]].(locations.BaseLocation))
			duration += route[key].(locations.BaseLocation).Duration + int(walkTime)
		}
		duration += locations.WalkingTime(f.EndLocation, route[orderOfLocations[len(orderOfLocations)-1]].(locations.BaseLocation)) + route[orderOfLocations[len(orderOfLocations)-1]].(locations.BaseLocation).Duration
		if f.NumberOfInterval == 0 {
			duration += f.StartLocation.Duration
		}
	}
	return duration
}

func (f MultidaysConstraints) FinalRouteTime(route map[int]generic.Location, orderOfLocations []int) int {
	if route == nil || len(orderOfLocations) == 0 {
		return 0
	}
	duration := 0
	for i := 0; i < len(orderOfLocations)-1; i++ {
		key := orderOfLocations[i]
		walkTime := locations.WalkingTime(route[key].(locations.BaseLocation), route[orderOfLocations[i+1]].(locations.BaseLocation))
		duration += route[key].(locations.BaseLocation).Duration + int(walkTime)
	}
	duration += route[orderOfLocations[len(orderOfLocations)-1]].(locations.BaseLocation).Duration
	return duration

}

func (f MultidaysConstraints) Boundary(route map[int]generic.Location, orderOfLocations []int) bool {
	if len(orderOfLocations) == 0 {
		return false
	}

	duration := f.routeTime(route, orderOfLocations)

	if duration > f.TimeLimit[f.NumberOfInterval] {
		return false
	}
	return true
}

func (f MultidaysConstraints) LocationConstraints(location generic.Location, id int) bool {

	for _, i := range f.ForbiddenLocations {
		if i == id {
			return false
		}
	}

	if f.StartLocationDistances[id] > f.StartEndDistance && f.EndLocationDistance[id] > f.StartEndDistance {
		return false
	}

	return true
}

func (f MultidaysConstraints) SplitForDays(locationsID []int, allLocations []generic.Location) (map[int][]int, map[int][]int) {
	days := make(map[int][]int)
	currentDay := make([]int, 0)
	currentRoute := make(map[int]generic.Location)
	day := 1

	for i := 0; i < len(locationsID); i++ {
		currentDay = append(currentDay, locationsID[i])
		currentRoute[locationsID[i]] = allLocations[locationsID[i]]
		sum := f.FinalRouteTime(currentRoute, currentDay)
		if sum > f.DayTimeLimit {
			days[day] = currentDay[:len(currentDay)-1]
			currentDay = make([]int, 0)
			currentRoute = make(map[int]generic.Location)
			day++
			currentDay = append(currentDay, locationsID[i])
			currentRoute[locationsID[i]] = allLocations[locationsID[i]]
		}

	}
	days[day] = currentDay

	times := make(map[int][]int)
	for key, day := range days {
		times[key] = f.computeTimeLimits(f.DayTimeLimit, allLocations, day)
	}

	return days, times
}

func (f MultidaysConstraints) computeTimeLimits(routeTimeLimit int, locs []generic.Location, compulsoryLocations []int) []int {
	locationsCount := make([]int, 0)
	minimumTime := make([]int, 0)
	sumLocationsCount := 0
	sumTime := 0

	for i := 0; i < len(compulsoryLocations)-1; i++ {
		keyStart := compulsoryLocations[i]
		keyEnd := compulsoryLocations[i+1]
		distance := locations.EuclidianDistance(locs[keyStart].(locations.BaseLocation), locs[keyEnd].(locations.BaseLocation))

		value := 0
		for id, loc := range locs {
			for _, j := range compulsoryLocations {
				if j == id {
					continue
				}
			}

			if locations.EuclidianDistance(loc.(locations.BaseLocation), locs[keyEnd].(locations.BaseLocation)) <= distance ||
				locations.EuclidianDistance(loc.(locations.BaseLocation), locs[keyStart].(locations.BaseLocation)) <= distance {
				value++
			}
		}
		locationsCount = append(locationsCount, value)
		sumLocationsCount += value

		time :=
			locations.WalkingTime(locs[keyStart].(locations.BaseLocation), locs[keyEnd].(locations.BaseLocation))
		if i == 0 {
			time += locs[keyStart].(locations.BaseLocation).Duration + locs[keyEnd].(locations.BaseLocation).Duration
		} else {
			time += locs[keyEnd].(locations.BaseLocation).Duration
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
