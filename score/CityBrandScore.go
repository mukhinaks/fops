package score

import (
	//	"math"
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type CityBrandScore struct {
	TimeLimit                int
	StartID                  int
	InstagramVisitorsMax     float64
	TripAdvisorVisitorsMax   float64
	FacebookRatingMax        float64
	FoursquareRatingVotesMax float64
	StartLocation            generic.Point
	//MaximumDistanceToStart float64
}

func (f CityBrandScore) Init(locs []generic.Point) generic.Score {
	instagramVisitorsMax := 0.0
	tripAdvisorVisitorsMax := 0.0
	facebookRatingMax := 0.0
	foursquareRatingVotesMax := 0.0

	start := locs[f.StartID].(points.CityBrandLocation)
	//maxDistance := 0.0

	for _, loc := range locs {
		location := loc.(points.CityBrandLocation)
		drop := false
		for _, category := range location.Categories {
			if category == "Restaurant" {
				drop = true
				break
			}
		}
		if drop {
			continue
		}

		if location.InstagramVisitorsNumber >= instagramVisitorsMax {
			instagramVisitorsMax = location.InstagramVisitorsNumber
		}
		if location.FacebookRating >= facebookRatingMax {
			facebookRatingMax = location.FacebookRating
		}

		if location.TripAdvisorReviewsNumber >= tripAdvisorVisitorsMax {
			tripAdvisorVisitorsMax = location.TripAdvisorReviewsNumber
		}

		if location.FoursquareRatingVotes >= foursquareRatingVotesMax {
			foursquareRatingVotesMax = location.FoursquareRatingVotes
		}

		/*
			distance := points.EuclidianDistance(location, start)
			if distance >= maxDistance {
				maxDistance = distance
			}
		*/
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

	if facebookRatingMax == 0 {
		f.FacebookRatingMax = 1
	} else {

		f.FacebookRatingMax = facebookRatingMax
	}

	f.StartLocation = start
	//f.MaximumDistanceToStart = maxDistance
	return f
}

func (f CityBrandScore) SinglePointScore(route map[int]generic.Point, orderOfLocations []int,
	location generic.Point, id int) float64 {
	distanceCoefficient := 0.0

	if len(orderOfLocations) != 0 && id != f.StartID {
		idLastLocation := orderOfLocations[len(orderOfLocations)-1]
		loc1 := route[idLastLocation].(points.CityBrandLocation)
		loc2 := location.(points.CityBrandLocation)
		distance := points.EuclidianDistance(loc1, loc2)
		distanceCoefficient = 1.0

		for i := 0; i < len(orderOfLocations)-1; i++ {
			if points.EuclidianDistance(loc2, route[orderOfLocations[i]].(points.CityBrandLocation)) <= distance {
				distanceCoefficient = 0
			}
		}

	} else {
		distanceCoefficient = 1
	}
	loc := location.(points.CityBrandLocation)
	score := float64(len(loc.CityBrand))/10.0 + ((loc.FoursquareRating/10.0)*loc.FoursquareRatingVotes/f.FoursquareRatingVotesMax+
		(loc.TripAdvisorRating/5.0)*loc.TripAdvisorReviewsNumber/f.TripAdvisorVisitorsMax+
		loc.InstagramVisitorsNumber/f.InstagramVisitorsMax+
		loc.FacebookRating/f.FacebookRatingMax/2+
		0)/1.0 +
		distanceCoefficient/4

	return score
}

func (f CityBrandScore) RouteScore(route map[int]generic.Point, orderOfLocations []int) float64 {
	_, ok := route[f.StartID]
	if !ok {
		route[f.StartID] = f.StartLocation
	}

	routeScore := 0.0
	for i, key := range orderOfLocations {
		routeScore += f.SinglePointScore(route, orderOfLocations[:i+1], route[key], key)
	}
	return routeScore
}

func (f CityBrandScore) ComputeRouteTimeFromSample(locationsID []int, allLocations []generic.Point) int {
	lastLocation := allLocations[len(locationsID)-1].(points.CityBrandLocation)
	time := lastLocation.Duration

	for i := 0; i < len(locationsID)-1; i++ {
		loc1 := allLocations[locationsID[i]].(points.CityBrandLocation)
		loc2 := allLocations[locationsID[i+1]].(points.CityBrandLocation)
		time += loc1.Duration + points.WalkingTime(loc1, loc2)
	}

	return time
}

func (f CityBrandScore) SinglePointScoreWithoutPositionDependance(location generic.Point, id int) float64 {

	loc := location.(points.CityBrandLocation)
	for _, category := range loc.Categories {
		if category == "Restaurant" {
			return 0
		}
	}
	score := float64(len(loc.CityBrand))/10.0 + ((loc.FoursquareRating/10.0)*loc.FoursquareRatingVotes/f.FoursquareRatingVotesMax+
		(loc.TripAdvisorRating/5.0)*loc.TripAdvisorReviewsNumber/f.TripAdvisorVisitorsMax+
		loc.InstagramVisitorsNumber/f.InstagramVisitorsMax+
		loc.FacebookRating/f.FacebookRatingMax+
		loc.WikipediaPage/100)/5.0

	return score
}

func (f CityBrandScore) UpdateScore(route map[int]generic.Point, orderOfPoints []int, locs map[int]generic.Point) generic.Score {
	return f
}
