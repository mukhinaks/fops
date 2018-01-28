package generic

type Location interface {
}

type Locations interface {
	Init(solver *Solver) Locations
	GetAllLocations() []Location
	GetCurrentLocations() map[int]Location
}
