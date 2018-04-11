package main

import (
	"fmt"

	"github.com/mukhinaks/fops/algorithm/aco"
	"github.com/mukhinaks/fops/algorithm/rga"
	"github.com/mukhinaks/fops/constraints"
	"github.com/mukhinaks/fops/generic"
	"github.com/mukhinaks/fops/points"
	"github.com/mukhinaks/fops/score"
)

type AvailableAlgortihms struct {
	ACO string
	RGA string
}

func (f AvailableAlgortihms) Init() AvailableAlgortihms {
	f.ACO = "ACO"
	f.RGA = "RGA"
	return f
}

// SolveClassicalOP solves classic Orienteering Problem.
// Result is optimal path with highest total score from start to end node considering giving time budget.
func SolveClassicalOP(configPath string, startID int, endID int, timeLimit int, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
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

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)

	return
}

// SolveClassicalOP solves classic Orienteering Problem.
// Result is optimal path with highest total score from start to end node considering giving time budget.
func SolveClassicalOPByRGA(configPath string, startID int, endID int, timeLimit int, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   rga.RGA{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
	pathAlgorithm := rga.RGA{}

	c := &constraints.OPConstraints{}

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	locations := solver.Points.GetAllPoints()
	intervalRoute := make(map[int]generic.Point)

	intervalRoute[startID] = locations[startID]
	intervalRoute[endID] = locations[endID]

	intervalKeys := []int{startID, endID}

	pathAlgorithm = pathAlgorithm.Init(&solver).(rga.RGA)
	solver.Algorithm = pathAlgorithm.SetInitialRoute(intervalRoute, intervalKeys)

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

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)

	return
}

// SolveClassicalOPWithReferencePath solves classic Orienteering Problem.
// Result is the best route from first to last node of reference path with respect of original path's time budget.
func SolveClassicalOPWithReferencePath(configPath string, referencePath []int, fileName string) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
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
	c.TimeLimit = c.ComputeRouteTimeFromSample(referencePath, solver.Points.GetAllPoints())
	solver.Constraints = c

	result, order, _ := solver.NextInterval()

	for _, k := range order {
		finalRoute[k] = result[k]
		finalOrder = append(finalOrder, k)
	}

	finalOrder = append(finalOrder, sc.EndID)
	finalRoute[sc.EndID] = result[sc.EndID]
	finalRoute[sc.StartID] = result[sc.StartID]

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
}

// SolveTDOP solves the Time Dependent Orienteering Problem.
// Result is optimal path with highest total score from start to end node.
// Transition time is normally distributed where mean is product of distance between two nodes and walking velocity (4 km/h).
func SolveTDOP(configPath string, startID int, endID int, timeLimit int, startTime int, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
	pathAlgorithm := &aco.ACO{}

	c := &constraints.TDOPConstraints{}

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
	c.StartTime = startTime
	solver.Constraints = c

	result, order, _ := solver.NextInterval()

	for _, k := range order {
		finalRoute[k] = result[k]
		finalOrder = append(finalOrder, k)
	}

	finalOrder = append(finalOrder, sc.EndID)
	finalRoute[sc.EndID] = result[sc.EndID]
	finalRoute[sc.StartID] = result[sc.StartID]

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
	return
}

// SolveTDOP solves the Time Dependent Orienteering Problem.
// Result is optimal path with highest total score from start to end node.
// Transition time is normally distributed where mean is product of distance between two nodes and walking velocity (4 km/h).
func SolveTDOPByRGA(configPath string, startID int, endID int, timeLimit int, startTime int, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   rga.RGA{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
	pathAlgorithm := rga.RGA{}

	c := &constraints.TDOPConstraints{}

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	locations := solver.Points.GetAllPoints()
	intervalRoute := make(map[int]generic.Point)

	intervalRoute[startID] = locations[startID]
	intervalRoute[endID] = locations[endID]

	intervalKeys := []int{startID, endID}

	pathAlgorithm = pathAlgorithm.Init(&solver).(rga.RGA)
	solver.Algorithm = pathAlgorithm.SetInitialRoute(intervalRoute, intervalKeys)

	finalRoute := make(map[int]generic.Point)
	finalOrder := []int{startID}

	sc.StartID = startID
	sc.EndID = endID
	solver.Score = sc

	c.StartID = startID
	c.EndID = endID
	c.TimeLimit = timeLimit
	c.StartTime = startTime
	solver.Constraints = c

	result, order, _ := solver.NextInterval()

	for _, k := range order {
		finalRoute[k] = result[k]
		finalOrder = append(finalOrder, k)
	}

	finalOrder = append(finalOrder, sc.EndID)
	finalRoute[sc.EndID] = result[sc.EndID]
	finalRoute[sc.StartID] = result[sc.StartID]

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
	return
}

// SolveOPCV solves Orienteering Problem with Compulsory Vertices.
// Resulting path consists all locations from the set of compulsory locations.
func SolveOPCV(configPath string, compulsoryLocations []int, routeTimeLimit int, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
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

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
	return
}

// SolveOPCV solves Orienteering Problem with Compulsory Vertices.
// Resulting path consists all locations from the set of compulsory locations.
func SolveOPCVByRGA(configPath string, compulsoryLocations []int, routeTimeLimit int, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   rga.RGA{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
	//pathAlgorithm := &aco.ACO{}

	finalRoute := make(map[int]generic.Point)
	finalOrder := []int{compulsoryLocations[0]}
	c := constraints.EnrichmentConstraints{}
	c.ForbiddenLocations = compulsoryLocations
	c.CompulsoryLocations = compulsoryLocations
	c.RouteTimeLimit = routeTimeLimit
	pathAlgorithm := rga.RGA{}

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

		locations := solver.Points.GetAllPoints()
		intervalRoute := make(map[int]generic.Point)

		intervalRoute[sc.StartID] = locations[sc.StartID]
		intervalRoute[sc.EndID] = locations[sc.EndID]

		intervalKeys := []int{sc.StartID, sc.EndID}

		pathAlgorithm = pathAlgorithm.Init(&solver).(rga.RGA)
		solver.Algorithm = pathAlgorithm.SetInitialRoute(intervalRoute, intervalKeys)

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

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
	return
}

// SolveOPCVForMultipleDaysRoute solves Orienteering Problem with Compulsory Vertices.
// The set of compulsory locations is split in predeined number of days.
// WARNING: if all compulsory locations cannot be visited in defined time budget time budget will be expanded!
// WARNING: if compulsory locations can be visited in less number of days the shorter route will be created.
func SolveOPCVForMultipleDaysRoute(configPath string, compulsoryLocations []int, dayTimeLimit int, dayNumber int, fileName string) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
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

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
}

// SolveOPTW solves Orienteering Problem with Time Windows.
// Resulting path contains locations according to their open hours.
func SolveOPTW(configPath string, startID int, endID int, timeLimit int, startTime int, dayOfWeek string, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}
	pathAlgorithm := &aco.ACO{}

	c := &constraints.OPTWConstraints{}
	c.TimeLimit = timeLimit
	c.DayOfWeek = dayOfWeek
	c.StartTime = startTime

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	sc.StartID = startID
	sc.EndID = endID
	solver.Score = sc

	c.StartID = startID
	c.EndID = endID
	solver.Constraints = c

	result, order, _ := solver.NextInterval()
	finalRoute := make(map[int]generic.Point)
	finalOrder := []int{startID}

	for _, k := range order {
		finalRoute[k] = result[k]
		finalOrder = append(finalOrder, k)
	}

	finalOrder = append(finalOrder, sc.EndID)
	finalRoute[sc.EndID] = result[sc.EndID]
	finalRoute[sc.StartID] = result[sc.StartID]

	// Commented part prints time of visit in final route
	/*
		currentTime := c.StartTime
		for i, k := range finalOrder {
			if i > 0 {
				visit := finalRoute[finalOrder[i-1]].(points.BaseLocation).Duration + points.WalkingTime(finalRoute[finalOrder[i-1]].(points.BaseLocation), finalRoute[k].(points.BaseLocation))
				currentTime = c.TimeUpdate(currentTime, visit)
			}
			fmt.Println(k, finalRoute[k].(points.BaseLocation).Title, currentTime)
		}
	*/

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
	return
}

// SolveOPTW solves Orienteering Problem with Time Windows.
// Resulting path contains locations according to their open hours.
func SolveOPTWByRGA(configPath string, startID int, endID int, timeLimit int, startTime int, dayOfWeek string, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   rga.RGA{},
	}

	locs := points.BaseLocations{}
	sc := score.SimpleScore{}

	c := &constraints.OPTWConstraints{}
	c.TimeLimit = timeLimit
	c.DayOfWeek = dayOfWeek
	c.StartTime = startTime

	pathAlgorithm := rga.RGA{}

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	locations := solver.Points.GetAllPoints()
	intervalRoute := make(map[int]generic.Point)

	intervalRoute[startID] = locations[startID]
	intervalRoute[endID] = locations[endID]

	intervalKeys := []int{startID, endID}

	pathAlgorithm = pathAlgorithm.Init(&solver).(rga.RGA)
	solver.Algorithm = pathAlgorithm.SetInitialRoute(intervalRoute, intervalKeys)

	sc.StartID = startID
	sc.EndID = endID
	solver.Score = sc

	c.StartID = startID
	c.EndID = endID
	solver.Constraints = c

	result, order, _ := solver.NextInterval()
	finalRoute := make(map[int]generic.Point)
	finalOrder := []int{startID}

	for _, k := range order {
		finalRoute[k] = result[k]
		finalOrder = append(finalOrder, k)
	}

	finalOrder = append(finalOrder, sc.EndID)
	finalRoute[sc.EndID] = result[sc.EndID]
	finalRoute[sc.StartID] = result[sc.StartID]

	// Commented part prints time of visit in final route
	/*
		currentTime := c.StartTime
		for i, k := range finalOrder {
			if i > 0 {
				visit := finalRoute[finalOrder[i-1]].(points.BaseLocation).Duration + points.WalkingTime(finalRoute[finalOrder[i-1]].(points.BaseLocation), finalRoute[k].(points.BaseLocation))
				currentTime = c.TimeUpdate(currentTime, visit)
			}
			fmt.Println(k, finalRoute[k].(points.BaseLocation).Title, currentTime)
		}
	*/

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
	return
}

// SolveOPFP solves Orienteering Problem with Functional Profits.
// Resulting path is constructed with respect of location position in final route.
func SolveOPFP(configPath string, startID int, endID int, timeLimit int, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}

	locs := points.BaseLocations{}
	sc := score.OPFPScore{}
	pathAlgorithm := &aco.ACO{}

	c := &constraints.OPFPConstraints{}

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

	finalScore = solver.Score.RouteScore(finalRoute, finalOrder)
	timePath = c.FinalRouteTime(finalRoute, finalOrder)

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
	return
}

// SolveOPFPByNonameAlgorithm solves Orienteering Problem with Functional Profits.
// Resulting path is constructed with respect of location position in final route.
func SolveOPFPByNonameAlgorithm(configPath string, startID int, endID int, timeLimit int, fileName string) (finalScore float64, timePath int) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   rga.RGA{},
	}

	locs := points.BaseLocations{}
	sc := score.OPFPScore{}
	pathAlgorithm := rga.RGA{}

	c := &constraints.OPFPConstraints{}

	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	locations := solver.Points.GetAllPoints()
	intervalRoute := make(map[int]generic.Point)

	intervalRoute[startID] = locations[startID]
	intervalRoute[endID] = locations[endID]

	intervalKeys := []int{startID, endID}

	pathAlgorithm = pathAlgorithm.Init(&solver).(rga.RGA)
	solver.Algorithm = pathAlgorithm.SetInitialRoute(intervalRoute, intervalKeys)

	sc.StartID = startID
	sc.EndID = endID
	solver.Score = sc

	c.StartID = startID
	c.EndID = endID
	c.TimeLimit = timeLimit
	solver.Constraints = c

	result, order, _ := solver.NextInterval()

	finalScore = solver.Score.RouteScore(result, order)
	timePath = c.FinalRouteTime(result, order)

	/*
		for _, i := range order {
			fmt.Println(result[i].(points.BaseLocation).Title)
		}
	*/
	locs.WriteLocationsToJSON(result, order, fileName)
	return
}

// SolveOPFPConsideringCityBrand solves Orienteering Problem with Functional Profits.
// Resulting path contains locations which represent city brand.
func SolveOPFPConsideringCityBrand(configPath string, timeLimit int, startTime int, dayOfWeek int, fileName string) {
	solver := generic.Solver{
		Points:      points.BaseLocations{},
		Score:       score.SimpleScore{},
		Constraints: constraints.EnrichmentConstraints{},
		Algorithm:   aco.ACO{},
	}
	eatTime := 90

	locs := points.CityBrandLocations{}
	sc := score.CityBrandScore{}
	pathAlgorithm := &aco.ACO{}

	c := &constraints.CityBrandConstraints{}
	c.TimeLimit = timeLimit - eatTime
	c.DayOfWeek = dayOfWeek
	c.StartTime = startTime
	//c.ForbiddenLocations = []int{787, 849, 841, 757, 831, 781, 7, 826, 855, 791, 838, 102, 853, 857, 4, 630,
	//129, 843, 814, 628, 138, 27, 107, 566, 835, 741, 20, 5, 17, 689, 804, 96, 789, 48}
	//c.ForbiddenLocations = []int{215, 253, 163, 164, 249, 217, 167, 194, 169, 3, 197, 17, 223, 8, 15, 250, 2, 225, 208, 193, 248, 206, 237, 131, 1, 245, 192, 220, 239}

	//c.ForbiddenLocations = []int{8479, 9855, 9813, 9862, 9684, 9914, 9741, 9646, 9915}
	c.ForbiddenLocations = []int{14958, 14981, 15000, 15301, 15325, 15591, 14885, 14694, 14709, 14843, 15578, 15611, 586, 15321, 15135, 15484,
		15878, 15330, 15374, 15329, 930, 15313, 15287, 15473, 15528, 15579, 15868, 15259, 15866, 15469, 14874, 14875, 14893, 14935,
		15582, 15382, 15294, 15009, 14758, 11907, 11956, 13475, 14620, 543, 14944, 15262, 14898, 14974, 15034, 15026, 14720,
		14688, 14569, 14986, 14821, 14696, 14702, 15366, 11919, 15335, 15338, 15250}
	solver.Points = locs
	solver.Constraints = c
	solver.Score = sc
	solver.Algorithm = pathAlgorithm
	solver.Start(configPath)

	maxScore := 0.0
	startID := 0

	endID := 0

	tmpLocations := solver.Points.GetAllPoints()
	tmpSC := sc.Init(tmpLocations)

	for i, location := range tmpLocations {
		locationScore := tmpSC.(score.CityBrandScore).SinglePointScoreWithoutPositionDependance(location, i)
		// Saransk: && i != 215 && i != 197
		// Rostov: && i != 4
		//Kaliningrad: && i != 3 && i != 796
		//EKB: && i != 1380
		//Sochi: && i != 18 && i != 1686 && i != 2283 && i != 1556 && i != 2044
		//&& i != 10
		//NN && i != 1803
		if locationScore >= maxScore && i != 13 {
			maxScore = locationScore
			startID = i
		}
	}

	sc.StartID = startID
	solver.Score = sc

	c.StartID = startID
	solver.Constraints = c

	result, order, _ := solver.NextInterval()
	finalRoute := make(map[int]generic.Point)
	finalOrder := []int{startID}

	for _, k := range order {
		finalRoute[k] = result[k]
		fmt.Println(k, result[k].(points.CityBrandLocation).Title)
		finalOrder = append(finalOrder, k)
	}

	finalRoute[sc.StartID] = result[sc.StartID]

	fmt.Println("final score:", solver.Score.RouteScore(finalRoute, finalOrder), len(finalOrder))
	fmt.Println("final time:", c.FinalRouteTime(finalRoute, finalOrder))
	fmt.Println(sc.StartID)
	fmt.Println(result[sc.StartID])

	// ADD RESTARAUNT
	eatConstraints := &constraints.RestarauntsConstraints{}

	solver.Constraints = eatConstraints
	solver.Start(configPath)
	fmt.Println(finalOrder)

	locationsNumber := len(finalOrder)
	for i := 0; i < locationsNumber-1; i++ {

		startID = finalOrder[i]
		endID = finalOrder[i+1]
		walkedTime := c.FinalRouteTime(finalRoute, finalOrder[:i+1])

		sc.StartID = startID
		//sc.EndID = endID
		solver.Score = sc

		eatConstraints.StartID = startID
		eatConstraints.EndID = endID
		eatConstraints.TimeLimit = eatTime + finalRoute[startID].(points.CityBrandLocation).Duration + finalRoute[endID].(points.CityBrandLocation).Duration + 30
		eatConstraints.DayOfWeek = dayOfWeek
		eatConstraints.StartTime = eatConstraints.TimeUpdate(startTime, walkedTime)

		solver.Constraints = eatConstraints
		result, order, _ := solver.NextInterval()

		for _, k := range order {
			l := result[k].(points.CityBrandLocation)
			l.IntervalNumber = -1
			finalRoute[k] = l
			finalOrder = append(finalOrder, k)
			eatConstraints.ForbiddenLocations = append(eatConstraints.ForbiddenLocations, k)
		}
	}

	for i := 0; i < 10; i++ {

		startID = finalOrder[0]
		endID = finalOrder[locationsNumber-1]

		sc.StartID = startID
		solver.Score = sc

		eatConstraints.StartID = startID
		eatConstraints.EndID = endID
		eatConstraints.TimeLimit = eatTime
		eatConstraints.DayOfWeek = dayOfWeek
		eatConstraints.StartTime = startTime

		solver.Constraints = eatConstraints
		result, order, _ := solver.NextInterval()

		for _, k := range order {
			l := result[k].(points.CityBrandLocation)
			l.IntervalNumber = -1
			finalRoute[k] = l
			finalOrder = append(finalOrder, k)
			eatConstraints.ForbiddenLocations = append(eatConstraints.ForbiddenLocations, k)
		}
	}

	locs.WriteLocationsToJSON(finalRoute, finalOrder, fileName)
}
