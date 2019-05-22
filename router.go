package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func startRouter() {
	gmt, err := time.LoadLocation("Iceland")
	if err != nil {
		log.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/get/{tags}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET", r.URL)
		output, err := generateCalendar(mux.Vars(r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		// Reference time is Mon Jan 2 15:04:05 MST 2006
		// Expcted format is Wed, 22 May 2019 08:11:07 GMT
		w.Header().Set("Last-Modified", time.Now().In(gmt).Add(-1*time.Second).Format("Mon, 2 Jan 2006 15:04:05 GMT"))
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
		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		// Reference time is Mon Jan 2 15:04:05 MST 2006
		// Expcted format is Wed, 22 May 2019 08:11:07 GMT
		w.Header().Set("Last-Modified", time.Now().In(gmt).Add(-1*time.Second).Format("Mon, 2 Jan 2006 15:04:05 GMT"))

		w.Write(output)
	})

	m.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("client/img"))))
	m.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("client/css"))))
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/index.html")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:10944", m))

}
