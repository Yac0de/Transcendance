package main

import (
	"log"
	"net/http"
	"strconv"
	"websocket/controllers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(hub *controllers.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	id, err := strconv.ParseUint(r.Header.Get("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		return
	}

	client := &controllers.Client{
		Id:   id,
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 1024),
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func main() {
	hub := controllers.NewHub()
	go hub.Run()

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
