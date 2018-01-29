package generic

type Point interface {
}

type Points interface {
	Init(solver *Solver) Points
	GetAllPoints() []Point
	GetCurrentPoints() map[int]Point
}
