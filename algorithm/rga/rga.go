package rga

import (
	"github.com/mukhinaks/fops/generic"
)

type RGA struct {
	route  map[int]generic.Point
	keys   []int
	score  float64
	solver *generic.Solver
}

func (rga RGA) Init(solver *generic.Solver) generic.PathAlgorithm {
	rga.solver = solver

	return rga
}

func (rga RGA) SetInitialRoute(route map[int]generic.Point, keys []int) generic.PathAlgorithm {
	rga.route = route
	rga.keys = keys

	return rga
}

func (rga RGA) CreateRoute() (map[int]generic.Point, []int, float64) {

	for {
		flag, candidate, keys := rga.SelectBestCandidateFromAllIntervals()
		if flag {
			rga.route = candidate
			rga.keys = keys
		} else {
			rga.score = rga.solver.Score.RouteScore(rga.route, rga.keys)
			break
		}
	}

	return rga.route, rga.keys, rga.score
}

func (rga RGA) SelectBestCandidateFromAllIntervals() (bool, map[int]generic.Point, []int) {
	var bestRoute map[int]generic.Point
	var bestOrder []int
	var bestScore float64
	flag := false

	candidatesLocations := rga.solver.Points.GetCurrentPoints()
	for i := 0; i < len(rga.keys)-1; i++ {
		idx, candidate, candidateKeys := rga.InsertLocationInInterval(rga.keys[i], rga.keys[i+1], candidatesLocations)
		if idx != -1 {
			score := rga.solver.Score.RouteScore(candidate, candidateKeys)

			if score > bestScore {
				flag = true
				bestRoute = candidate
				bestOrder = candidateKeys
			}
		}
	}

	return flag, bestRoute, bestOrder
}

func (rga RGA) InsertLocationInInterval(startID int, endID int, locations map[int]generic.Point) (int, map[int]generic.Point, []int) {

	constraint := rga.solver.Constraints.UpdateConstraint(nil, []int{startID, endID}, rga.solver.Points.GetAllPoints())
	actualLocations := constraint.ReducePoints(nil, nil, locations) // noname.solver.Points.GetPointsInArea(startID, endID) //
	currentScore := rga.solver.Score.UpdateScore(rga.route, rga.keys, actualLocations)

	if len(actualLocations) == 0 {
		return -1, nil, nil
	}

	maxScore := 0.0
	maxScoreID := -1
	var bestRoute map[int]generic.Point
	var bestOrder []int

	for key, location := range actualLocations {
		_, inRoute := rga.route[key]
		if inRoute {
			continue
		}
		intervalRoute := make(map[int]generic.Point)

		intervalRoute[key] = location

		intervalKeys := []int{key}
		locationScore := currentScore.SinglePointScore(intervalRoute, intervalKeys, location, key)
		if locationScore > maxScore {
			candidate := make(map[int]generic.Point)
			for index, value := range rga.route {
				candidate[index] = value
			}
			candidate[key] = location

			candidateKeys := make([]int, 0)
			for _, value := range rga.keys {
				candidateKeys = append(candidateKeys, value)
				if value == startID {
					candidateKeys = append(candidateKeys, key)
				}
			}

			if rga.solver.Constraints.Boundary(candidate, candidateKeys) {
				maxScore = locationScore
				maxScoreID = key
				bestOrder = candidateKeys
				bestRoute = candidate

			}

		}
	}

	return maxScoreID, bestRoute, bestOrder
}
