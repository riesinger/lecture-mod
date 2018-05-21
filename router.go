package main

import (
	"log"
	"net/http"
)

func startRouter() {
	mux := http.NewServeMux()
	mux.HandleFunc("/rapla/get/", func(w http.ResponseWriter, r *http.Request) {
		mtx.RLock()
		w.Write([]byte(miLectures))
		mtx.RUnlock()
	})
	mux.Handle("/rapla/", http.StripPrefix("/rapla/", http.FileServer(http.Dir("client"))))


	log.Fatal(http.ListenAndServe("0.0.0.0:10944", mux))

}
