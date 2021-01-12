package extractors

import (
	"github.com/weissleb/peloton-tableau-connector/service/clients"
	"time"
	"fmt"
	"reflect"
	"github.com/weissleb/peloton-tableau-connector/service/peloservice"
)

type ClientSession struct {
	Client  clients.HttpClientInterface
	Session peloservice.UserSession
}

type WorkoutsSummary struct {
	TotalWorkouts           int       `json:"TotalWorkouts"`
	LastWorkoutStartTimeUTC time.Time `json:"LastWorkoutStartTimeUTC"`
}

// structs for user friendly workouts
type Workouts []Workout

type Workout struct {
	ExtractTimeUTC     time.Time `csv:"ExtractTimeUTC"`
	Id                 string    `csv:"Id"`
	StartTimeSeconds   uint64    `csv:"StartTimeSeconds"`
	TimeZone           string    `csv:"TimeZone"`
	StartTimeUTC       time.Time `csv:"StartTimeUTC"`
	StartTime          time.Time `csv:"StartTime"`
	Output             float32   `csv:"Output"`
	WasPR              bool      `csv:"WasPR"`
	CurrentPR          bool      `csv:"CurrentPR"`
	Type               string    `csv:"Type"`
	Status             string    `csv:"Status"`
	RideId             string    `csv:"RidId"`
	RideTitle          string    `csv:"RideTitle"`
	RideDifficulty     float32   `csv:"RideDifficulty"`
	RideLevel          string    `csv:"RideLevel"`
	Instructor         string    `csv:"Instructor"`
	InstructorImageURL string    `csv:"InstructorImageURL"`
	DurationSeconds    uint16    `csv:"DurationSeconds"`
	HasWeights         bool      `csv:"HasWeights"`
}

func (w Workouts) GetAsRecords(withHeader bool) [][]string {
	var records [][]string
	if withHeader {
		t := reflect.TypeOf(Workout{})
		header := make([]string, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			header[i] = t.Field(i).Name
		}
		records = append(records, header)
	}

	layout := "2006-01-02 15:04:05"
	for _, workout := range w {
		records = append(records, []string{
			workout.ExtractTimeUTC.Format(layout),
			workout.Id,
			fmt.Sprintf("%d", workout.StartTimeSeconds),
			workout.TimeZone,
			workout.StartTimeUTC.Format(layout),
			workout.StartTime.Format(layout),
			fmt.Sprintf("%f", workout.Output),
			fmt.Sprintf("%d", boolToDig(workout.WasPR)),
			fmt.Sprintf("%d", boolToDig(workout.CurrentPR)),
			workout.Type,
			workout.Status,
			workout.RideId,
			workout.RideTitle,
			fmt.Sprintf("%f", workout.RideDifficulty),
			workout.RideLevel,
			workout.Instructor,
			workout.InstructorImageURL,
			fmt.Sprintf("%d", workout.DurationSeconds),
			fmt.Sprintf("%d", boolToDig(workout.HasWeights)),
		})
	}

	return records
}

// Len is part of sort.Interface.
func (w Workouts) Len() int {
	return len(w)
}

// Swap is part of sort.Interface.
func (w Workouts) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// Less is part of sort.Interface. We use StartTimeSeconds as the value to sort by.
func (w Workouts) Less(i, j int) bool {
	return w[i].StartTimeSeconds < w[j].StartTimeSeconds
}

func boolToDig(boolVal bool) int8 {
	if boolVal {
		return 1
	}
	return 0
}
