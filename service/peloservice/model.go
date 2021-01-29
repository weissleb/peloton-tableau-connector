package peloservice

// Contains structs for workouts.
//
// Users should first obtain a list of `Workouts`s.  For each workout, the ride can then be obtained
// from the `WorkoutDetail`.  More information is then available about the ride.
//
// 1. Obtain list of `Workouts`s.
// 2. For each workout, get the `Ride` from `WorkoutDetail`.
// 3. Obtain the `RideDetail`.

type UserSession struct {
	UserId  string
	Cookies string
}

type session struct {
	SessionId string `json:"session_id"`
	UserId    string `json:"user_id"`
}

// List of workouts exported from /api/user/{userId}/workout_history_csv
type exportedWorkouts []exportedWorkout

type exportedWorkout struct {
	StartTime          string  `csv:"Workout Timestamp"`
	LiveOrOnDemand     string  `csv:"Live/On-Demand"`
	Instructor         string  `csv:"Instructor Name"`
	LengthMinutes      int     `csv:"Length (minutes)"`
	FitnessDiscipline  string  `csv:"Fitness Discipline"`
	ClassType          string  `csv:"Type"`
	ClassTitle         string  `csv:"Title"`
	TotalOutput        int     `csv:"Total Output"`
	AvgWatts           int     `csv:"Avg. Watts"`
	AvgResistance      string  `csv:"Avg. Resistance"`
	AvgCadenceRPM      int     `csv:"Avg. Cadence (RPM)"`
	AvgSpeedMPH        float64 `csv:"Avg. Speed (mph)"`
	AvgSpeedKPH        float64 `csv:"Avg. Speed (kph)"`
	DistanceMiles      float64 `csv:"Distance (mi)"`
	DistanceKilometers float64 `csv:"Distance (km)"`
	CaloriesBurned     int     `csv:"Calories Burned"`
	AvgHeartRate       float64 `csv:"Avg. Heartrate"`
}

// List of workouts returned from /api/user/{userId}/workouts
type workouts struct {
	Workouts  []workout `json:"data"`
	PageCount uint16    `json:"page_count"`
	HasNext   bool      `json:"show_next"`
}

type workout struct {
	Id               string  `json:"id"`
	StartTimeSeconds uint64  `json:"start_time"`
	Timezone         string  `json:"timezone"`
	CurrentPr        bool    `json:"is_total_work_personal_record"`
	Wtype            string  `json:"fitness_discipline"`
	Status           string  `json:"status"`
	Output           float64 `json:"total_work"`
	Ride             ride    `json:"ride"`
}

// Workouts details returned from /api/workout/{id}
// Struct removed as it was not needed.

type ride struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Duration   uint16 `json:"duration"`
	Instructor struct {
		Name     string `json:"name"`
		ImageURL string `json:"about_image_url"`
	} `json:"instructor"`
	Difficulty_Rating float64 `json:"difficulty_rating_avg"`
	Difficulty_Level  string  `json:"difficulty_level"`
	Equipmenttags     []struct {
		Slug string `json:"slug"`
	} `json:"equipment_tags"`
}

// Ride details returned from /api/ride/{id}/details
type rideDetail struct {
	Ride ride `json:"ride"`
}

// Summaries such as distance, calories, etc. are available at https://api.onepeloton.com/api/workout/{id}/performance_graph?every_n=1800
