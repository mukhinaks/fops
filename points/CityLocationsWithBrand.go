package points

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mukhinaks/fops/generic"
)

type CityBrandLocations struct {
	Points []CityBrandLocation
	solver *generic.Solver
}

type CityBrandLocation struct {
	AdditionalCategories     []string         `json:"additional_categories"`
	Address                  string           `json:"address"`
	Categories               []string         `json:"categories"`
	CityBrand                []string         `json:"city_brand"`
	Duration                 int              `json:"duration"`
	FacebookCheckins         float64          `json:"facebook_checkins"`
	FacebookRating           float64          `json:"facebook_rating"`
	FoursquareCheckinsCount  float64          `json:"foursquare_checkinsCount"`
	FoursquareRating         float64          `json:"foursquare_rating"`
	FoursquareRatingVotes    float64          `json:"foursquare_ratingVotes"`
	FoursquareUserCount      float64          `json:"foursquare_userCount"`
	Image                    string           `json:"image"`
	InstagramTitle           string           `json:"instagram_title"`
	InstagramVisitorsNumber  float64          `json:"instagram_visitorsNumber"`
	Lat                      float64          `json:"lat"`
	Lng                      float64          `json:"lng"`
	OpenHours                map[string][]int `json:"open_hours"`
	Title                    string           `json:"title"`
	TripAdvisorRating        float64          `json:"tripAdvisor_rating"`
	TripAdvisorReviewsNumber float64          `json:"tripAdvisor_reviewsNumber"`
	WikipediaPage            float64          `json:"wikipedia_page"`
	WikipediaTitle           string           `json:"wikipedia_title"`
	X                        float64          `json:"x"`
	Y                        float64          `json:"y"`
	IntervalNumber           int
}

func (locations CityBrandLocations) Init(solver *generic.Solver) generic.Points {
	locations.solver = solver
	data, err := locations.readLocations(solver.Configuration["DataPath"].(string))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	locations.Points = data
	fmt.Println(len(data))
	return locations
}

func (l *CityBrandLocation) String() (string, error) {
	return string(l.Title), nil
}

func (locations CityBrandLocations) readLocations(pathToJSON string) ([]CityBrandLocation, error) {
	raw, err := ioutil.ReadFile(pathToJSON)
	if err != nil {
		return nil, err
	}

	data := make([]CityBrandLocation, 1)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (locations CityBrandLocations) GetAllPoints() []generic.Point {
	allLocations := make([]generic.Point, 0)
	for _, l := range locations.Points {
		allLocations = append(allLocations, l)
	}
	return allLocations
}

func (locations CityBrandLocations) GetCurrentPoints() map[int]generic.Point {
	currentLocations := make(map[int]generic.Point)
	for idx, location := range locations.Points {
		if locations.solver.Constraints.SinglePointConstraints(location, idx) {
			currentLocations[idx] = location
		}
	}
	return currentLocations
}

func (locations CityBrandLocations) GetPointsInArea(startID int, endID int) map[int]generic.Point {
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

func (l CityBrandLocations) WriteLocationsToJSON(route map[int]generic.Point, order []int, filePath string) {
	locations := make([]CityBrandLocation, 0)
	for _, idx := range order {
		locations = append(locations, route[idx].(CityBrandLocation))
	}
	locationsJSON, _ := json.Marshal(locations)

	ioutil.WriteFile(filePath, locationsJSON, 0644)
}
