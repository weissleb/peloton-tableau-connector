package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"fmt"
	"github.com/weissleb/peloton-tableau-connector/config"
	"io/ioutil"
	"encoding/gob"
	"encoding/json"
	"strings"
	"time"
)

// user holds a users account information
type User struct {
	UserName   string
	FailedAuth bool
}

// tpl holds all parsed templates
var tpl *template.Template

func init() {

	gob.Register(User{})
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", WdcHandler)
	r.HandleFunc("/login", authHandler)
	r.HandleFunc("/cycling/schema/{table}", cyclingSchema)
	r.HandleFunc("/cycling/data/{table}", cyclingData)

	// start server
	fmt.Println(config.Banner)
	fmt.Printf("connector is at %s://%s:%s\n", config.Protocol, config.Host, config.Port)
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, r))
}

func WdcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	user := &User{
		UserName: "unknown",
		FailedAuth: false,
	}

	params := r.URL.Query()
	userName, ok := params["user"]
	if ok && len(userName[0]) > 0 {
		user.UserName = userName[0]
	}
	redirectCause, ok := params["redirectCause"]
	if ok && redirectCause[0] == "authFailed" {
		user.FailedAuth = true
	}

	tpl.ExecuteTemplate(w, "pelotonWDC.gohtml", user)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	var username, password string
	username = r.FormValue("username")
	password = r.FormValue("password")

	requestUrl := "http://localhost:30000/auth"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "username": "%s",
    "password": "%s"
}`, username, password))

	client := &http.Client{
	}
	req, err := http.NewRequest(method, requestUrl, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Print("error: could not authenticate")
		http.Redirect(w, r, "/?redirectCause=authFailed&user=" + username, http.StatusFound)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	respData := &struct {
		UserToken string `json:"user_token"`
	}{}

	err = json.Unmarshal(body, respData)
	if err != nil {
		log.Fatal(err.Error())
	}
	userToken := respData.UserToken
	log.Printf("got token with length %v", len(userToken))

	// put the token and username into cookies
	expiration := time.Now().Add(time.Hour)

	tokenCookie := http.Cookie{
		Name: "peloton_wdc_token",
		Value: userToken,
		Expires: expiration}

	userCookie := http.Cookie{
		Name:    "peloton_wdc_user",
		Value:   username,
		Expires: expiration,
	}
	http.SetCookie(w, &tokenCookie)
	http.SetCookie(w, &userCookie)

	http.Redirect(w, r, "/?user=" + username, http.StatusFound)
}

func cyclingSchema(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	client := &http.Client{}
	vars := mux.Vars(r)
	table, _ := vars["table"]
	url := "http://localhost:30000/cycling/schema?tables=" + table
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(body)
}

func cyclingData(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	client := &http.Client{}
	vars := mux.Vars(r)
	table, _ := vars["table"]
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		log.Print("error: did not find Authentication header")
	}
	if strings.Index(authHeader, "Bearer") != 0 {
		log.Print("error: the Authentication header is not a Bearer token")
	}
	url := "http://localhost:30000/cycling/data/" + table
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", authHeader)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(body)
}
