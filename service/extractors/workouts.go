package extractors

import (
	"time"
	"sort"
	"log"
	"sync"
	"strings"
	"github.com/weissleb/peloton-tableau-connector/config"
	"github.com/weissleb/peloton-tableau-connector/service/peloservice"
	"github.com/luci/luci-go/common/runtime/goroutine"
)

var layout = "2006-01-02 15:04:05"

var (
	workoutsCache map[string]struct {
		expire   time.Time
		workouts Workouts
	}
)

func init() {
	workoutsCache = make(map[string]struct {
		expire   time.Time
		workouts Workouts
	})
}

// Figure out a way to not have to pass in HttpClientInteface here.
func ExtractCyclingWorkouts(client PelotonClient) (Workouts, error) {

	user := client.GetSessionUser()
	if workouts, ok := workoutsCache[user]; ok {
		if time.Now().Before(workouts.expire) {
			log.Print("returning non-expired workouts cache hit for user " + user)
			return workouts.workouts, nil
		}
		log.Print("deleting expired workouts cache hit for user " + user)
		delete(workoutsCache, user)
	}

	workouts := Workouts{}
	page := uint16(0)
	hasNext := true

	// The underlying Peloton API which is called by `peloservice.GetWorkouts` returns results with paging, so we
	// also will query workouts in pages.
	for ; hasNext; page++ {
		// Let's add a retry loop.
		_workouts, err := peloservice.GetWorkouts(*client.getHttpClient(), *client.getUserSession(), page)
		if err != nil {
			return workouts, err
		}
		if config.LogLevel == "DEBUG" {
			log.Printf("DEBUG: got workouts for page %v of %v", page, _workouts.PageCount-1)
		}
		hasNext = _workouts.HasNext
		// We can set this config to false if we want to test by only getting a single page.
		// Working through all the pages could be a little slow because of subsequent peloservice for ride data.
		if !config.PeloAllPages {
			log.Print("only getting first page")
			hasNext = false
		}

		if len(_workouts.Workouts) == 0 {
			return nil, nil
		}

		for _, w := range _workouts.Workouts {
			if w.Wtype != "cycling" {
				continue
			}
			workout := Workout{
				ExtractTimeUTC:   time.Now().UTC(),
				Id:               w.Id,
				StartTimeSeconds: w.StartTimeSeconds,
				StartTimeUTC:     time.Unix(int64(w.StartTimeSeconds), 0).UTC(),
				TimeZone:         w.Timezone,
				Output:           w.Output,
				WasPR:            false, // set later
				HasWeights:       false, // set later
				CurrentPR:        w.CurrentPr,
				Type:             w.Wtype,
				Status:           w.Status,
				RideId:           w.Ride.Id,
				RideTitle:        w.Ride.Title,
				DurationSeconds:  w.Ride.Duration,
				RideDifficulty:   w.Ride.Difficulty_Rating,
				RideLevel:        w.Ride.Difficulty_Level,
			}
			loc, _ := time.LoadLocation(w.Timezone)
			workout.StartTime = time.Unix(int64(w.StartTimeSeconds), 0).In(loc)

			workouts = append(workouts, workout)
		}

		var waitGroup sync.WaitGroup
		for i, _ := range workouts {
			waitGroup.Add(1)
			go func(i int) {
				defer waitGroup.Done()
				rideId := workouts[i].RideId
				if config.LogLevel == "DEBUG" {
					log.Printf("DEBUG: (goroutine %v) getting ride detail for %s", goroutine.CurID(), rideId)
				}
				// Let's add a retry loop.
				ride, err := peloservice.GetRide(*client.getHttpClient(), *client.userSession, rideId)
				if err != nil {
					log.Fatal(err)
				}
				workouts[i].Instructor = ride.Instructor.Name
				if config.LogLevel == "DEBUG" {
					log.Printf("instructor is %s", ride.Instructor.Name)
				}
				workouts[i].InstructorImageURL = ride.Instructor.ImageURL
				// set HasWeights
				for _, equipmenttag := range ride.Equipmenttags {
					if strings.Contains(equipmenttag.Slug, "weights") {
						workouts[i].HasWeights = true
					}
				}
			}(i)
		}

		waitGroup.Wait()
	}

	// sort the workouts
	if config.LogLevel == "DEBUG" {
		log.Printf("DEBUG: sorting %d workouts", workouts.Len())
	}
	sort.Sort(workouts)

	// iterate the workouts and set WasPR
	prMap := make(map[uint16]float32)
	for i, workout := range workouts {
		if prMap[workout.DurationSeconds] == 0 {
			prMap[workout.DurationSeconds] = workout.Output
			continue
		}

		if prMap[workout.DurationSeconds] < workout.Output {
			prMap[workout.DurationSeconds] = workout.Output
			workouts[i].WasPR = true
			//log.Printf("Got PR for %d min ride on %s.", workout.DurationSeconds/60, workout.StartTime.Format(layout))
		}
	}

	if config.UseWorkoutCache {
		log.Print("caching workouts for user " + user)
		workoutsCache[user] = struct {
			expire   time.Time
			workouts Workouts
		}{
			expire:   time.Now().Add(time.Second * time.Duration(config.CacheExpireSeconds)),
			workouts: workouts,
		}
	}

	if config.LogLevel == "DEBUG" {
		for i, _ := range workouts {
			log.Printf("workout %s has instructor %s", workouts[i].Id, workouts[i].Instructor)
		}
	}

	return workouts, nil
}

func GetCyclingWorkoutsSummary(client PelotonClient) (WorkoutsSummary, error) {
	workouts, err := ExtractCyclingWorkouts(client)
	if err != nil {
		return WorkoutsSummary{}, err
	}

	workoutCount := workouts.Len()
	var lastWorkout time.Time
	for _, workout := range workouts {
		if lastWorkout.Before(workout.StartTimeUTC) {
			lastWorkout = workout.StartTimeUTC
		}
	}

	return WorkoutsSummary{
		TotalWorkouts:           workoutCount,
		LastWorkoutStartTimeUTC: lastWorkout,
	}, nil
}
