package peloservice

import (
	"testing"
	"encoding/json"
	"github.com/gocarina/gocsv"
)

func TestJson(t *testing.T) {
	workoutsResponse := []byte(`
{
  "data": [
    {
      "id": "4af299ce8a794ae29996b762f353fab4",
      "is_total_work_personal_record": false,
      "start_time": 1599492654,
      "status": "COMPLETE",
      "workout_type": "cycling",
      "total_work": 482432.95
    }
  ]
}`)

	rideDetailResponse := []byte(`
{
  "ride": {
    "equipment_tags": [
      {
		"id": "0f5f1ff2d6c647cf98d599ed90ad72d3",
		"name": "Light Weights",
		"slug": "light_weights",
		"icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/equipment-icons/light_weights.png"
	  }
    ]
  }
}
`)

	var (
		err        error
		workouts   workouts
		rideDetail rideDetail
		want       interface{}
		got        interface{}
	)

	if err = json.Unmarshal(workoutsResponse, &workouts); err != nil {
		t.Fatal(err)
	}

	if err = json.Unmarshal(rideDetailResponse, &rideDetail); err != nil {
		t.Fatal(err)
	}

	want = "light_weights"
	got = rideDetail.Ride.Equipmenttags[0].Slug
	if got != want {
		t.Fatalf("got: %s; want %s\n", got, want)
	}
}

func TestCsvUS(t *testing.T) {
	workoutsResponse := []byte(`
Workout Timestamp,Live/On-Demand,Instructor Name,Length (minutes),Fitness Discipline,Type,Title,Class Timestamp,Total Output,Avg. Watts,Avg. Resistance,Avg. Cadence (RPM),Avg. Speed (mph),Distance (mi),Calories Burned,Avg. Heartrate,Avg. Incline,Avg. Pace (min/mi)
2018-02-22 17:25 (EST),Live,,15,Cycling,Scenic Ride,15 min Venice Scenic Ride,,46,52,25%,93,11.87,2.96,63,,,
`)

	var (
		err        error
		workouts   exportedWorkouts
		want       interface{}
		got        interface{}
	)

	if err = gocsv.UnmarshalBytes(workoutsResponse, &workouts); err != nil {
		t.Fatal(err)
	}

	want = "15 min Venice Scenic Ride"
	got = workouts[0].ClassTitle
	if got != want {
		t.Fatalf("got: %s; want %s\n", got, want)
	}
}

func TestCsvMetric(t *testing.T) {
	workoutsResponse := []byte(`
Workout Timestamp,Live/On-Demand,Instructor Name,Length (minutes),Fitness Discipline,Type,Title,Class Timestamp,Total Output,Avg. Watts,Avg. Resistance,Avg. Cadence (RPM),Avg. Speed (kph),Distance (km),Calories Burned,Avg. Heartrate,Avg. Incline,Avg. Pace (min/km)
2020-12-15 21:18 (CET),On Demand,Irene Scholz,30,Cycling,Theme,30 min Festive Ride,2020-12-15 19:24 (CET),284,158,44%,87,31.29,15.63,384,,,
`)

	var (
		err        error
		workouts   exportedWorkouts
		want       interface{}
		got        interface{}
	)

	if err = gocsv.UnmarshalBytes(workoutsResponse, &workouts); err != nil {
		t.Fatal(err)
	}

	want = 15.63
	got = workouts[0].DistanceKilometers
	if got != want {
		t.Fatalf("got: %s; want %s\n", got, want)
	}
}
