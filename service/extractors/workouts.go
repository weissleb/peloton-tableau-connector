package extractors

import (
	"time"
	"sort"
	"log"
	"github.com/weissleb/peloton-tableau-connector/config"
	"github.com/weissleb/peloton-tableau-connector/service/peloservice"
	"strings"
	"strconv"
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

// Gets dataset from the peloservice and transforms it for our needs.
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
	extractTime := time.Now().UTC()
	exportedWorkouts, err := peloservice.GetExportedWorkouts(*client.getHttpClient(), *client.getUserSession())
	if err != nil {
		return workouts, err
	}
	if len(exportedWorkouts) == 0 {
		return nil, nil
	}

	layout := "2006-01-02 15:04 (MST)"
	for _, w := range exportedWorkouts {
		if w.FitnessDiscipline != "Cycling" {
			continue
		}
		startTime, _ := time.Parse(layout, w.StartTime)
		timeZone, _ := startTime.Zone()
		avgResistenceInt, _ := strconv.Atoi(strings.TrimRight(w.AvgResistance, "%s"))
		avgResistence := float64(avgResistenceInt) / 100.00
		workouts = append(workouts, Workout{
			ExtractTimeUTC:     extractTime,
			StartTime:          startTime,
			TimeZone:           timeZone,
			StartTimeUTC:       startTime.UTC(),
			Type:               w.ClassType,
			RideTitle:          w.ClassTitle,
			Instructor:         w.Instructor,
			RideLengthMinutes:  w.LengthMinutes,
			Output:             w.TotalOutput,
			AvgWatts:           w.AvgWatts,
			AvgResistance:      avgResistence,
			AvgCadenceRPM:      w.AvgCadenceRPM,
			AvgSpeedMPH:        w.AvgSpeedMPH,
			AvgSpeedKPH:        w.AvgSpeedKPH,
			DistanceMiles:      w.DistanceMiles,
			DistanceKilometers: w.DistanceKilometers,
			CaloriesBurned:     w.CaloriesBurned,
			AvgHeartRate:       w.AvgHeartRate,
		})
	}

	/*
		 * This section pulls the workouts from the API, rather than using the CSV export.
		 * I may add some of this back later to assign the Id, RideDifficulty, RideLevel, InstructorImageURL,
		 * and HasWeights fields.

		page := uint16(0)
		hasNext := true

		// The underlying Peloton API which is called by `peloservice.GetWorkouts` returns results with paging, so we
		// also will query workouts in pages.
		for ; hasNext; page++ {
			// Let's add a retry loop.
			exportedWorkouts, err := peloservice.GetWorkouts(*client.getHttpClient(), *client.getUserSession(), page)
			if err != nil {
				return workouts, err
			}
			if config.LogLevel == "DEBUG" {
				log.Printf("DEBUG: got workouts for page %v of %v", page, exportedWorkouts.PageCount-1)
			}
			hasNext = exportedWorkouts.HasNext
			// We can set this config to false if we want to test by only getting a single page.
			// Working through all the pages could be a little slow because of subsequent peloservice for ride data.
			if !config.PeloAllPages {
				log.Print("only getting first page")
				hasNext = false
			}

			if len(exportedWorkouts.Workouts) == 0 {
				return nil, nil
			}

			for _, w := range exportedWorkouts.Workouts {
				if w.Wtype != "cycling" {
					continue
				}
				// Gather up workouts so we can later assign Id, RideDifficulty, RideLevel,
				// InstructorImageURL and then HasWeights.
			}

			// The following is really only necessary to get the HasWeights flag.
			// It gets instructor URL too, but meh.

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
	*/

	// sort the workouts by StartTime
	// iterate the workouts and set WasPR
	sort.Sort(workouts)
	prMap := make(map[int]int)
	for i, workout := range workouts {
		if prMap[workout.RideLengthMinutes] == 0 {
			prMap[workout.RideLengthMinutes] = workout.Output
			continue
		}

		if prMap[workout.RideLengthMinutes] < workout.Output {
			prMap[workout.RideLengthMinutes] = workout.Output
			workouts[i].WasPR = true
			//log.Printf("Got PR for %d min ride on %s.", workout.DurationSeconds/60, workout.StartTime.Format(layout))
		}
	}

	// sort workouts in reverse order by StartTime
	// find the first (i.e. last) WasPR, and set it to CurrentPR
	sort.Sort(sort.Reverse(workouts))
	crSet := make(map[int]bool)
	for i, workout := range workouts {
		if crSet[workout.RideLengthMinutes] {
			continue
		}
		if workouts[i].WasPR {
			workouts[i].CurrentPR = true
			crSet[workout.RideLengthMinutes] = true
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
