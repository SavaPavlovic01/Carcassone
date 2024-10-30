package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	gm GameManager
	rm RoomManager
}

func newServer() Server {
	rm := RoomManager{
		rooms:       map[string]*Room{},
		makeRoomReq: make(chan RoomRequest),
	}
	gm := newGameManager(&rm)
	return Server{
		gm: gm,
		rm: rm,
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s Server) serveWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	player := newPlayer(conn, "", s.rm.makeRoomReq, s.gm.gameReqs)

	player.start()
}
