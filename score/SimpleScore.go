package score

import (
	//	"math"
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
)

type SimpleScore struct {
	TimeLimit            int
	StartID              int
	EndID                int
	InstagramVisitorsMax float64

	StartLocationDistances map[int]float64
	EndLocationDistance    map[int]float64
	StartEndDistance       float64
	StartLocation          generic.Point
	EndLocation            generic.Point
}

func (f SimpleScore) Init(locs []generic.Point) generic.Score {
	instagramVisitorsMax := 0.0

	start := locs[f.StartID].(points.BaseLocation)
	end := locs[f.EndID].(points.BaseLocation)

	f.StartLocationDistances = make(map[int]float64)
	f.EndLocationDistance = make(map[int]float64)

	for idx, loc := range locs {
		location := loc.(points.BaseLocation)
		if location.InstagramVisitorsNumber > instagramVisitorsMax {
			instagramVisitorsMax = location.InstagramVisitorsNumber
		}

		f.StartLocationDistances[idx] = points.EuclidianDistance(start, location)
		f.EndLocationDistance[idx] = points.EuclidianDistance(end, location)
	}
	f.StartEndDistance = points.EuclidianDistance(start, end)

	f.InstagramVisitorsMax = instagramVisitorsMax

	f.StartLocation = start
	f.EndLocation = end
	return f
}

func (f SimpleScore) SinglePointScore(route map[int]generic.Point, orderOfLocations []int,
	location generic.Point, id int) float64 {

	loc := location.(points.BaseLocation)

	return loc.InstagramVisitorsNumber
}

func (f SimpleScore) RouteScore(route map[int]generic.Point, orderOfLocations []int) float64 {
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

func (f SimpleScore) UpdateScore(route map[int]generic.Point, orderOfPoints []int, locs map[int]generic.Point) generic.Score {
	return f
}
