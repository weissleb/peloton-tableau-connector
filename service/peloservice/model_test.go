package peloservice

import (
	"testing"
	"encoding/json"
)

func Test(t *testing.T) {
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
