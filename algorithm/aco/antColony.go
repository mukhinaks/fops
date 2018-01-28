package aco

import (
	"math"

	"github.com/mukhinaks/fops/generic"
)

type ACO struct {
	pheromones            map[int]map[int]float64
	fadeness              float64
	attractivenessControl float64
	pheromoneControl      float64
	antsNumber            float64
	iterations            int
	currentIterations     int
	solver                *generic.Solver
	numberOfChannels      int //
}

func (colony ACO) Init(solver *generic.Solver) generic.PathAlgorithm {
	colony.solver = solver
	colony.fadeness = solver.Configuration["Fadeness"].(float64)
	colony.attractivenessControl = solver.Configuration["AttractivenessControl"].(float64)
	colony.pheromoneControl = solver.Configuration["PheromoneControl"].(float64)

	colony.iterations = int(solver.Configuration["Iterations"].(float64))
	colony.pheromones = make(map[int](map[int]float64))
	colony.antsNumber = solver.Configuration["AntsNumber"].(float64)
	colony.numberOfChannels = int(solver.Configuration["NumberOfChannels"].(float64))
	return colony
}

func (colony ACO) CreateRoute() (map[int]generic.Location, []int, float64) {
	colony.pheromones = make(map[int](map[int]float64))

	var bestRoute map[int]generic.Location
	var bestOrder []int
	var bestScore float64
	candidatesLocations := colony.solver.Locations.GetCurrentLocations()
	antsNumber := int(math.Min(colony.antsNumber*float64(len(colony.solver.Locations.GetAllLocations())),
		float64(len(candidatesLocations))*colony.antsNumber)) + 1
	antsNumberPerChannel := int(antsNumber/colony.numberOfChannels) + 1

	if antsNumberPerChannel == 0 {
		antsNumberPerChannel = 1
	}

	for i := 0; i < colony.iterations; i++ {

		allAntsPheromones := make([]map[int](map[int]float64), colony.numberOfChannels)
		allBestRoutes := make([]map[int]generic.Location, colony.numberOfChannels)
		allBestOrders := make([][]int, colony.numberOfChannels)
		allBestScores := make([]float64, colony.numberOfChannels)

		sem := make(chan bool, colony.numberOfChannels)
		for k := 0; k < colony.numberOfChannels; k++ {
			go func(k int) {
				ants := make([]Ant, antsNumberPerChannel)

				deltaPheromones := make(map[int](map[int]float64))

				var localBestRoute map[int]generic.Location
				var localBestOrder []int
				var localBestScore float64

				for _, ant := range ants {
					ant.Init(&candidatesLocations, colony)

					finalRoute, finalOrder, finalScore, antPheromones := ant.GetRoute()

					if finalScore >= localBestScore {
						localBestRoute = finalRoute
						localBestOrder = finalOrder
						localBestScore = finalScore
					}

					for startKey, elem := range antPheromones {
						for endKey, delta := range elem {
							_, ok := deltaPheromones[startKey]
							if ok {
								value, exists := deltaPheromones[startKey][endKey]
								if exists {
									deltaPheromones[startKey][endKey] = value + delta/(float64(antsNumberPerChannel)*float64(colony.numberOfChannels))
								} else {
									deltaPheromones[startKey][endKey] = delta / (float64(antsNumberPerChannel) * float64(colony.numberOfChannels))
								}
							} else {
								deltaPheromones[startKey] = make(map[int]float64)
								deltaPheromones[startKey][endKey] = delta / (float64(antsNumberPerChannel) * float64(colony.numberOfChannels))
							}
						}

					}
				}

				allAntsPheromones[k] = deltaPheromones
				allBestRoutes[k] = localBestRoute
				allBestOrders[k] = localBestOrder
				allBestScores[k] = localBestScore
				sem <- true
			}(k)
		}
		for i := 0; i < colony.numberOfChannels; i++ {
			<-sem
		}

		for i, score := range allBestScores {
			if score >= bestScore {
				bestRoute = allBestRoutes[i]
				bestOrder = allBestOrders[i]
				bestScore = score
			}
		}
		iterationDeltaPheromones := make(map[int](map[int]float64))
		for _, deltaPheromones := range allAntsPheromones {
			for startKey, elem := range deltaPheromones {
				for endKey, delta := range elem {
					_, ok := iterationDeltaPheromones[startKey]
					if ok {
						value, exists := iterationDeltaPheromones[startKey][endKey]
						if exists {
							iterationDeltaPheromones[startKey][endKey] = value + delta
						} else {
							iterationDeltaPheromones[startKey][endKey] = delta
						}
					} else {
						iterationDeltaPheromones[startKey] = make(map[int]float64)
						iterationDeltaPheromones[startKey][endKey] = delta
					}
				}

			}
		}

		colony.UpdatePheromones(iterationDeltaPheromones)
		colony.currentIterations++
	}

	return bestRoute, bestOrder, bestScore
}

func (colony *ACO) UpdatePheromones(deltaPheromones map[int](map[int]float64)) {
	for startKey, elem := range colony.pheromones {
		for endKey, pheromone := range elem {
			colony.pheromones[startKey][endKey] = pheromone * colony.fadeness
		}
	}

	for startKey, elem := range deltaPheromones {
		_, ok := colony.pheromones[startKey]
		if !ok {
			colony.pheromones[startKey] = make(map[int]float64)
		}
		for endKey, delta := range elem {
			_, ok := colony.pheromones[startKey][endKey]
			if ok {
				colony.pheromones[startKey][endKey] = colony.pheromones[startKey][endKey] + delta
			} else {
				colony.pheromones[startKey][endKey] =
					math.Pow(colony.fadeness, float64(colony.currentIterations)) + delta
			}
		}
	}
}

func (colony ACO) GetRawLocations() []generic.Location {
	return colony.solver.Locations.GetAllLocations()
}
