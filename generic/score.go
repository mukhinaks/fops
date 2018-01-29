package generic

type Score interface {
	Init(points []Point) Score
	SinglePointScore(route map[int]Point, orderOfPoints []int, place Point, id int) float64
	RouteScore(route map[int]Point, orderOfPoints []int) float64
}
