package generic

type Score interface {
	Init(locations []Location) Score
	LocationScore(route map[int]Location, orderOfLocations []int, place Location, id int) float64
	RouteScore(route map[int]Location, orderOfLocations []int) float64
}
