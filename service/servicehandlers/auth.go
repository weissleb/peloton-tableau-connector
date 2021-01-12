package servicehandlers

import (
	"net/http"
	"log"
	"encoding/json"
	"encoding/base64"
	"io/ioutil"
	"strconv"
	"fmt"
	"github.com/weissleb/peloton-tableau-connector/service/extractors"
	"github.com/weissleb/peloton-tableau-connector/config"
)

func PostUserSession(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)
	formData := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(body, formData)
	if err != nil {
		log.Fatal(err.Error())
	}

	// os.Getenv("PELO_USER"), os.Getenv("PELO_PASS")
	pelotonClient, err := extractors.NewPelotonClient(httpClient, formData.Username, formData.Password)
	if err != nil {
		status, err := strconv.Atoi(err.Error())
		if err != nil {
			status = http.StatusInternalServerError
		}
		log.Printf("could not authenticate, got %v status from Peloton", status)
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response, err := json.Marshal(struct {
		UserToken string `json:"user_token"`
	}{
		UserToken: base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", pelotonClient.GetSessionUser(), pelotonClient.GetSessionCookie()))),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(response)
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)
	check := &struct {
		RestRequiresAuth bool `json:"rest_requires_auth"`
	}{
		RestRequiresAuth: true,
	}

	if !config.RequireAuth {
		check.RestRequiresAuth = false
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response, err := json.Marshal(check)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(response)
}
