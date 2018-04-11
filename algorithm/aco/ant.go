package aco

import (
	"math"
	"math/rand"
	"time"

	"github.com/mukhinaks/fops/generic"
)

type Ant struct {
	route map[int]generic.Point
	score float64
	keys  []int

	locations *map[int]generic.Point
	colony    ACO
}

func (ant *Ant) Init(locations *map[int]generic.Point, colony ACO) {
	rand.Seed(time.Now().UnixNano())
	ant.locations = locations
	ant.colony = colony
}

func (ant *Ant) NextLocation() (bool, int) {
	probabilities := make(map[int]float64)
	probabilitiesSum := 0.0

	actualLocations := ant.colony.solver.Constraints.ReducePoints(ant.route, ant.keys, *ant.locations)
	currentScore := ant.colony.solver.Score.UpdateScore(ant.route, ant.keys, actualLocations)

	if len(actualLocations) == 0 {
		return false, -1
	}

	for key, location := range actualLocations {
		_, inRoute := ant.route[key]
		if inRoute {
			continue
		}

		pheromone := math.Pow(ant.colony.fadeness, float64(ant.colony.currentIterations))
		if len(ant.keys) > 0 {
			if data, isFinish := ant.colony.pheromones[ant.keys[len(ant.keys)-1]][key]; isFinish {
				pheromone = data
			}
		}
		locationScore := currentScore.SinglePointScore(ant.route, ant.keys, location, key)

		probability := math.Pow(pheromone, ant.colony.pheromoneControl) * math.Pow(locationScore, ant.colony.attractivenessControl)
		//if len(probabilities) < 30 {
		probabilities[key] = probability
		probabilitiesSum += probability
		/*
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
		*/
	}

	if len(probabilities) == 0 {
		return false, -1
	}

	randomNumber := rand.Float64() * probabilitiesSum
	/*
		candidate := make(map[int]generic.Point)
		for key, value := range ant.route {
			candidate[key] = value
		}
	*/
	cumProbabiltySum := 0.0
	index := -1
	for key, prob := range probabilities {
		cumProbabiltySum += prob
		if cumProbabiltySum >= randomNumber {
			index = key
			break
		}
	}
	if index == -1 {
		return false, -1
	}
	ant.route[index] = actualLocations[index]

	return ant.colony.solver.Constraints.Boundary(ant.route, append(ant.keys, index)), index
}

func (ant *Ant) GetRoute() {

	ant.route = make(map[int]generic.Point)
	ant.keys = make([]int, 0)

	for {
		flag, index := ant.NextLocation()
		if flag {
			ant.keys = append(ant.keys, index)
		} else {
			ant.score = ant.colony.solver.Score.RouteScore(ant.route, ant.keys)
			return
		}
	}
}
