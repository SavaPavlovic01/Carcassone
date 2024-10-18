package main

import (
	"log"
	"net/http"
)

func main() {
	rm := RoomManager{}
	rm.rooms = map[string]Room{}
	rm.makeRoomReq = make(chan RoomRequest)
	go rm.start()
	http.HandleFunc("/", rm.serveWS)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
