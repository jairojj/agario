package main

import (
	"log"
	"net/http"
)

func main() {
	hub := newHub()
	go hub.run()

	log.Println("Server started on localhost:3000")

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
