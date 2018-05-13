package main

import (
	"log"
	"net/http"
)

func startRouter() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ai", func(w http.ResponseWriter, r *http.Request) {
		mtx.RLock()
		w.Write([]byte(aiLectures))
		mtx.RUnlock()
	})
	mux.HandleFunc("/mi", func(w http.ResponseWriter, r *http.Request) {
		mtx.RLock()
		w.Write([]byte(miLectures))
		mtx.RUnlock()
	})


	log.Fatal(http.ListenAndServe("127.0.0.1:10944", mux))

}
