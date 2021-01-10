package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"math/rand"
	"fmt"
	"github.com/weissleb/peloton-tableau-connector/config"
	"io/ioutil"
	"github.com/gorilla/sessions"
	"os"
	"encoding/gob"
	"encoding/json"
	"strings"
	"time"
)

// user holds a users account information
type User struct {
	Username      string
	Authenticated bool
	UserToken     string
}

// store will hold all session data
var store *sessions.CookieStore

// tpl holds all parsed templates
var tpl *template.Template

func init() {
	authKeyOne := []byte(os.Getenv("FRONTEND_SESSION_KEY"))
	encryptionKeyOne := []byte(os.Getenv("FRONTEND_SESSION_KEY"))

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		//HttpOnly: true,
	}

	gob.Register(User{})

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

var layout = "2006-01-02 15:04:05"

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/peloton-wdc", WdcHandler)
	r.HandleFunc("/login", authHandler)
	r.HandleFunc("/token", tokenHandler).Methods(http.MethodPost)
	r.HandleFunc("/cycling/schema/{table}", cyclingSchema)
	r.HandleFunc("/cycling/data/{table}", cyclingData)

	// start server
	fmt.Println(config.Banner)
	fmt.Println("connector is at http://localhost:" + config.Port)
	http.ListenAndServe(":"+config.Port, r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	t := template.Must(template.New("example").ParseFiles("templates/example.html"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	rando := rand.Int()
	log.Print(rando)
	t.ExecuteTemplate(w, "example.html", struct {
		Rando int
	}{
		Rando: rando,
	})
}

func WdcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.ExecuteTemplate(w, "pelotonWDC.gohtml", nil)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	var username, password string
	if r.Method == http.MethodGet {
		username = os.Getenv("PELO_USER")
		password = os.Getenv("PELO_PASS")
	} else if r.Method == http.MethodPost {
		username = r.FormValue("username")
		password = r.FormValue("password")
	}

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
		//http.Redirect(w, r, "/forbidden", http.StatusFound)
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
	log.Printf("DEBUG UserToken from user_token = %s", userToken)

	// put the token into a cookie named "peloton_wdc_test"
	expiration := time.Now().Add(time.Minute) //365 * 24 * time.Hour
	cookie := http.Cookie{Name: "peloton_wdc_test", Value: userToken, Expires: expiration}
	log.Printf("DEGUB: cookie exipiration = %s", cookie.Expires.Format(layout))
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/peloton-wdc", http.StatusFound)
}

// TODO: I don't think I'm calling this.  Should be safe to remove it.
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request to %s %s", r.Method, r.URL.Path)

	// receive the cookie data in the POST request
	formData := &struct {
		Cookie string `json:"cookie"`
	}{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(body, formData)
	if err != nil {
		log.Fatal(err.Error())
	}

	// de-serialize it into a User object


	// return the UserToken
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
	log.Print("DEBUG authHeader = " + authHeader)
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