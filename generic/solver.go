package generic

import (
	"github.com/mukhinaks/fops/misc"
)

type Solver struct {
	Algorithm     PathAlgorithm
	Score         Score
	Points        Points
	Constraints   Constraints
	Configuration map[string]interface{}
}

func (solver *Solver) Start(configPath string) {
	solver.Configuration = misc.ReadConfig(configPath)
	solver.Points = solver.Points.Init(solver)
	points := solver.Points.GetAllPoints()
	solver.Algorithm = solver.Algorithm.Init(solver)
	solver.Score = solver.Score.Init(points)
	solver.Constraints = solver.Constraints.Init(points)
}

func (solver *Solver) NextInterval() (map[int]Point, []int, float64) {
	points := solver.Points.GetAllPoints()
	solver.Score = solver.Score.Init(points)
	solver.Constraints = solver.Constraints.Init(points)
	return solver.Algorithm.CreateRoute()
}
