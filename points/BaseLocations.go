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

	OpenHours map[string][]int `json:"open_hours"`

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
	data, err := locations.readLocations(solver.Configuration["DataPath"].(string))
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

func (locations BaseLocations) readLocations(pathToJSON string) ([]BaseLocation, error) {
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

func (l BaseLocations) WriteLocationsToJSON(route map[int]generic.Point, order []int, filePath string) {
	locations := make([]BaseLocation, 0)
	for _, idx := range order {
		locations = append(locations, route[idx].(BaseLocation))
	}
	locationsJSON, _ := json.Marshal(locations)

	ioutil.WriteFile(filePath, locationsJSON, 0644)
}

func (locations BaseLocations) GetPointsInArea(startID int, endID int) map[int]generic.Point {
	start := locations.Points[startID]
	end := locations.Points[endID]
	distance := EuclidianDistance(start, end)

	currentLocations := make(map[int]generic.Point)
	for idx, location := range locations.Points {
		if EuclidianDistance(location, start) <= distance && EuclidianDistance(location, end) <= distance {
			currentLocations[idx] = location
		}
	}
	return currentLocations
}

func (locations BaseLocations) FindClosestPoint(lat float64, lon float64) generic.Point {
	minDistance := math.MaxFloat64
	closestPoint := locations.Points[0]

	for _, location := range locations.Points {
		if HaversineDistance(location.Lat, location.Lng, lat, lon) <= minDistance {
			closestPoint = location
		}
	}
	return closestPoint
}
