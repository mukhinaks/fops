package generic

type Constraints interface {
	Init(locations []Point) Constraints
	Boundary(route map[int]Point, orderOfPoints []int) bool
	SinglePointConstraints(place Point, id int) bool
}
