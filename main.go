package main

import (
	"fmt"

	"github.com/mukhinaks/fops/constraints"
	"github.com/mukhinaks/fops/points"

	"github.com/mukhinaks/fops/algorithm/aco"
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/score"
)

func main() {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.OPScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}

	ExperimentClassicalOP(solver, "", 21, 19, 200, "classic-op")
	ExperimentClassicalOPWithReference(solver, "", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "classic-op-with-ref")
	ExperimentOPCV(solver, "", []int{0, 31, 21, 19}, 240, "opcv")
	ExperimentOPCVMultipleDays(solver, "", []int{152, 3, 106, 105, 51, 63, 9, 127, 157, 158, 11, 13, 5191}, 600, 2, "opcv-multiple-days")

}

func ExperimentClassicalOP(solver generic.Solver, configPath string, startID int, endID int, timeLimit int, fileName string) {
	fmt.Println("-----------------")
	fmt.Println("Classical OP")

	locs := points.BaseLocations{}
	sc := score.OPScore{}
	pathAlgorithm := &aco.ACO{}

	c := &constraints.OPConstraints{}

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	finalRoute := make(map[int]generic.Point)
	finalOrder := []int{startID}

	sc.StartID = startID
	sc.EndID = endID
	solver.Score = sc

	c.StartID = startID
	c.EndID = endID
	c.TimeLimit = timeLimit
	solver.Constraints = c

	result, order, _ := solver.NextInterval()

	for _, k := range order {
		finalRoute[k] = result[k]
		finalOrder = append(finalOrder, k)
	}

	finalOrder = append(finalOrder, sc.EndID)
	finalRoute[sc.EndID] = result[sc.EndID]
	finalRoute[sc.StartID] = result[sc.StartID]

	fmt.Println("final score:", solver.Score.RouteScore(finalRoute, finalOrder))
	fmt.Println("final time:", c.FinalRouteTime(finalRoute, finalOrder))

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName+".json")
}

func ExperimentClassicalOPWithReference(solver generic.Solver, configPath string, referencePath []int, fileName string) {
	fmt.Println("-----------------")
	fmt.Println("Classical OP with reference path")

	locs := points.BaseLocations{}
	sc := score.OPScore{}
	pathAlgorithm := &aco.ACO{}

	c := &constraints.OPConstraints{}

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	finalRoute := make(map[int]generic.Point)
	finalOrder := []int{referencePath[0]}

	sc.StartID = referencePath[0]
	sc.EndID = referencePath[len(referencePath)-1]
	solver.Score = sc

	c.StartID = sc.StartID
	c.EndID = sc.EndID
	c.TimeLimit = sc.ComputeRouteTimeFromSample(referencePath, solver.Points.GetAllPoints())
	solver.Constraints = c

	result, order, _ := solver.NextInterval()

	for _, k := range order {
		finalRoute[k] = result[k]
		finalOrder = append(finalOrder, k)
	}

	finalOrder = append(finalOrder, sc.EndID)
	finalRoute[sc.EndID] = result[sc.EndID]
	finalRoute[sc.StartID] = result[sc.StartID]

	fmt.Println("final score:", solver.Score.RouteScore(finalRoute, finalOrder))
	fmt.Println("final time:", c.FinalRouteTime(finalRoute, finalOrder))

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName+".json")
}

func ExperimentOPCV(solver generic.Solver, configPath string, compulsoryLocations []int, routeTimeLimit int, fileName string) {
	fmt.Println("-----------------")
	fmt.Println("OPCV")

	locs := points.BaseLocations{}
	sc := score.OPScore{}
	pathAlgorithm := &aco.ACO{}

	finalRoute := make(map[int]generic.Point)
	finalOrder := []int{compulsoryLocations[0]}
	c := constraints.EnrichmentConstraints{}
	c.ForbiddenLocations = compulsoryLocations
	c.CompulsoryLocations = compulsoryLocations
	c.RouteTimeLimit = routeTimeLimit
	pathAlgorithm = &aco.ACO{}

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	for i := 0; i < len(compulsoryLocations)-1; i++ {
		sc.StartID = compulsoryLocations[i]
		sc.EndID = compulsoryLocations[i+1]
		solver.Score = sc

		c.NumberOfInterval = i

		solver.Constraints = c

		result, order, _ := solver.NextInterval()

		for _, k := range order {
			finalRoute[k] = result[k]
			finalOrder = append(finalOrder, k)
			c.ForbiddenLocations = append(c.ForbiddenLocations, k)
		}

		finalOrder = append(finalOrder, sc.EndID)
		finalRoute[sc.EndID] = result[sc.EndID]
		finalRoute[sc.StartID] = result[sc.StartID]
	}

	fmt.Println("final score:", solver.Score.RouteScore(finalRoute, finalOrder))
	fmt.Println("final time:", c.FinalRouteTime(finalRoute, finalOrder))
	fmt.Println("original time:", c.FinalRouteTime(finalRoute, compulsoryLocations))
	fmt.Println("original score:", solver.Score.RouteScore(finalRoute, compulsoryLocations))
	locs.WriteLocationsToJSON(finalRoute, compulsoryLocations, fileName+"-origin.json")
	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName+".json")

}

func ExperimentOPCVMultipleDays(solver generic.Solver, configPath string, compulsoryLocations []int, dayTimeLimit int, dayNumber int, fileName string) {
	fmt.Println("-----------------")
	fmt.Println("OPCV for Mutiple Days")

	locs := points.BaseLocations{}
	sc := score.OPScore{}
	pathAlgorithm := &aco.ACO{}

	finalRoute := make(map[int]generic.Point)
	finalOrder := make([]int, 0)

	c := constraints.MultidaysConstraints{}
	c.ForbiddenLocations = compulsoryLocations
	c.CompulsoryLocations = compulsoryLocations
	c.DayTimeLimit = dayTimeLimit
	c.DaysNumber = dayNumber

	pathAlgorithm = &aco.ACO{}

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)
	days, times := c.SplitForDays(c.CompulsoryLocations, solver.Points.GetAllPoints())
	for i := 1; i <= c.DaysNumber; i++ {
		if len(days[i]) == 0 {
			continue
		}
		c.CurrentDay = i
		c.TimeLimit = times[i]
		c.CompulsoryLocations = days[i]
		finalOrder = append(finalOrder, days[i][0])

		for j := 0; j < len(days[i])-1; j++ {
			sc.StartID = days[i][j]
			sc.EndID = days[i][j+1]
			solver.Score = sc

			c.NumberOfInterval = j
			solver.Constraints = c

			result, order, _ := solver.NextInterval()

			for _, k := range order {
				finalOrder = append(finalOrder, k)
				c.ForbiddenLocations = append(c.ForbiddenLocations, k)
			}

			finalOrder = append(finalOrder, sc.EndID)
			end := result[sc.EndID].(points.BaseLocation)
			end.ID = c.CurrentDay
			finalRoute[sc.EndID] = end
			start := result[sc.StartID].(points.BaseLocation)
			start.ID = c.CurrentDay
			finalRoute[sc.StartID] = start

			for _, k := range order {
				l := result[k].(points.BaseLocation)
				l.ID = c.CurrentDay
				finalRoute[k] = l
			}
		}
	}

	fmt.Println("final score:", solver.Score.RouteScore(finalRoute, finalOrder))
	fmt.Println("final time:", c.FinalRouteTime(finalRoute, finalOrder))
	fmt.Println("original time:", c.FinalRouteTime(finalRoute, compulsoryLocations))
	fmt.Println("original score:", solver.Score.RouteScore(finalRoute, compulsoryLocations))
	locs.WriteLocationsToJSON(finalRoute, compulsoryLocations, fileName+"-origin.json")
	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName+".json")

}
