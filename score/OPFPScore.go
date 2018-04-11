package score

import (

	//	"math"
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type OPFPScore struct {
	TimeLimit                int
	StartID                  int
	EndID                    int
	InstagramVisitorsMax     float64
	TripAdvisorVisitorsMax   float64
	FoursquareRatingVotesMax float64
	StartLocationDistances   map[int]float64
	EndLocationDistance      map[int]float64
	StartEndDistance         float64
	StartLocation            generic.Point
	EndLocation              generic.Point
}

func (f OPFPScore) Init(locs []generic.Point) generic.Score {
	instagramVisitorsMax := 0.0
	tripAdvisorVisitorsMax := 0.0
	foursquareRatingVotesMax := 0.0
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
		if location.FoursquareRatingVotes > foursquareRatingVotesMax {
			foursquareRatingVotesMax = location.FoursquareRatingVotes
		}
		f.StartLocationDistances[idx] = points.EuclidianDistance(start, location)
		f.EndLocationDistance[idx] = points.EuclidianDistance(end, location)
	}
	f.StartEndDistance = points.EuclidianDistance(start, end)

	if instagramVisitorsMax == 0 {
		f.InstagramVisitorsMax = 1
	} else {
		f.InstagramVisitorsMax = instagramVisitorsMax
	}

	if tripAdvisorVisitorsMax == 0 {
		f.TripAdvisorVisitorsMax = 1
	} else {
		f.TripAdvisorVisitorsMax = tripAdvisorVisitorsMax
	}

	if foursquareRatingVotesMax == 0 {
		f.FoursquareRatingVotesMax = 1
	} else {
		f.FoursquareRatingVotesMax = foursquareRatingVotesMax
	}

	f.StartLocation = start
	f.EndLocation = end
	return f
}

func (f OPFPScore) SinglePointScore(route map[int]generic.Point, orderOfLocations []int,
	location generic.Point, id int) float64 {
	distanceCoefficient := 0.0

	if len(orderOfLocations) != 0 && id != f.EndID && id != f.StartID {
		positionInRoute := -1
		for i, idx := range orderOfLocations {
			if id == idx {
				positionInRoute = i
			}
		}

		previousLocation := f.StartLocation
		if positionInRoute > 0 {
			previousLocation = route[orderOfLocations[positionInRoute-1]].(points.BaseLocation)
		}

		nextLocation := f.EndLocation
		if positionInRoute < len(orderOfLocations)-1 && positionInRoute != -1 {
			nextLocation = route[orderOfLocations[positionInRoute+1]].(points.BaseLocation)
		}

		distanceCoefficient = points.EuclidianDistance(previousLocation, nextLocation) / (points.EuclidianDistance(previousLocation, location) +
			points.EuclidianDistance(location, nextLocation))

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

func (f OPFPScore) RouteScore(route map[int]generic.Point, orderOfLocations []int) float64 {

	_, ok := route[f.StartID]
	if !ok {
		route[f.StartID] = f.StartLocation
	}
	_, ok = route[f.EndID]
	if !ok {
		route[f.EndID] = f.EndLocation
	}

	routeScore := 0.0
	for _, key := range orderOfLocations {
		routeScore += f.SinglePointScore(route, orderOfLocations, route[key], key)
	}
	return routeScore
}

func (f OPFPScore) UpdateScore(route map[int]generic.Point, orderOfPoints []int, locs map[int]generic.Point) generic.Score {

	instagramVisitorsMax := 0.0
	tripAdvisorVisitorsMax := 0.0
	foursquareRatingVotesMax := 0.0

	for _, loc := range locs {
		location := loc.(points.BaseLocation)
		if location.InstagramVisitorsNumber > instagramVisitorsMax {
			instagramVisitorsMax = location.InstagramVisitorsNumber
		}
		if location.TripAdvisorReviewsNumber > tripAdvisorVisitorsMax {
			tripAdvisorVisitorsMax = location.TripAdvisorReviewsNumber
		}
		if location.FoursquareRatingVotes > foursquareRatingVotesMax {
			foursquareRatingVotesMax = location.FoursquareRatingVotes
		}
	}

	if instagramVisitorsMax == 0 {
		f.InstagramVisitorsMax = 1
	} else {
		f.InstagramVisitorsMax = instagramVisitorsMax
	}

	if tripAdvisorVisitorsMax == 0 {
		f.TripAdvisorVisitorsMax = 1
	} else {
		f.TripAdvisorVisitorsMax = tripAdvisorVisitorsMax
	}

	if foursquareRatingVotesMax == 0 {
		f.FoursquareRatingVotesMax = 1
	} else {
		f.FoursquareRatingVotesMax = foursquareRatingVotesMax
	}

	return f
}
