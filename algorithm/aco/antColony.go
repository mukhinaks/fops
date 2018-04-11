package aco

import (
	"math"
	"runtime"
	"sync"

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

type Deltas struct {
	start int
	end   int
	delta float64
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
	runtime.GOMAXPROCS(12)
	return colony
}

func (colony ACO) CreateRoute() (map[int]generic.Point, []int, float64) {
	colony.pheromones = make(map[int](map[int]float64))

	var bestRoute map[int]generic.Point
	var bestOrder []int
	var bestScore float64
	candidatesLocations := colony.solver.Points.GetCurrentPoints()
	antsNumber := int(float64(len(candidatesLocations))*colony.antsNumber) + 1

	for i := 0; i < colony.iterations; i++ {

		iterationPheromones := make([]Deltas, 0)

		sem := make(chan bool, antsNumber)
		mux := sync.Mutex{}
		for k := 0; k < antsNumber; k++ {
			go func(i int) {
				ant := Ant{}
				ant.Init(&candidatesLocations, colony)
				ant.GetRoute()

				mux.Lock()
				if ant.score >= bestScore {
					bestRoute = ant.route
					bestOrder = ant.keys
					bestScore = ant.score
				}

				for key := 0; key < len(ant.keys)-1; key++ {
					startKey := ant.keys[key]
					endKey := ant.keys[key+1]

					iterationPheromones = append(iterationPheromones, Deltas{startKey, endKey, ant.score / float64(antsNumber)})
				}
				mux.Unlock()
				sem <- true
			}(k)
		}

		//t := time.Now()
		for i := 0; i < antsNumber; i++ {
			<-sem
		}
		//		fmt.Println(time.Since(t))

		colony.UpdatePheromones(iterationPheromones)
		colony.currentIterations++
	}

	return bestRoute, bestOrder, bestScore
}

func (colony *ACO) UpdatePheromones(allAntsPheromones []Deltas) {
	for startKey, elem := range colony.pheromones {
		for endKey, pheromone := range elem {
			colony.pheromones[startKey][endKey] = pheromone * colony.fadeness
		}
	}

	for _, deltas := range allAntsPheromones {
		delta := deltas.delta
		_, ok := colony.pheromones[deltas.start]
		if ok {
			value, ok := colony.pheromones[deltas.start][deltas.end]
			if ok {
				colony.pheromones[deltas.start][deltas.end] = value + delta
			} else {
				colony.pheromones[deltas.start][deltas.end] =
					math.Pow(colony.fadeness, float64(colony.currentIterations)) + delta
			}
		} else {
			colony.pheromones[deltas.start] = make(map[int]float64)
			colony.pheromones[deltas.start][deltas.end] =
				math.Pow(colony.fadeness, float64(colony.currentIterations)) + delta
		}
	}
	/*
		for _, deltaPheromones := range allAntsPheromones {
			for startKey, elem := range deltaPheromones {
				for endKey, delta := range elem {
					_, ok := colony.pheromones[startKey]
					if ok {
						value, ok := colony.pheromones[startKey][endKey]
						if ok {
							colony.pheromones[startKey][endKey] = value + delta
						} else {
							colony.pheromones[startKey][endKey] =
								math.Pow(colony.fadeness, float64(colony.currentIterations)) + delta
						}
					} else {
						colony.pheromones[startKey] = make(map[int]float64)
						colony.pheromones[startKey][endKey] =
							math.Pow(colony.fadeness, float64(colony.currentIterations)) + delta
					}

				}

			}
		}
	*/

}
