package score

import (
	//	"math"
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type OPScore struct {
	TimeLimit              int
	StartID                int
	EndID                  int
	InstagramVisitorsMax   float64
	TripAdvisorVisitorsMax float64
	StartLocationDistances map[int]float64
	EndLocationDistance    map[int]float64
	StartEndDistance       float64
	StartLocation          generic.Point
	EndLocation            generic.Point
}

func (f OPScore) Init(locs []generic.Point) generic.Score {
	instagramVisitorsMax := 0.0
	tripAdvisorVisitorsMax := 0.0
	start := locs[f.StartID].(points.BaseLocation)
	end := locs[f.EndID].(points.BaseLocation)

	f.StartLocationDistances = make(map[int]float64)
	f.EndLocationDistance = make(map[int]float64)

	for idx, loc := range locs {
		location := loc.(points.BaseLocation)
		if location.InstagramVisitorsNumber > instagramVisitorsMax {
			instagramVisitorsMax = location.InstagramVisitorsNumber
		}
		if location.TripAdvisorReviewsNumber > tripAdvisorVisitorsMax {
			tripAdvisorVisitorsMax = location.TripAdvisorReviewsNumber
		}
		f.StartLocationDistances[idx] = points.EuclidianDistance(start, location)
		f.EndLocationDistance[idx] = points.EuclidianDistance(end, location)
	}
	f.StartEndDistance = points.EuclidianDistance(start, end)

	f.InstagramVisitorsMax = instagramVisitorsMax
	f.TripAdvisorVisitorsMax = tripAdvisorVisitorsMax
	f.StartLocation = start
	f.EndLocation = end
	return f
}

func (f OPScore) SinglePointScore(route map[int]generic.Point, orderOfLocations []int,
	location generic.Point, id int) float64 {
	distanceCoefficient := 0.0

	if len(orderOfLocations) != 0 && id != f.EndID && id != f.StartID {
		idLastLocation := orderOfLocations[len(orderOfLocations)-1]
		loc1 := route[idLastLocation].(points.BaseLocation)
		loc2 := location.(points.BaseLocation)
		distanceCoefficient = f.EndLocationDistance[idLastLocation] / (points.EuclidianDistance(loc1, loc2) +
			f.EndLocationDistance[id])
		if points.WalkingTime(loc1, loc2) > int(float64(f.EndLocationDistance[idLastLocation])/66.7) {
			distanceCoefficient = 0.0
		}
	} else {
		distanceCoefficient = f.StartEndDistance / (f.StartLocationDistances[id] + f.EndLocationDistance[id])
	}
	loc := location.(points.BaseLocation)
	score := loc.OfficialGuide + loc.FoursquareRating/10.0 +
		(loc.TripAdvisorRating/5.0)*(loc.TripAdvisorReviewsNumber/f.TripAdvisorVisitorsMax) +
		loc.InstagramVisitorsNumber/f.InstagramVisitorsMax +
		distanceCoefficient

	return score
}

func (f OPScore) RouteScore(route map[int]generic.Point, orderOfLocations []int) float64 {
	_, ok := route[f.StartID]
	if !ok {
		route[f.StartID] = f.StartLocation
	}
	_, ok = route[f.EndID]
	if !ok {
		route[f.EndID] = f.EndLocation
	}

	routeScore := 0.0
	for i, key := range orderOfLocations {
		routeScore += f.SinglePointScore(route, orderOfLocations[:i+1], route[key], key)
	}
	return routeScore
}

func (f OPScore) ComputeRouteTimeFromSample(locationsID []int, allLocations []generic.Point) int {
	lastLocation := allLocations[len(locationsID)-1].(points.BaseLocation)
	time := lastLocation.Duration

	for i := 0; i < len(locationsID)-1; i++ {
		loc1 := allLocations[locationsID[i]].(points.BaseLocation)
		loc2 := allLocations[locationsID[i+1]].(points.BaseLocation)
		time += loc1.Duration + points.WalkingTime(loc1, loc2)
	}

	return time
}
