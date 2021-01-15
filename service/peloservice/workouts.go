package peloservice

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/weissleb/peloton-tableau-connector/service/clients"
	"github.com/weissleb/peloton-tableau-connector/config"
)

func GetWorkouts(client clients.HttpClientInterface, session UserSession, page uint16) (workouts, error) {

	var err error
	workouts := workouts{}

	limit := config.PeloPageLimit

	path := fmt.Sprintf("/api/user/%s/workouts", session.UserId)
	requesturl := BaseUrl + path
	method := "GET"

	var (
		req  *http.Request
		res  *http.Response
		body []byte
	)

	req, err = http.NewRequest(method, requesturl, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Cookie", fmt.Sprintf(session.Cookies))

	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", page))
	q.Add("limit", fmt.Sprintf("%d", limit))
	q.Add("joins", "ride")
	req.URL.RawQuery = q.Encode()

	res, err = client.Do(req)
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return workouts, err
	}
	json.Unmarshal(body, &workouts)
	return workouts, nil
}

func GetRide(client clients.HttpClientInterface, session UserSession, id string) (ride, error) {

	var err error
	rideDetail := rideDetail{}

	path := fmt.Sprintf("/api/ride/%s/details", id)
	requesturl := BaseUrl + path
	method := "GET"

	var (
		req  *http.Request
		res  *http.Response
		body []byte
	)

	req, err = http.NewRequest(method, requesturl, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Cookie", fmt.Sprintf(session.Cookies))

	res, err = client.Do(req)
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return ride{}, err
	}
	json.Unmarshal(body, &rideDetail)
	return rideDetail.Ride, nil
}
