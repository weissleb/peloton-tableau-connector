package extractors

import (
	"fmt"
	"github.com/weissleb/peloton-tableau-connector/config"
	"github.com/weissleb/peloton-tableau-connector/service/peloservice"
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

//var layout = "2006-01-02 15:04:05"

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

	workouts := Workouts{}
	user := client.GetSessionUser()

	if config.UseWorkoutCache {
		if workouts, ok := workoutsCache[user]; ok {
			if time.Now().Before(workouts.expire) {
				log.Print("returning non-expired workouts cache hit for user " + user)
				return workouts.workouts, nil
			}
			log.Print("deleting expired workouts cache hit for user " + user)
			delete(workoutsCache, user)
		}
	}

	type extraFields struct {
		Id                    string
		RideDifficulty        float64
		RideImageUrl          string
		StartTime             time.Time
		EndTime               time.Time
		TimeZone              string
		StriveScore           float64
		HeartRateZone1Seconds int
		HeartRateZone1Percent float64
		HeartRateZone2Seconds int
		HeartRateZone2Percent float64
		HeartRateZone3Seconds int
		HeartRateZone3Percent float64
		HeartRateZone4Seconds int
		HeartRateZone4Percent float64
		HeartRateZone5Seconds int
		HeartRateZone5Percent float64
	}
	apiWorkoutsMapAbbreviation := make(map[string]extraFields) // Old data come using abbreviation and no daylight savings :-(
	apiWorkoutsMapOffset := make(map[string]extraFields)       // Newer data come using offset.

	// Go get the workouts from the API.
	page := uint16(0)
	hasNext := true
	layout := "2006-01-02 15:04 (MST)"

	// The underlying Peloton API which is called by `peloservice.GetWorkouts` returns results with paging, so we
	// also will query workouts in pages.
	for ; hasNext; page++ {
		// TODO: Let's add a retry loop.
		apiWorkouts, err := peloservice.GetWorkouts(*client.getHttpClient(), *client.getUserSession(), page)
		if err != nil {
			return workouts, err
		}
		if config.LogLevel == "DEBUG" {
			log.Printf("DEBUG: got workouts for page %v of %v", page, apiWorkouts.PageCount-1)
		}
		hasNext = apiWorkouts.HasNext
		// We can set this config to false if we want to test by only getting a single page.
		// Working through all the pages could be a little slow because of subsequent peloservice for ride data.
		if !config.PeloAllPages {
			log.Print("only getting first page")
			hasNext = false
		}

		if len(apiWorkouts.Workouts) == 0 {
			return nil, nil
		}

		// The following is really only necessary to get the HasWeights flag.
		// It gets instructor URL too, but meh.
		/*
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
					workouts[i].RideImageUrl = ride.Instructor.ImageURL
					// set HasWeights
					for _, equipmenttag := range ride.Equipmenttags {
						if strings.Contains(equipmenttag.Slug, "weights") {
							workouts[i].HasWeights = true
						}
					}
				}(i)
			}

			waitGroup.Wait()
		*/

		dayLayout := "2006-01-02"
		minuteLayout := "2006-01-02 15:04"
		for _, workout := range apiWorkouts.Workouts {
			loc, _ := time.LoadLocation(workout.Timezone)
			st := time.Unix(int64(workout.StartTimeSeconds), 0).In(loc)
			et := time.Unix(int64(workout.EndTimeSeconds), 0).In(loc)
			z, _ := st.Zone()
			totalHRZone := workout.EffortZones.HeartRateZoneDurations.HrZone1Seconds +
				workout.EffortZones.HeartRateZoneDurations.HrZone2Seconds +
				workout.EffortZones.HeartRateZoneDurations.HrZone3Seconds +
				workout.EffortZones.HeartRateZoneDurations.HrZone4Seconds +
				workout.EffortZones.HeartRateZoneDurations.HrZone5Seconds
			// Store keys as day, output, title for matching start time provided with abbreviation and no daylight savings.
			key := fmt.Sprintf("%s %.0f %s",
				st.Format(dayLayout), math.Round(workout.Output/1000), strings.ToLower(workout.Ride.Title))
			apiWorkoutsMapAbbreviation[key] = extraFields{
				Id:                    workout.Id,
				RideDifficulty:        workout.Ride.Difficulty_Rating,
				RideImageUrl:          workout.Ride.ImageURL,
				StartTime:             st,
				EndTime:               et,
				TimeZone:              z,
				StriveScore:           workout.EffortZones.TotalEffortPoints,
				HeartRateZone1Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone1Seconds,
				HeartRateZone1Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone1Seconds, totalHRZone),
				HeartRateZone2Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone2Seconds,
				HeartRateZone2Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone2Seconds, totalHRZone),
				HeartRateZone3Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone3Seconds,
				HeartRateZone3Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone3Seconds, totalHRZone),
				HeartRateZone4Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone4Seconds,
				HeartRateZone4Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone4Seconds, totalHRZone),
				HeartRateZone5Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone5Seconds,
				HeartRateZone5Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone5Seconds, totalHRZone),
			}

			// Store keys as minute, title for matching starting properly provided with offset.
			key = fmt.Sprintf("%s %s",
				st.Format(minuteLayout), strings.ToLower(workout.Ride.Title))
			apiWorkoutsMapOffset[key] = extraFields{
				Id:                    workout.Id,
				RideDifficulty:        workout.Ride.Difficulty_Rating,
				RideImageUrl:          workout.Ride.ImageURL,
				StartTime:             st,
				EndTime:               et,
				TimeZone:              z,
				StriveScore:           workout.EffortZones.TotalEffortPoints,
				HeartRateZone1Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone1Seconds,
				HeartRateZone1Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone1Seconds, totalHRZone),
				HeartRateZone2Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone2Seconds,
				HeartRateZone2Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone2Seconds, totalHRZone),
				HeartRateZone3Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone3Seconds,
				HeartRateZone3Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone3Seconds, totalHRZone),
				HeartRateZone4Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone4Seconds,
				HeartRateZone4Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone4Seconds, totalHRZone),
				HeartRateZone5Seconds: workout.EffortZones.HeartRateZoneDurations.HrZone5Seconds,
				HeartRateZone5Percent: divZero(workout.EffortZones.HeartRateZoneDurations.HrZone5Seconds, totalHRZone),
			}
		}
	}

	extractTime := time.Now().UTC()
	exportedWorkouts, err := peloservice.GetExportedWorkouts(*client.getHttpClient(), *client.getUserSession())
	if err != nil {
		return workouts, err
	}
	if len(exportedWorkouts) == 0 {
		return nil, nil
	}

	pAbbreviation := "\\d{4}-\\d{2}-\\d{2}\\s\\d{2}:\\d{2}\\s[(]\\w+[)]"
	rAbbreviation := regexp.MustCompile(pAbbreviation)

	pOffset := "\\d{4}-\\d{2}-\\d{2}\\s\\d{2}:\\d{2}\\s[(].?\\d+[)]"
	rOffset := regexp.MustCompile(pOffset)

	var extras extraFields
	var foundExtras bool
	for _, w := range exportedWorkouts {
		if w.FitnessDiscipline != "Cycling" {
			continue
		}
		startTime, _ := time.Parse(layout, w.StartTime)
		timeZone, _ := startTime.Zone()
		avgResistenceInt, _ := strconv.Atoi(strings.TrimRight(w.AvgResistance, "%s"))
		avgResistence := float64(avgResistenceInt) / 100.00

		if rAbbreviation.MatchString(w.StartTime) {
			key := fmt.Sprintf("%s %d %s",
				w.StartTime[:10], w.TotalOutput, strings.ToLower(w.ClassTitle))
			extras, foundExtras = apiWorkoutsMapAbbreviation[key]
		} else if rOffset.MatchString(w.StartTime) {
			key := fmt.Sprintf("%s %s",
				w.StartTime[:16], strings.ToLower(w.ClassTitle))
			extras, foundExtras = apiWorkoutsMapOffset[key]
		}

		if !foundExtras {
			log.Printf("warning, did not find api workouts for %s on at %s", w.ClassTitle, w.StartTime)
		} else {
			startTime = extras.StartTime
			timeZone = extras.TimeZone
		}

		// Distance based scenic rides report "None" for length minutes, 
		// and unfortunately "Scenic Ride" for name, so our current key
		// structure doesn't work.  Just set length to 0 for now.
		var lengthMinutes int
		if lengthMinutes, err = strconv.Atoi(w.LengthMinutes); err != nil {
			lengthMinutes = 0
		}

		workouts = append(workouts, Workout{
			Id:                    extras.Id,
			ExtractTimeUTC:        extractTime,
			StartTime:             startTime,
			TimeZone:              timeZone,
			StartTimeUTC:          startTime.UTC(),
			Type:                  w.ClassType,
			RideTitle:             w.ClassTitle,
			Instructor:            w.Instructor,
			RideLengthMinutes:     lengthMinutes,
			RideDifficulty:        extras.RideDifficulty,
			RideImageUrl:          extras.RideImageUrl,
			Output:                w.TotalOutput,
			AvgWatts:              w.AvgWatts,
			AvgResistance:         avgResistence,
			AvgCadenceRPM:         w.AvgCadenceRPM,
			AvgSpeedMPH:           w.AvgSpeedMPH,
			AvgSpeedKPH:           w.AvgSpeedKPH,
			DistanceMiles:         w.DistanceMiles,
			DistanceKilometers:    w.DistanceKilometers,
			CaloriesBurned:        w.CaloriesBurned,
			AvgHeartRate:          w.AvgHeartRate,
			StriveScore:           extras.StriveScore,
			HeartRateZone1Seconds: extras.HeartRateZone1Seconds,
			HeartRateZone1Percent: extras.HeartRateZone1Percent,
			HeartRateZone2Seconds: extras.HeartRateZone2Seconds,
			HeartRateZone2Percent: extras.HeartRateZone2Percent,
			HeartRateZone3Seconds: extras.HeartRateZone3Seconds,
			HeartRateZone3Percent: extras.HeartRateZone3Percent,
			HeartRateZone4Seconds: extras.HeartRateZone4Seconds,
			HeartRateZone4Percent: extras.HeartRateZone4Percent,
			HeartRateZone5Seconds: extras.HeartRateZone5Seconds,
			HeartRateZone5Percent: extras.HeartRateZone5Percent,
		})
	}

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

	if config.LogLevel == "DEBUG" {
		log.Printf("DEBUG: total workouts = %d, api workouts in map = %d", workouts.Len(), len(apiWorkoutsMapAbbreviation))
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

func divZero(a, b int) float64 {
	if b == 0 {
		return float64(b)
	}
	return float64(a) / float64(b)
}
