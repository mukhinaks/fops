package generic

type Constraints interface {
	Init(locations []Point) Constraints
	Boundary(route map[int]Point, orderOfPoints []int) bool
	SinglePointConstraints(place Point, id int) bool
	ReducePoints(route map[int]Point, orderOfLocations []int, locations map[int]Point) map[int]Point
	UpdateConstraint(route map[int]Point, orderOfPoints []int, locations []Point) Constraints
}
