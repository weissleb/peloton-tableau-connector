package peloservice

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"github.com/weissleb/peloton-tableau-connector/service/clients"
)

func GetSession(client clients.HttpClientInterface, username string, password string) (UserSession, error) {

	var err error
	userSession := UserSession{}
	session := session{}

	path := fmt.Sprintf("/auth/login")
	requesturl := BaseUrl + path
	method := "POST"

	payload := strings.NewReader(
		fmt.Sprintf("{\"username_or_email\": \"%s\", \"password\": \"%s\"}", username, password))

	var (
		req  *http.Request
		res  *http.Response
		body []byte
	)

	req, err = http.NewRequest(method, requesturl, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/vnd.api+json")

	res, err = client.Do(req)
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if res.StatusCode == http.StatusOK {
		var cookies = make([]string, len(res.Cookies()))
		for i, cookie := range res.Cookies() {
			cookies[i] = cookie.Name + "=" + cookie.Value
		}
		userSession.Cookies = strings.Join(cookies, ";")

		if err != nil {
			return userSession, err
		}
		json.Unmarshal(body, &session)
		userSession.UserId = session.UserId
		return userSession, nil
	}

	return UserSession{}, errors.New(fmt.Sprintf("%d", res.StatusCode))


}
