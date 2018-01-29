package aco

import (
	"math"
	"math/rand"
	"time"

	"github.com/mukhinaks/fops/generic"
)

type Ant struct {
	route           map[int]generic.Point
	score           float64
	keys            []int
	deltaPheromones map[int](map[int]float64)
	locations       *map[int]generic.Point
	colony          ACO
}

func (ant *Ant) Init(locations *map[int]generic.Point, colony ACO) {
	rand.Seed(time.Now().UnixNano())
	ant.locations = locations
	ant.colony = colony
}

func (ant *Ant) NextLocation() (bool, map[int]generic.Point, int) {
	probabilities := make(map[int]float64)
	probabilitiesSum := 0.0

	if len(*ant.locations) == 0 {
		return false, nil, -1
	}
	for key, location := range *ant.locations {
		_, inRoute := ant.route[key]
		if inRoute {
			continue
		}

		pheromone := math.Pow(ant.colony.fadeness, float64(ant.colony.currentIterations))
		if len(ant.keys) > 0 {
			if _, isStart := ant.colony.pheromones[ant.keys[len(ant.keys)-1]]; isStart {
				if data, isFinish := ant.colony.pheromones[ant.keys[len(ant.keys)-1]][key]; isFinish {
					pheromone = data
				}
			}
		}
		locationScore := ant.colony.solver.Score.SinglePointScore(ant.route, ant.keys, location, key)

		probability := math.Pow(pheromone, ant.colony.pheromoneControl) * math.Pow(locationScore, ant.colony.attractivenessControl)
		if len(probabilities) < 30 {
			probabilities[key] = probability
			probabilitiesSum += probability
		} else {

			for idx, prob := range probabilities {
				if prob < probability {
					delete(probabilities, idx)
					probabilitiesSum -= prob
					probabilities[key] = probability
					probabilitiesSum += probability
					break
				}
			}
		}
	}
	for key, prob := range probabilities {
		probabilities[key] = prob / probabilitiesSum
	}

	randomNumber := rand.Float64()
	candidate := make(map[int]generic.Point)
	for key, value := range ant.route {
		candidate[key] = value
	}
	cumProbabiltySum := 0.0
	index := 0
	for key, prob := range probabilities {
		cumProbabiltySum += prob
		index = key
		if cumProbabiltySum >= randomNumber {
			break
		}
	}
	candidate[index] = (*ant.locations)[index]
	return ant.colony.solver.Constraints.Boundary(candidate, append(ant.keys, index)), candidate, index
}

func (ant *Ant) GetRoute() (map[int]generic.Point, []int, float64, map[int](map[int]float64)) {

	ant.route = make(map[int]generic.Point)
	ant.keys = make([]int, 0)

	for {
		flag, candidate, index := ant.NextLocation()
		if flag {
			ant.route = candidate
			ant.keys = append(ant.keys, index)
		} else {
			ant.score = ant.colony.solver.Score.RouteScore(ant.route, ant.keys)
			break
		}
	}

	ant.deltaPheromones = make(map[int](map[int]float64))

	for key := 0; key < len(ant.keys)-1; key++ {
		startKey := ant.keys[key]
		endKey := ant.keys[key+1]

		_, ok := ant.deltaPheromones[startKey]
		if !ok {
			ant.deltaPheromones[startKey] = make(map[int]float64)
		}
		ant.deltaPheromones[startKey][endKey] = ant.score / float64(len(ant.route))
	}

	return ant.route, ant.keys, ant.score, ant.deltaPheromones
}
