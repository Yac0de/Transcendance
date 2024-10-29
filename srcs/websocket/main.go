package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query().Get("id")
		_, err := strconv.ParseUint(values, 10, 64)
		if values == "" || err != nil {
			log.Printf("error decoding id query: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r.Header.Add("id", values)
		serveWs(hub, w, r)
	})
	log.Println("Server started on :4001")
	err := http.ListenAndServe(":4001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
