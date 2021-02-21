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
	Id                string    `csv:"Id"`
	StartTime         time.Time `csv:"StartTime"`
	TimeZone          string    `csv:"TimeZone"`
	StartTimeUTC      time.Time `csv:"StartTimeUTC"`
	WasPR             bool      `csv:"WasPR"`
	CurrentPR         bool      `csv:"CurrentPR"`
	Type              string    `csv:"Type"`
	RideTitle         string    `csv:"RideTitle"`
	RideDifficulty    float64   `csv:"RideDifficulty"`
	Instructor        string    `csv:"Instructor"`
	RideImageUrl      string    `csv:"RideImageUrl"`
	RideLengthMinutes int       `csv:"RideLengthMinutes"`
	//HasWeights         bool      `csv:"HasWeights"`
	Output             int     `csv:"Output"`
	AvgWatts           int     `csv:"AvgWatts"`
	AvgResistance      float64 `csv:"AvgResistance"`
	AvgCadenceRPM      int     `csv:"AvgCadence"`
	AvgSpeedMPH        float64 `csv:"AvgSpeedMPH"`
	AvgSpeedKPH        float64 `csv:"AvgSpeedKPH"`
	DistanceMiles      float64 `csv:"DistanceMiles"`
	DistanceKilometers float64 `csv:"DistanceKilometers"`
	CaloriesBurned     int     `csv:"CaloriesBurned"`
	AvgHeartRate       float64 `csv:"AvgHeartRate"`
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
			workout.StartTime.Format(layout),
			workout.TimeZone,
			workout.StartTimeUTC.Format(layout),
			fmt.Sprintf("%d", boolToDig(workout.WasPR)),
			fmt.Sprintf("%d", boolToDig(workout.CurrentPR)),
			workout.Type,
			workout.RideTitle,
			fmt.Sprintf("%.2f", workout.RideDifficulty),
			workout.Instructor,
			workout.RideImageUrl,
			fmt.Sprintf("%d", workout.RideLengthMinutes),
			//fmt.Sprintf("%d", boolToDig(workout.HasWeights)),
			fmt.Sprintf("%d", workout.Output),
			fmt.Sprintf("%d", workout.AvgWatts),
			fmt.Sprintf("%.2f", workout.AvgResistance),
			fmt.Sprintf("%d", workout.AvgCadenceRPM),
			fmt.Sprintf("%.2f", workout.AvgSpeedMPH),
			fmt.Sprintf("%.2f", workout.AvgSpeedKPH),
			fmt.Sprintf("%.2f", workout.DistanceMiles),
			fmt.Sprintf("%.2f", workout.DistanceKilometers),
			fmt.Sprintf("%d", workout.CaloriesBurned),
			fmt.Sprintf("%.2f", workout.AvgHeartRate),
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
	return w[i].StartTimeUTC.Before(w[j].StartTimeUTC)
}

func boolToDig(boolVal bool) int8 {
	if boolVal {
		return 1
	}
	return 0
}
