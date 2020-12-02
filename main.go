package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"math/rand"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8889", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
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