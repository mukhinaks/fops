package points

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/mukhinaks/fops/generic"
)

type BaseLocations struct {
	Points []BaseLocation
	solver *generic.Solver
}

type BaseLocation struct {
	Address                 string   `json:"address"`
	Category                []string `json:"category"`
	Duration                int      `json:"duration"`
	FoursquareCheckinsCount float64  `json:"foursquare_checkinsCount"`
	FoursquareRating        float64  `json:"foursquare_rating"`
	FoursquareRatingVotes   float64  `json:"foursquare_ratingVotes"`
	FoursquareUserCount     float64  `json:"foursquare_userCount"`
	InstagramVisitorsList   []string `json:"instagram_visitorsList"`
	InstagramVisitorsNumber float64  `json:"instagram_visitorsNumber"`
	Lat                     float64  `json:"lat"`
	Lng                     float64  `json:"lng"`
	OfficialGuide           float64  `json:"officialGuide"`

	OpenHours struct {
		Num0 []int `json:"0"`
		Num1 []int `json:"1"`
		Num2 []int `json:"2"`
		Num3 []int `json:"3"`
		Num4 []int `json:"4"`
		Num5 []int `json:"5"`
		Num6 []int `json:"6"`
	} `json:"open_hours"`

	Title                    string  `json:"title"`
	TripAdvisorLink          string  `json:"tripAdvisor_link"`
	TripAdvisorRating        float64 `json:"tripAdvisor_rating"`
	TripAdvisorReviewsNumber float64 `json:"tripAdvisor_reviewsNumber"`
	X                        float64 `json:"x"`
	Y                        float64 `json:"y"`
	ID                       int     `json:"id"`
}

func (locations BaseLocations) Init(solver *generic.Solver) generic.Points {
	locations.solver = solver
	data, err := readLocations(solver.Configuration["DataPath"].(string))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	locations.Points = data
	return locations
}

func (l *BaseLocation) String() (string, error) {
	return string(l.Title), nil
}

func readLocations(pathToJSON string) ([]BaseLocation, error) {
	raw, err := ioutil.ReadFile(pathToJSON)
	if err != nil {
		return nil, err
	}

	data := make([]BaseLocation, 1)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (locations BaseLocations) GetAllPoints() []generic.Point {
	allLocations := make([]generic.Point, 0)
	for _, l := range locations.Points {
		allLocations = append(allLocations, l)
	}
	return allLocations
}

func (locations BaseLocations) GetCurrentPoints() map[int]generic.Point {
	currentLocations := make(map[int]generic.Point)
	for idx, location := range locations.Points {
		if locations.solver.Constraints.SinglePointConstraints(location, idx) {
			currentLocations[idx] = location
		}
	}
	return currentLocations
}

func EuclidianDistance(loc1 BaseLocation, loc2 BaseLocation) float64 {
	return distanceToPoint(loc1.X, loc1.Y, loc2.X, loc2.Y)
}

func EuclidianDistanceToLineSegment(start BaseLocation, end BaseLocation, location BaseLocation) float64 {
	return distanceToLine(start.X, start.Y, end.X, end.Y, location.X, location.Y)
}

func WalkingTime(loc1 BaseLocation, loc2 BaseLocation) int {
	return int(distanceToPoint(loc1.X, loc1.Y, loc2.X, loc2.Y) / 66.7)
}

func distanceToPoint(loc1Lat float64, loc1Lng float64, loc2Lat float64, loc2Lng float64) float64 {
	result := math.Sqrt(math.Pow((loc1Lat-loc2Lat), 2) + math.Pow((loc1Lng-loc2Lng), 2))
	return result
}

func scalarDot(loc1Lat float64, loc1Lng float64, loc2Lat float64, loc2Lng float64) float64 {
	result := (loc1Lat * loc2Lat) + (loc1Lng * loc2Lng)
	return result
}

func distanceToLine(loc1Lat float64, loc1Lng float64, loc2Lat float64, loc2Lng float64, newLat float64, newLng float64) float64 {
	vector1Lat := loc1Lat - loc2Lat
	vector1Lng := loc1Lng - loc2Lng
	vector2Lat := loc1Lat - newLat
	vector2Lng := loc1Lng - newLng

	scalarProduct := scalarDot(vector1Lat, vector1Lng, vector2Lat, vector2Lng)
	if scalarProduct <= 0 {
		return distanceToPoint(newLat, newLng, loc1Lat, loc1Lng)
	}

	length := scalarDot(vector1Lat, vector1Lng, vector1Lat, vector1Lng)
	if length <= scalarProduct {
		return distanceToPoint(newLat, newLng, loc2Lat, loc2Lng)
	}

	b := scalarProduct / length
	locLat := loc1Lat + b*vector1Lat
	locLng := loc1Lng + b*vector1Lng

	return distanceToPoint(newLat, newLng, locLat, locLng)
}

func (l BaseLocations) WriteLocationsToJSON(route map[int]generic.Point, order []int, filePath string) {
	locations := make([]BaseLocation, 0)
	for _, idx := range order {
		locations = append(locations, route[idx].(BaseLocation))
	}
	locationsJSON, _ := json.Marshal(locations)

	ioutil.WriteFile(filePath, locationsJSON, 0644)
}
