package main

import (
	"log"
	"net/http"
)

func startRouter() {
	mux := http.NewServeMux()
	mux.HandleFunc("/rapla/ai", func(w http.ResponseWriter, r *http.Request) {
		mtx.RLock()
		w.Write([]byte(aiLectures))
		mtx.RUnlock()
	})
	mux.HandleFunc("/rapla/mi", func(w http.ResponseWriter, r *http.Request) {
		mtx.RLock()
		w.Write([]byte(miLectures))
		mtx.RUnlock()
	})
	mux.HandleFunc("/rapla", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("See <a href='/rapla/ai'>AI calendar</a> or <a href='/rapla/mi'>MI calendar</a>"))
	})


	log.Fatal(http.ListenAndServe("0.0.0.0:10944", mux))

}
