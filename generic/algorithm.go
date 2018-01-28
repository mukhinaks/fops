package generic

type PathAlgorithm interface {
	Init(solver *Solver) PathAlgorithm
	CreateRoute() (map[int]Location, []int, float64)
	GetRawLocations() []Location
}
