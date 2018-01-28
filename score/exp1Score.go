package score

import (
	//	"math"
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/locations"
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
	StartLocation          generic.Location
	EndLocation            generic.Location
}

func (f OPScore) Init(locs []generic.Location) generic.Score {
	instagramVisitorsMax := 0.0
	tripAdvisorVisitorsMax := 0.0
	start := locs[f.StartID].(locations.BaseLocation)
	end := locs[f.EndID].(locations.BaseLocation)

	f.StartLocationDistances = make(map[int]float64)
	f.EndLocationDistance = make(map[int]float64)

	for idx, loc := range locs {
		location := loc.(locations.BaseLocation)
		if location.InstagramVisitorsNumber > instagramVisitorsMax {
			instagramVisitorsMax = location.InstagramVisitorsNumber
		}
		if location.TripAdvisorReviewsNumber > tripAdvisorVisitorsMax {
			tripAdvisorVisitorsMax = location.TripAdvisorReviewsNumber
		}
		f.StartLocationDistances[idx] = locations.EuclidianDistance(start, location)
		f.EndLocationDistance[idx] = locations.EuclidianDistance(end, location)
	}
	f.StartEndDistance = locations.EuclidianDistance(start, end)

	f.InstagramVisitorsMax = instagramVisitorsMax
	f.TripAdvisorVisitorsMax = tripAdvisorVisitorsMax
	f.StartLocation = start
	f.EndLocation = end
	return f
}

func (f OPScore) LocationScore(route map[int]generic.Location, orderOfLocations []int,
	location generic.Location, id int) float64 {
	distanceCoefficient := 0.0

	if len(orderOfLocations) != 0 && id != f.EndID && id != f.StartID {
		idLastLocation := orderOfLocations[len(orderOfLocations)-1]
		loc1 := route[idLastLocation].(locations.BaseLocation)
		loc2 := location.(locations.BaseLocation)
		distanceCoefficient = f.EndLocationDistance[idLastLocation] / (locations.EuclidianDistance(loc1, loc2) +
			f.EndLocationDistance[id])
		if locations.WalkingTime(loc1, loc2) > int(float64(f.EndLocationDistance[idLastLocation])/66.7) {
			distanceCoefficient = 0.0
		}
	} else {
		distanceCoefficient = f.StartEndDistance / (f.StartLocationDistances[id] + f.EndLocationDistance[id])
	}
	loc := location.(locations.BaseLocation)
	score := loc.OfficialGuide + loc.FoursquareRating/10.0 +
		(loc.TripAdvisorRating/5.0)*(loc.TripAdvisorReviewsNumber/f.TripAdvisorVisitorsMax) +
		loc.InstagramVisitorsNumber/f.InstagramVisitorsMax +
		distanceCoefficient

	return score
}

func (f OPScore) RouteScore(route map[int]generic.Location, orderOfLocations []int) float64 {
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
		routeScore += f.LocationScore(route, orderOfLocations[:i+1], route[key], key)
	}
	return routeScore
}

func (f OPScore) ComputeRouteTimeFromSample(locationsID []int, allLocations []generic.Location) int {
	lastLocation := allLocations[len(locationsID)-1].(locations.BaseLocation)
	time := lastLocation.Duration

	for i := 0; i < len(locationsID)-1; i++ {
		loc1 := allLocations[locationsID[i]].(locations.BaseLocation)
		loc2 := allLocations[locationsID[i+1]].(locations.BaseLocation)
		time += loc1.Duration + locations.WalkingTime(loc1, loc2)
	}

	return time
}
