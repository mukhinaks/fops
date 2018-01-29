package generic

type PathAlgorithm interface {
	Init(solver *Solver) PathAlgorithm
	CreateRoute() (map[int]Point, []int, float64)
}
