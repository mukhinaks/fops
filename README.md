# fops

## How to launch:

Clone the repository:\
`git clone https://github.com/mukhinaks/fops.git`

Edit *config.json* file if neccessary:\
AntsNumber - number of ants per point\
Fadeness - pheromone fadeness, for further information read [algorithm description](https://en.wikipedia.org/wiki/Ant_colony_optimization_algorithms) \
Iterations - number of iterations for ant colony optimization\
AttractivenessControl - influence of point score in probability computation\
PheromoneControl - influence of pheromone value in probability computation\
DataPath - path to dataset\
NumberOfChannels - parameter for parallel launch\
TimeLimit - currently not used\

Build the framework (run terminal in FOPS directory):\
`go build`

Start the itinerary construction:\
`./fops`

## Output result
The output route is JSON file, where each element contains all information about location.
```json
[
    {
        "address":"St. Petersburg, Isaakiyevskaya Square",
        "category":[
            "Sights \u0026 Landmarks",
            "Museums \u0026 Libraries"
        ],
        "duration":120,
        "foursquare_checkinsCount":14805,
        "foursquare_rating":9.4,
        "foursquare_ratingVotes":1238,
        "foursquare_userCount":16254,
        "instagram_visitorsList":null,
        "instagram_visitorsNumber":148834,
        "lat":59.933013,
        "lng":30.307442,
        "officialGuide":1,
        "open_hours":{
            "0":null,
            "1":[
              1030,
              1730
            ],
            "2":null,
            "3":[
              1030,
              1730
            ],
            "4":[
              1030,
              1730
            ],
            "5":[
              1030,
              1730
            ],
            "6":[
              1030,
              1730
            ]
        },
        "title":"St.Isaac's Square",
        "tripAdvisor_link":"https://www.tripadvisor.com/Attraction_Review-g298507-d300132-Reviews-St_Isaac_s_Cathedral_State_Museum_Memorial-St_Petersburg_Northwestern_District.html",
        "tripAdvisor_rating":4.5,
        "tripAdvisor_reviewsNumber":8622,
        "x":3373809.0106866737,
        "y":8384839.049077872,
        "id":0
    }
]
```

## Dataset
The data is publicly available [here](https://goo.gl/q9T2pr).
