package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func startRouter() {
	m := mux.NewRouter()
	m.HandleFunc("/rapla/get/{tags}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET", r.URL)
		output, err := generateCalendar(mux.Vars(r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		w.Write(output)
	})
	m.HandleFunc("/rapla/get/", func(w http.ResponseWriter, r *http.Request) {
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
	m.Handle("/rapla/", http.StripPrefix("/rapla/", http.FileServer(http.Dir("client"))))

	log.Fatal(http.ListenAndServe("0.0.0.0:10944", m))

}
