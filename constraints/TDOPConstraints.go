package constraints

import (
	"math/rand"

	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type TDOPConstraints struct {
	TimeLimit int
	StartID   int
	EndID     int

	StartLocation     points.BaseLocation
	EndLocation       points.BaseLocation
	SpeedDistribution map[int][]float64
	StartTime         int
}

func (f *TDOPConstraints) Init(locs []generic.Point) generic.Constraints {
	start := locs[f.StartID].(points.BaseLocation)
	end := locs[f.EndID].(points.BaseLocation)

	f.SpeedDistribution = make(map[int][]float64)

	for idx := range locs {
		f.SpeedDistribution[idx] = make([]float64, 24)
		for i := 0; i < 24; i++ {
			f.SpeedDistribution[idx][i] = rand.NormFloat64()
		}
	}

	f.StartLocation = start
	f.EndLocation = end

	return f
}

func (f *TDOPConstraints) TimeUpdate(currentTime int, duration int) int {
	hours := int(currentTime/100) + int(duration/60)
	minutes := (currentTime - int(currentTime/100)*100) + (duration - int(duration/60)*60)
	for minutes >= 60 {
		hours++
		minutes -= 60
	}
	for hours >= 24 {
		hours -= 24
	}

	return minutes + hours*100
}

func (f *TDOPConstraints) updatedTime(walkingTime int, stockID int, currentTime int) int {
	hour := currentTime / 100
	//fmt.Println(stockID, currentTime, hour)
	newTime := int(float64(walkingTime)*0.05*f.SpeedDistribution[stockID][hour]) + walkingTime
	if newTime <= 0 {
		return walkingTime
	}
	return newTime
}

func (f *TDOPConstraints) routeTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	time := f.StartTime
	if route == nil {
		duration = f.StartLocation.Duration + f.EndLocation.Duration + f.updatedTime(points.WalkingTime(f.StartLocation, f.EndLocation), f.EndID, time)
	} else {
		loc := route[orderOfLocations[0]].(points.BaseLocation)
		duration = f.StartLocation.Duration + f.updatedTime(points.WalkingTime(f.StartLocation, loc), orderOfLocations[0], time)
		time = f.TimeUpdate(f.StartTime, duration)
		for i := 0; i < len(orderOfLocations)-1; i++ {
			key := orderOfLocations[i]
			if key == f.StartID || key == f.EndID {
				continue
			}
			duration += route[key].(points.BaseLocation).Duration
			time = f.TimeUpdate(f.StartTime, duration)
			walkTime := f.updatedTime(points.WalkingTime(route[key].(points.BaseLocation), route[orderOfLocations[i+1]].(points.BaseLocation)), orderOfLocations[i+1], time)
			duration += int(walkTime)
			time = f.TimeUpdate(f.StartTime, duration)
		}
		duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation).Duration
		time = f.TimeUpdate(f.StartTime, duration)
		duration += f.updatedTime(points.WalkingTime(f.EndLocation, route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation)), f.EndID, time) +
			f.EndLocation.Duration

	}
	return duration
}

func (f *TDOPConstraints) FinalRouteTime(route map[int]generic.Point, orderOfLocations []int) int {
	duration := 0
	time := f.StartTime
	if route == nil {
		return 0
	}
	for i := 0; i < len(orderOfLocations)-1; i++ {
		key := orderOfLocations[i]
		duration += route[key].(points.BaseLocation).Duration
		time = f.TimeUpdate(f.StartTime, duration)
		walkTime := f.updatedTime(points.WalkingTime(route[key].(points.BaseLocation), route[orderOfLocations[i+1]].(points.BaseLocation)), orderOfLocations[i+1], time)
		duration += int(walkTime)
		time = f.TimeUpdate(f.StartTime, duration)
	}
	duration += route[orderOfLocations[len(orderOfLocations)-1]].(points.BaseLocation).Duration

	return duration
}

func (f *TDOPConstraints) Boundary(route map[int]generic.Point, orderOfLocations []int) bool {
	duration := f.routeTime(route, orderOfLocations)

	if duration > f.TimeLimit {
		return false
	}
	return true
}

func (f *TDOPConstraints) ReducePoints(route map[int]generic.Point, orderOfLocations []int, locations map[int]generic.Point) map[int]generic.Point {
	return locations
}

func (f *TDOPConstraints) SinglePointConstraints(location generic.Point, id int) bool {
	if id == f.StartID || id == f.EndID {
		return false
	}

	time := f.StartLocation.Duration
	t := f.TimeUpdate(f.StartTime, time)

	time += f.updatedTime(points.WalkingTime(f.StartLocation, location), id, t) + location.(points.BaseLocation).Duration
	t = f.TimeUpdate(f.StartTime, time)

	time += f.updatedTime(points.WalkingTime(f.EndLocation, location.(points.BaseLocation)), f.EndID, t) +
		f.EndLocation.Duration

	if time > f.TimeLimit {
		return false
	}

	return true
}

func (f *TDOPConstraints) UpdateConstraint(route map[int]generic.Point, orderOfPoints []int, locations []generic.Point) generic.Constraints {
	return f
}
