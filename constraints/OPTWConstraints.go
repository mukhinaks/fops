package constraints

import (
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

// OPTWConstraints implements Constraint interface for solving orienteering problem with time windows
type OPTWConstraints struct {
	TimeLimit          int
	StartID            int
	EndID              int
	StartLocation      points.BaseLocation
	EndLocation        points.BaseLocation
	StartTime          int
	DayOfWeek          string
	ForbiddenLocations []int
}

func (f *OPTWConstraints) Init(locs []generic.Point) generic.Constraints {
	start := locs[f.StartID].(points.BaseLocation)
	end := locs[f.EndID].(points.BaseLocation)

	f.StartLocation = start
	f.EndLocation = end

	return f
}

func (f *OPTWConstraints) routeTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		duration = f.StartLocation.Duration
	} else {
		loc := route[orderOfLocations[0]].(points.BaseLocation)
		duration = f.StartLocation.Duration + points.WalkingTime(f.StartLocation, loc)
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			if key == f.StartID || key == f.EndID {
				continue
			}
			walkTime := points.WalkingTime(route[key].(points.BaseLocation), route[orderOfLocations[i+1]].(points.BaseLocation))
			duration += route[key].(points.BaseLocation).Duration + int(walkTime)
		}
		duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation).Duration + f.EndLocation.Duration + points.WalkingTime(route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation), f.EndLocation)

	}

	return duration
}

func (f *OPTWConstraints) TimeUpdate(currentTime int, duration int) int {
	hours := int(currentTime/100) + int(duration/60)
	minutes := (currentTime - int(currentTime/100)*100) + (duration - int(duration/60)*60)
	for minutes >= 60 {
		hours++
		minutes -= 60
	}

	return minutes + hours*100
}

func (f *OPTWConstraints) FinalRouteTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	if route == nil {
		return 0
	}
	for i := 0; i < len(orderOfLocations)-1; i++ {
		key := orderOfLocations[i]
		walkTime := points.WalkingTime(route[key].(points.BaseLocation), route[orderOfLocations[i+1]].(points.BaseLocation))
		duration += route[key].(points.BaseLocation).Duration + int(walkTime)
	}
	duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation).Duration

	return duration
}

func (f *OPTWConstraints) Boundary(route map[int]generic.Point, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)
	if duration > f.TimeLimit {
		return false
	}

	time := f.TimeUpdate(f.StartTime, f.StartLocation.Duration)
	time = f.TimeUpdate(time, points.WalkingTime(f.StartLocation, route[orderOfLocations[0]].(points.BaseLocation)))

	for i := 0; i <= len(orderOfLocations)-1; i++ {
		key := orderOfLocations[i]
		if key == f.StartID || key == f.EndID {
			continue
		}
		location := route[key].(points.BaseLocation)

		openHours, ok := location.OpenHours[f.DayOfWeek]
		if !ok {
			return false
		}

		startWorking := openHours[0]
		endWorking := openHours[1]
		if endWorking < startWorking {
			endWorking += 2400
		}

		if openHours[0] <= time && openHours[1] >= f.TimeUpdate(time, location.Duration) {
			time = f.TimeUpdate(time, location.Duration)
			if i != len(orderOfLocations)-1 {
				walkTime := points.WalkingTime(location, route[orderOfLocations[i+1]].(points.BaseLocation))
				time = f.TimeUpdate(time, walkTime)
			}
		} else {
			return false
		}
	}

	if len(orderOfLocations) > 0 {
		walkTime := points.WalkingTime(f.EndLocation, route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation))
		time = f.TimeUpdate(time, walkTime)
		openHours, ok := f.EndLocation.OpenHours[f.DayOfWeek]
		if !ok {
			return false
		}

		startWorking := openHours[0]
		endWorking := openHours[1]
		if endWorking < startWorking {
			endWorking += 2400
		}

		if openHours[0] <= time && openHours[1] >= f.TimeUpdate(time, f.EndLocation.Duration) {
			return true

		}

		return false
	}
	return true
}

func (f *OPTWConstraints) ReducePoints(route map[int]generic.Point, orderOfLocations []int, locations map[int]generic.Point) map[int]generic.Point {
	return locations
}

func (f *OPTWConstraints) SinglePointConstraints(location generic.Point, id int) bool {
	for _, i := range f.ForbiddenLocations {
		if i == id {
			return false
		}
	}

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

func (f *OPTWConstraints) UpdateConstraint(route map[int]generic.Point, orderOfPoints []int, locations []generic.Point) generic.Constraints {
	return f
}
