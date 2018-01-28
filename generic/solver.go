package generic

import (
	"github.com/mukhinaks/fops/misc"
)

type Solver struct {
	Algorithm     PathAlgorithm
	Score         Score
	Locations     Locations
	Constraints   Constraints
	Configuration map[string]interface{}
}

func (solver *Solver) Start(configPath string) {
	solver.Configuration = misc.ReadConfig(configPath)
	solver.Locations = solver.Locations.Init(solver)
	locations := solver.Locations.GetAllLocations()
	solver.Algorithm = solver.Algorithm.Init(solver)
	solver.Score = solver.Score.Init(locations)
	solver.Constraints = solver.Constraints.Init(locations)
}

func (solver *Solver) NextInterval() (map[int]Location, []int, float64) {
	locations := solver.Locations.GetAllLocations()
	solver.Score = solver.Score.Init(locations)
	solver.Constraints = solver.Constraints.Init(locations)
	return solver.Algorithm.CreateRoute()
}
