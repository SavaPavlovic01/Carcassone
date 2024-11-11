package main

import (
	"log"
	"net/http"
)

func main() {
	server := newServer()
	go server.rm.start()
	go server.gm.start()

	static := http.Dir("./frontend/dist")

	http.HandleFunc("/ws", server.serveWS)
	http.Handle("/", http.FileServer(static))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
