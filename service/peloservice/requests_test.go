package peloservice

import (
	"log"
	"testing"
	"net/http"
	"io/ioutil"
	"bytes"
	"reflect"
	"github.com/weissleb/peloton-tableau-connector/service/clients"
)

var httpClient clients.HttpClientInterface

func init() {
	log.Printf("initializing mock client", )
	httpClient = &clients.MockClient{}
}

// These tests are crude and poorly constructed.
// Please, please make them better before actually shipping!
func TestGetSession(t *testing.T) {

	clients.GetDoFunc = func(req *http.Request) (response *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewReader([]byte(`
{
    "session_id": "not real",
    "user_id": "not real",
    "user_data": {
        "location": "Noblesville, IN",
        "total_pending_followers": 0,
        "created_at": 1518835623,
        "customized_max_heart_rate": 176,
        "cycling_ftp_source": "ftp_manual_source",
        "subscription_credits_used": 0,
        "name": "Brian Weissler",
        "default_heart_rate_zones": [
            0,
            114.4,
            132.0,
            149.6,
            167.2
        ],
        "contract_agreements": [
            {
                "contract_type": "privacy_policy",
                "contract_id": "",
                "contract_created_at": 1600101574,
                "bike_contract_url": "https://s3.us-east-2.amazonaws.com/contract-terms-html/prod/en-us/privacy_1600101574.html",
                "tread_contract_url": "https://s3.us-east-2.amazonaws.com/contract-terms-html/prod/en-us/privacy_1600101574.html",
                "agreed_at": 1577710450,
                "contract_display_name": "Privacy Policy"
            },
            {
                "contract_type": "subscription_terms",
                "contract_id": "",
                "contract_created_at": 1598562045,
                "bike_contract_url": "https://s3.us-east-2.amazonaws.com/contract-terms-html/prod/en-us/subscription_tos_1598562045.html",
                "tread_contract_url": "https://s3.us-east-2.amazonaws.com/contract-terms-html/prod/en-us/subscription_tos_1598562045.html",
                "agreed_at": 1539004459,
                "contract_display_name": "Membership Terms"
            },
            {
                "contract_type": "terms_of_service",
                "contract_id": "",
                "contract_created_at": 1600103588,
                "bike_contract_url": "https://s3.us-east-2.amazonaws.com/contract-terms-html/prod/en-us/tos_1600103588.html",
                "tread_contract_url": "https://s3.us-east-2.amazonaws.com/contract-terms-html/prod/en-us/tos_1600103588.html",
                "agreed_at": 1577710450,
                "contract_display_name": "Terms of Service"
            }
        ],
        "phone_number": "",
        "is_fitbit_authenticated": true,
        "facebook_access_token": "",
        "cycling_ftp_workout_id": null,
        "email": "fake",
        "quick_hits": {
            "quick_hits_enabled": true,
            "speed_shortcuts": null,
            "incline_shortcuts": null
        },
        "id": "not real",
        "is_demo": false,
        "last_name": "Weissler",
        "total_pedaling_metric_workouts": 409,
        "birthday": 0,
        "v1_referrals_made": 0,
        "username": "",
        "obfuscated_email": "dfa9b31328c193f271e596b5e6ac51704168b99345a8c7c8ef78948c42064798",
        "first_name": "Brian",
        "member_groups": [],
        "referral_code": "6212dab3e0bc45ae8b5305e9c2786f27",
        "total_non_pedaling_metric_workouts": 11,
        "estimated_cycling_ftp": 157,
        "cycling_workout_ftp": 0,
        "workout_counts": [
            {
                "name": "Yoga",
                "slug": "yoga",
                "count": 0,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/zero-yoga2.png"
            },
            {
                "name": "Stretching",
                "slug": "stretching",
                "count": 3,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/nonzero-stretching2.png"
            },
            {
                "name": "Strength",
                "slug": "strength",
                "count": 1,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/nonzero-strength2.png"
            },
            {
                "name": "Bootcamp",
                "slug": "circuit",
                "count": 0,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/zero-circuit2.png"
            },
            {
                "name": "Running",
                "slug": "running",
                "count": 0,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/zero-running2.png"
            },
            {
                "name": "Cycling",
                "slug": "cycling",
                "count": 414,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/nonzero-cycling2.png"
            },
            {
                "name": "Walking",
                "slug": "walking",
                "count": 0,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/zero-walking2.png"
            },
            {
                "name": "Cardio",
                "slug": "cardio",
                "count": 0,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/zero-cardio2.png"
            },
            {
                "name": "Meditation",
                "slug": "meditation",
                "count": 2,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/workout-count-icons/nonzero-meditation2.png"
            }
        ],
        "gender": "male",
        "has_signed_waiver": false,
        "is_internal_beta_tester": false,
        "weight": 170.0,
        "height": 70.0,
        "paired_devices": [
            {
                "name": "TREKZ Titanium by AfterShokz",
                "paired_device_type": "audio",
                "serial_number": "20:74:CF:28:EA:1B"
            },
            {
                "name": "Device: 65398",
                "paired_device_type": "heart_rate_monitor",
                "serial_number": "65398"
            }
        ],
        "instructor_id": null,
        "customized_heart_rate_zones": [],
        "can_charge": true,
        "created_country": null,
        "total_followers": 40,
        "default_max_heart_rate": 176,
        "is_complete_profile": true,
        "total_following": 13,
        "image_url": "https://s3.amazonaws.com/peloton-profile-images/8d06aabc8986c967c4bf2abe790a366d34566714/2eea52a163b444bf80ec9729bd05df1d",
        "is_profile_private": false,
        "is_provisional": false,
        "facebook_id": "",
        "block_explicit": false,
        "referrals_made": 0,
        "is_external_beta_tester": false,
        "last_workout_at": 1600435742,
        "is_strava_authenticated": false,
        "has_active_device_subscription": true,
        "has_active_digital_subscription": false,
        "subscription_credits": 0,
        "total_workouts": 420,
        "cycling_ftp": 142,
        "middle_initial": ""
    },
    "pubsub_session": {}
}`))),
		}, nil
	}

	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    UserSession
		wantErr bool
	}{
		{
			name: "Test getting a session.",
			args: args{
				username: "username",
				password: "password",
			},
			want: UserSession{
				UserId: "not real",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSession(httpClient, tt.args.username, tt.args.password)
			if err != nil != tt.wantErr {
				t.Errorf("GetSession error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.UserId != tt.want.UserId {
				t.Errorf("UserId -> got: %s; want %s", got.UserId, tt.want.UserId)
			}
		})
	}
}

func TestGetWorkouts(t *testing.T) {

	clients.GetDoFunc = func(req *http.Request) (response *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewReader([]byte(`
{
    "data": [
		{
            "created_at": 1600435742,
            "device_type": "home_bike_v1",
            "end_time": 1600438549,
            "fitbit_id": "59007145",
            "fitness_discipline": "cycling",
            "has_pedaling_metrics": true,
            "has_leaderboard_metrics": true,
            "id": "cfc1ea921353437caf8630825dd6d8b6",
            "is_total_work_personal_record": true,
            "metrics_type": "cycling",
            "name": "Cycling Workouts",
            "peloton_id": "22f01cb3f09f47e8a523cffbc9367c65",
            "platform": "home_bike",
            "start_time": 1600435852,
            "strava_id": null,
            "status": "COMPLETE",
            "timezone": "America/New_York",
            "title": null,
            "total_work": 482432.95,
            "user_id": "not real",
            "workout_type": "class",
            "total_video_watch_time_seconds": 2791,
            "total_video_buffering_seconds": 0,
            "v2_total_video_watch_time_seconds": null,
            "v2_total_video_buffering_seconds": null,
            "created": 1600435742,
            "device_time_created_at": 1600421342,
            "effort_zones": null
        },
        {
            "created_at": 1600347510,
            "device_type": "home_bike_v1",
            "end_time": 1600349370,
            "fitbit_id": "55654796",
            "fitness_discipline": "cycling",
            "has_pedaling_metrics": true,
            "has_leaderboard_metrics": true,
            "id": "4f00f7c891e24584ba3505ff946556b9",
            "is_total_work_personal_record": false,
            "metrics_type": "cycling",
            "name": "Cycling Workouts",
            "peloton_id": "9c2198fe20ac400aa6cd9fc9536cbdca",
            "platform": "home_bike",
            "start_time": 1600347572,
            "strava_id": null,
            "status": "COMPLETE",
            "timezone": "America/New_York",
            "title": null,
            "total_work": 289132.36,
            "user_id": "not real",
            "workout_type": "class",
            "total_video_watch_time_seconds": 1831,
            "total_video_buffering_seconds": 0,
            "v2_total_video_watch_time_seconds": null,
            "v2_total_video_buffering_seconds": null,
            "created": 1600347510,
            "device_time_created_at": 1600333110,
            "effort_zones": null
        }
    ],
    "limit": 25,
    "page": 15,
    "total": 420,
    "count": 25,
    "page_count": 17,
    "show_previous": true,
    "show_next": true,
    "sort_by": "-created_at,-pk"
}`))),
		}, nil
	}

	type args struct {
		session UserSession
		page    uint16
	}
	tests := []struct {
		name    string
		args    args
		want    workouts
		wantErr bool
	}{
		{
			name: "Test getting workouts.",
			args: args{
				session: UserSession{
					Cookies: "not real",
					UserId:  "not real",
				},
				page: 15,
			},
			want: workouts{
				Workouts: []workout{
					{
						Id:               "cfc1ea921353437caf8630825dd6d8b6",
						CurrentPr:        true,
						StartTimeSeconds: 1600435852,
						Timezone:         "America/New_York",
						Wtype:            "cycling",
						Status:           "COMPLETE",
						Output:           482432.95,
					},
					{
						Id:               "4f00f7c891e24584ba3505ff946556b9",
						CurrentPr:        false,
						StartTimeSeconds: 1600347572,
						Timezone:         "America/New_York",
						Wtype:            "cycling",
						Status:           "COMPLETE",
						Output:           289132.36,
					},
				},
				PageCount: 17,
				HasNext:   true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWorkouts(httpClient, tt.args.session, tt.args.page)
			if err != nil != tt.wantErr {
				t.Errorf("GetWorkouts error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Workouts) != len(tt.want.Workouts) {
				t.Errorf("number of workouts -> got: %d; want %d", len(got.Workouts), len(tt.want.Workouts))
			}
			for i, workout := range got.Workouts {
				if !reflect.DeepEqual(workout, tt.want.Workouts[i]) {
					t.Errorf("workout %d -> got: %v; want %v", i, workout, tt.want.Workouts[i])
				}
			}
		})
	}
}

func TestGetRide(t *testing.T) {

	clients.GetDoFunc = func(req *http.Request) (response *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewReader([]byte(`
{
    "ride": {
        "class_type_ids": [
            "7579b9edbdf9464fa19eb58193897a73"
        ],
        "content_provider": "peloton",
        "content_format": "video",
        "description": "Join Kendall in this special Intervals & Arms ride featuring the dancehall musical stylings of Beenie Man and Bounty Killer. This upbeat, Jamaican-themed class will challenge you to push through alternating bursts of efforts and free weight segments.",
        "difficulty_estimate": 8.653,
        "overall_estimate": 0.9603,
        "difficulty_rating_avg": 8.653,
        "difficulty_rating_count": 14273,
        "difficulty_level": null,
        "duration": 2700,
        "equipment_ids": [
            "0f5f1ff2d6c647cf98d599ed90ad72d3"
        ],
        "equipment_tags": [
            {
                "id": "0f5f1ff2d6c647cf98d599ed90ad72d3",
                "name": "Light Weights",
                "slug": "light_weights",
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/equipment-icons/light_weights.png"
            }
        ],
        "extra_images": [],
        "fitness_discipline": "cycling",
        "fitness_discipline_display_name": "Cycling",
        "has_closed_captions": true,
        "has_pedaling_metrics": true,
        "home_peloton_id": "a89c3a607d5b4eee83ad83ee9907461d",
        "id": "9d03c9c251d74081b0f9c811ea31138d",
        "image_url": "https://s3.amazonaws.com/peloton-ride-images/02ccb5d3c70c82936fed06f9636cee838c207044/img_1599496590_6d90d0e9f45e4f01b64922b8205b6950.png",
        "instructor_id": "4904612965164231a37143805a387e40",
        "is_archived": true,
        "is_closed_caption_shown": true,
        "is_explicit": true,
        "has_free_mode": false,
        "is_live_in_studio_only": false,
        "language": "english",
        "origin_locale": "en-US",
        "length": 2829,
        "live_stream_id": "9d03c9c251d74081b0f9c811ea31138d-live",
        "live_stream_url": null,
        "location": "psny-studio-1",
        "metrics": [
            "heart_rate",
            "cadence",
            "calories"
        ],
        "original_air_time": 1599492120,
        "overall_rating_avg": 0.9603,
        "overall_rating_count": 16516,
        "pedaling_start_offset": 60,
        "pedaling_end_offset": 2760,
        "pedaling_duration": 2700,
        "rating": 0,
        "ride_type_id": "7579b9edbdf9464fa19eb58193897a73",
        "ride_type_ids": [
            "7579b9edbdf9464fa19eb58193897a73"
        ],
        "sample_vod_stream_url": null,
        "scheduled_start_time": 1599492600,
        "series_id": "283319daf8834b86a6205737001b0d56",
        "sold_out": false,
        "studio_peloton_id": "6f6279c75ea34802864ad0e1f67af6ba",
        "title": "45 min Intervals & Arms Ride: Dancehall",
        "total_ratings": 0,
        "total_in_progress_workouts": 45,
        "total_workouts": 33824,
        "vod_stream_url": "https://amd-vod.akamaized.net/vod/bike/09-2020/09072020-kendall_toole-1130am-drastic-1-9d03c9c251d74081b0f9c811ea31138d/HLS/master.m3u8",
        "vod_stream_id": "9d03c9c251d74081b0f9c811ea31138d-vod",
        "captions": [
            "en-US"
        ],
        "join_tokens": {
            "on_demand": ""
        },
        "instructor": {
            "id": "4904612965164231a37143805a387e40",
            "bio": "Kendall is a natural-born fighter who empowers you to find your voice and believe in your inner strength. Kendall is a dynamic athlete, with a background in professional cheerleading, gymnastics, dance and boxing. After college, Kendall began what she believed would be a dream job at a tech startup, but soon realized the culture wasn’t the right fit for her. To cope, Kendall began boxing, where she discovered that her true calling is empowering others to break down barriers through movement and sweat. In her classes, Kendall uses the power of positive energy and rock 'n' roll to inspire you to rise up and discover your highest self.",
            "short_bio": "Kendall is a natural-born fighter who empowers you to believe in your inner strength. As a dynamic athlete, coach and creative, Kendall uses the power of positive energy to inspire you to rise up and commit to your highest self. ",
            "coach_type": "peloton_coach",
            "is_filterable": true,
            "is_instructor_group": false,
            "is_visible": true,
            "list_order": 27,
            "featured_profile": true,
            "film_link": "",
            "facebook_fan_page": "",
            "music_bio": "",
            "spotify_playlist_uri": "",
            "background": "My goal is to help you discover your power. We have a choice when it comes to listening to the voice in our heads: it holds us back or helps us rise up. It just depends on how we train it.",
            "ordered_q_and_as": [
                [
                    "How Do You Motivate?",
                    "POSITIVE ENERGY! What we put out in the world comes back to us, so there is no time for negativity. We are all in this together for support and accountability, and I hope you enter my class feeling curious and proud."
                ],
                [
                    "Outside of Peloton",
                    "I am a boxer at heart and have been learning the art for over 7 years. Rock ‘n’ roll seeps through my soul, especially the classics, so you can also find me listening to vinyl or scrounging up tickets to a show. "
                ],
                [
                    "",
                    ""
                ]
            ],
            "instagram_profile": "",
            "strava_profile": "",
            "twitter_profile": "",
            "quote": "Rise up and greet your higher self",
            "username": "kmoneyNYC",
            "name": "Kendall Toole",
            "first_name": "Kendall",
            "last_name": "Toole",
            "user_id": "d7833badf6694abfafd92fe1dd51c12e",
            "life_style_image_url": "https://s3.amazonaws.com/workout-metric-images-prod/5d34ae7f03794b30b7e6b68374c8fbd4",
            "bike_instructor_list_display_image_url": null,
            "web_instructor_list_display_image_url": "https://s3.amazonaws.com/workout-metric-images-prod/6312d386e064462e97ed8775e78bc063",
            "ios_instructor_list_display_image_url": "https://s3.amazonaws.com/workout-metric-images-prod/dd19df34c0d04cf685e44db2ab67e248",
            "about_image_url": "https://s3.amazonaws.com/workout-metric-images-prod/fb0c6a5efaa142aba993047915423174",
            "image_url": "https://s3.amazonaws.com/workout-metric-images-prod/63cd471dddd6423294737112b32a8b78",
            "jumbotron_url": null,
            "jumbotron_url_dark": "https://s3.amazonaws.com/workout-metric-images-prod/ac7613daf3ad4c148bc9c0ebf111266a",
            "jumbotron_url_ios": "https://s3.amazonaws.com/workout-metric-images-prod/e3cf2cd1902e411d98983686af639483",
            "web_instructor_list_gif_image_url": null,
            "instructor_hero_image_url": "https://s3.amazonaws.com/workout-metric-images-prod/bb032c53ecc946ada7188e79db031d32",
            "fitness_disciplines": [
                "strength",
                "cardio",
                "cycling",
                "stretching"
            ]
        },
        "is_favorite": false,
        "total_user_workouts": 1,
        "total_following_workouts": 0,
        "leaderboard_filter_type": null
    },
    "playlist": {
        "id": "91a10e0bc6fb470e891c15dc44e02ba5",
        "ride_id": "9d03c9c251d74081b0f9c811ea31138d",
        "top_artists": [
            {
                "artist_id": "fc268e012a2f429d97ff7fceb045f9d1",
                "artist_name": "Beenie Man"
            },
            {
                "artist_id": "6fbc54ae39b343e48b01d5845014a7cd",
                "artist_name": "Mobb Deep"
            },
            {
                "artist_id": "a11b4f6a7a4d4c94b9f9694b69de2ee0",
                "artist_name": "Damian Marley"
            },
            {
                "artist_id": "b3c2a50ad57542c38ece4f968f24b4e8",
                "artist_name": "Busy Signal"
            }
        ],
        "songs": [
            {
                "id": "6ff1045b7d1d4421a650ce12c5776e82",
                "title": "Lodge",
                "artists": [
                    {
                        "artist_id": "70b1d9dd6aee4706ba7205d6b1ec6ad4",
                        "artist_name": "Bounty Killer"
                    }
                ],
                "album": {
                    "id": "ffaa99d56b4243f2af3d0fc0a70b5b44",
                    "image_url": "https://neurotic.azureedge.net/RR/AlbumImages/Catalog/1a3908d2-7295-4613-9092-679b5162962c/Product/5ba310de-a5de-4aa7-a1f5-bce6d4f0c477/big_054645134163.jpg",
                    "name": "Roots, Reality & Culture"
                },
                "cue_time_offset": 60,
                "start_time_offset": 60,
                "liked": false
            },
            {
                "id": "8fe3abcedcd84070b14d7f0242f36de5",
                "title": "Who Am I",
                "artists": [
                    {
                        "artist_id": "fc268e012a2f429d97ff7fceb045f9d1",
                        "artist_name": "Beenie Man"
                    }
                ],
                "album": {
                    "id": "fd580e7e9f884682944f5974de05aae2",
                    "image_url": "https://neurotic.azureedge.net/RR/AlbumImages/Catalog/1a3908d2-7295-4613-9092-679b5162962c/Product/c159e686-acec-4d52-8fbc-9002be35aef6/big_054645160568.jpg",
                    "name": "Best Of (collector's Edition)"
                },
                "cue_time_offset": 282,
                "start_time_offset": 282,
                "liked": false
            }
        ],
        "is_top_artists_shown": true,
        "is_playlist_shown": true,
        "is_in_class_music_shown": true
    },
    "averages": {
        "average_total_work": 274,
        "average_distance": 10.52,
        "average_calories": 388,
        "average_avg_power": 107,
        "average_avg_speed": 14.9,
        "average_avg_cadence": 63,
        "average_avg_resistance": 45
    },
    "segments": {
        "segment_list": [
            {
                "id": "7f7acfb23faa425b9af7fe45e4f98642",
                "length": 419,
                "start_time_offset": 0,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/segment-icons/warmup.png",
                "intensity_in_mets": 3.5,
                "metrics_type": "cycling",
                "icon_name": "warmup",
                "icon_slug": "warmup",
                "name": "Warmup"
            },
            {
                "id": "5624634987554da8b04953ce7b4ce20f",
                "length": 446,
                "start_time_offset": 419,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/segment-icons/cycling.png",
                "intensity_in_mets": 6.0,
                "metrics_type": "cycling",
                "icon_name": "cycling",
                "icon_slug": "cycling",
                "name": "Cycling"
            },
            {
                "id": "5af47f38910f4c86b4ee4372914c7098",
                "length": 416,
                "start_time_offset": 865,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/segment-icons/upper_body.png",
                "intensity_in_mets": 6.5,
                "metrics_type": "cycling",
                "icon_name": "upper_body",
                "icon_slug": "upper_body",
                "name": "Arms"
            },
            {
                "id": "58afeeb879784b0b82aca7476b6d3bc4",
                "length": 506,
                "start_time_offset": 1281,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/segment-icons/cycling.png",
                "intensity_in_mets": 6.0,
                "metrics_type": "cycling",
                "icon_name": "cycling",
                "icon_slug": "cycling",
                "name": "Cycling"
            },
            {
                "id": "81c316abeb5a4bdb94afdbe1e91bfdcb",
                "length": 437,
                "start_time_offset": 1787,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/segment-icons/upper_body.png",
                "intensity_in_mets": 6.5,
                "metrics_type": "cycling",
                "icon_name": "upper_body",
                "icon_slug": "upper_body",
                "name": "Arms"
            },
            {
                "id": "2b24e5b807bb4e5bb2d783616cd639b7",
                "length": 416,
                "start_time_offset": 2224,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/segment-icons/cycling.png",
                "intensity_in_mets": 6.0,
                "metrics_type": "cycling",
                "icon_name": "cycling",
                "icon_slug": "cycling",
                "name": "Cycling"
            },
            {
                "id": "ea8eec45ac4648308c1cb15661e3a783",
                "length": 60,
                "start_time_offset": 2640,
                "icon_url": "https://s3.amazonaws.com/static-cdn.pelotoncycle.com/segment-icons/cooldown.png",
                "intensity_in_mets": 3.5,
                "metrics_type": "cycling",
                "icon_name": "cooldown",
                "icon_slug": "cooldown",
                "name": "Cool Down"
            }
        ],
        "segment_category_distribution": {
            "Cycling Warmup": "0.15518518518518518",
            "cycling": "0.5066666666666667",
            "Cycling_Arms": "0.31592592592592594",
            "Cycling Cool Down": "0.022222222222222223"
        },
        "segment_body_focus_distribution": {
            "cardio": "1.0",
            "arms": "0.31592592592592594"
        }
    },
    "default_album_images": {
        "default_in_class_image_url": "https://s3.amazonaws.com/peloton-ride-images/DEFAULT_ALBUM_ART_IN_CLASS.svg",
        "default_class_detail_image_url": "https://s3.amazonaws.com/peloton-ride-images/DEFAULT_ALBUM_ART_CLASS_DETAIL.svg"
    },
    "excluded_platforms": [],
    "is_ftp_test": false,
    "disabled_leaderboard_filters": {
        "just_me": false,
        "following": false,
        "age_and_gender": false
    },
    "sampled_top_tags": null,
    "instructor_cues": [
        {
            "offsets": {
                "start": 60,
                "end": 108
            },
            "resistance_range": {
                "upper": 35,
                "lower": 25
            },
            "cadence_range": {
                "upper": 100,
                "lower": 80
            }
        },
        {
            "offsets": {
                "start": 109,
                "end": 134
            },
            "resistance_range": {
                "upper": 38,
                "lower": 27
            },
            "cadence_range": {
                "upper": 100,
                "lower": 80
            }
        }
    ],
    "target_class_metrics": {
        "target_graph_metrics": [
            {
                "graph_data": {
                    "upper": [
                        100,
                        100,
                        100
                    ],
                    "lower": [
                        80,
                        80,
                        80
                    ],
                    "average": [
                        90,
                        90,
                        90
                    ]
                },
                "max": 120,
                "min": 50,
                "average": 80,
                "type": "cadence"
            },
            {
                "graph_data": {
                    "upper": [
                        35,
                        35,
                        35
                    ],
                    "lower": [
                        25,
                        25,
                        25
                    ],
                    "average": [
                        30,
                        30,
                        30
                    ]
                },
                "max": 65,
                "min": 25,
                "average": 46,
                "type": "resistance"
            }
        ],
        "total_expected_output": {
            "expected_upper_output": 673,
            "expected_lower_output": 267
        }
    },
    "events": {
        "data": []
    }
}`))),
		}, nil
	}

	type args struct {
		session UserSession
		id      string
	}
	tests := []struct {
		name    string
		args    args
		want    ride
		wantErr bool
	}{
		{
			name: "Test getting ride.",
			args: args{
				session: UserSession{
					Cookies: "not real",
					UserId:  "not real",
				},
				id: "4af299ce8a794ae29996b762f353fab4",
			},
			want: ride{
				Id:                "9d03c9c251d74081b0f9c811ea31138d",
				Title:             "45 min Intervals & Arms Ride: Dancehall",
				Duration:          2700,
				Difficulty_Rating: 8.653,
				Instructor: struct {
					Name     string `json:"name"`
					ImageURL string `json:"about_image_url"`
				}{
					Name:     "Kendall Toole",
					ImageURL: "https://s3.amazonaws.com/workout-metric-images-prod/fb0c6a5efaa142aba993047915423174",
				},
				Equipmenttags: []struct{ Slug string `json:"slug"` }{
					{
						Slug: "light_weights",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRide(httpClient, tt.args.session, tt.args.id)
			if err != nil != tt.wantErr {
				t.Errorf("GetRide error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ride -> got: %v; want %v", got, tt.want)
			}
		})
	}
}
