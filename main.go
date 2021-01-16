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
	"github.com/weissleb/peloton-tableau-connector/service/servicehandlers"
	"os"
)

// user holds a users account information
type User struct {
	UserName   string
	FailedAuth bool
}

// tpl holds all parsed templates
var tpl *template.Template

var port string

func init() {

	gob.Register(User{})
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/home", homeHandler)
	r.HandleFunc("/", WdcHandler)
	r.Handle("/peloton-wdc", http.RedirectHandler("/", http.StatusFound))
	r.HandleFunc("/login", authHandler)
	r.HandleFunc("/cycling/schema/{table}", cyclingSchema)
	r.HandleFunc("/cycling/data/{table}", cyclingData)

	r.HandleFunc("/service/auth", servicehandlers.PostUserSession).
		Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/service/auth/check", servicehandlers.CheckAuth).
		Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/service/cycling/schema", servicehandlers.GetSimpleCyclingSchemas).
		Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/service/cycling/data/{table}", servicehandlers.GetCyclingDataJson).
		Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/service/cycling/summary", servicehandlers.GetCyclingDataSummary).
		Methods(http.MethodGet, http.MethodOptions)

	// start server
	fmt.Println(config.Banner)
	fmt.Printf("connector is on port %s\n", port)

	authMessage := "off"
	if config.RequireAuth {
		authMessage = "on"
	}
	log.Printf("authentication is %s", authMessage)

	cacheMessage := "off"
	if config.UseWorkoutCache {
		cacheMessage = "on"
	}
	log.Printf("caching of workouts is %s", cacheMessage)

	log.Fatal(http.ListenAndServe(":" + port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request)  {

	http.SetCookie(w, &http.Cookie{
		Name:  "peloton_wdc_host",
		Value: r.Host,
	})

	tpl.ExecuteTemplate(w, "home.gohtml", nil)
}

func WdcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	user := &User{
		UserName:   "unknown",
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

	proto := "http"
	if r.Host != "localhost:" + port {
		proto = "https"
	}
	requestUrl := fmt.Sprintf("%s://%s/service/auth", proto, r.Host)
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
		http.Redirect(w, r, "/?redirectCause=authFailed&user="+username, http.StatusFound)
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
		Name:    "peloton_wdc_token",
		Value:   userToken,
		Expires: expiration}

	userCookie := http.Cookie{
		Name:    "peloton_wdc_user",
		Value:   username,
		Expires: expiration,
	}
	http.SetCookie(w, &tokenCookie)
	http.SetCookie(w, &userCookie)

	http.Redirect(w, r, "/?user="+username, http.StatusFound)
}

func cyclingSchema(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	client := &http.Client{}
	vars := mux.Vars(r)
	table, _ := vars["table"]
	proto := "http"
	if r.Host != "localhost:" + port {
		proto = "https"
	}
	url := fmt.Sprintf("%s://%s/service/cycling/schema?tables=%s", proto, r.Host, table)
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
	proto := "http"
	if r.Host != "localhost:" + port {
		proto = "https"
	}
	url := fmt.Sprintf("%s://%s/service/cycling/data/%s", proto, r.Host, table)
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
