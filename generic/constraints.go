package generic

type Constraints interface {
	Init(locations []Location) Constraints
	Boundary(route map[int]Location, orderOfLocations []int) bool
	LocationConstraints(place Location, id int) bool
}
