package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type Player struct {
	conn   *websocket.Conn
	id     string
	inRoom bool

	makeRoomReq chan<- RoomRequest
}

func newPlayer(conn *websocket.Conn, id string, makeRoomReq chan RoomRequest) *Player {
	return &Player{
		conn:        conn,
		id:          id,
		inRoom:      false,
		makeRoomReq: makeRoomReq,
	}
}

func (p *Player) sendString(msg string) {
	resp := struct {
		Status string `json:"status"`
	}{Status: msg}
	w, _ := p.conn.NextWriter(websocket.TextMessage)
	defer w.Close()
	json.NewEncoder(w).Encode(resp)
}

func (p *Player) receiveMsg() (map[string]interface{}, error) {
	msgType, r, err := p.conn.NextReader()
	if msgType != websocket.TextMessage || err != nil {
		p.sendString("error on get reader")
		return nil, errors.New("error on get reader")
	}

	data := map[string]interface{}{}
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		p.sendString("error on read data")
		return nil, errors.New("on read")
	}

	fmt.Printf("%+v\n", data)
	return data, nil

}

// TODO: refactor
func (p *Player) start() {
	for {
		data, err := p.receiveMsg()
		if err != nil {
			p.sendString("error")
			continue
		}
		fmt.Printf("%+v\n", data)
		if getEventType(data) == createRoom {
			p.id = data["playerId"].(string)
			if p.inRoom {
				p.sendString("already in room")
				continue
			}
			p.makeRoomReq <- RoomRequest{player: p, reqType: makeNewRoom, room: p.id}

		}

		if getEventType(data) == joinRoom {
			p.id = data["playerId"].(string)
			if p.inRoom {
				p.sendString("already in room")
				continue
			}
			p.makeRoomReq <- RoomRequest{player: p, reqType: playerJoinRoom, room: data["roomId"].(string)}
		}

	}
}
