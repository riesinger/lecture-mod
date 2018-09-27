package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func startRouter() {
	m := mux.NewRouter()
	m.HandleFunc("/get/{tags}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET", r.URL)
		output, err := generateCalendar(mux.Vars(r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		w.Write(output)
	})
	m.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET", r.URL)
		m := make(map[string]string)
		output, err := generateCalendar(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		w.Write(output)
	})

	m.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("client/img"))))
	m.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("client/css"))))
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/index.html")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:10944", m))

}
