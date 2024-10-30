package main

import (
	"log"
	"net/http"
)

func main() {
	server := newServer()
	go server.rm.start()
	go server.gm.start()
	http.HandleFunc("/", server.serveWS)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
