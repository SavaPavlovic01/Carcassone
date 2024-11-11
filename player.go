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
	gameReq     chan<- GameRequest
}

func newPlayer(conn *websocket.Conn, id string, makeRoomReq chan RoomRequest, gameReq chan GameRequest) *Player {
	return &Player{
		conn:        conn,
		id:          id,
		inRoom:      false,
		makeRoomReq: makeRoomReq,
		gameReq:     gameReq,
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

func (p *Player) sendStruct(msg interface{}) {

	w, _ := p.conn.NextWriter(websocket.TextMessage)
	defer w.Close()
	json.NewEncoder(w).Encode(msg)
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

	//fmt.Printf("%+v\n", data)
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
		//fmt.Printf("%+v\n", data)
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

		if getEventType(data) == gameStarted {
			p.gameReq <- GameRequest{player: p, reqType: startGame, roomId: data["roomId"].(string)}
		}

		if getEventType(data) == tileAdded {
			sides := [4]tile_side{} // TODO: TERRIBLE HACK PLS FIX
			for i, x := range data["sides"].([]interface{}) {
				sides[i] = tile_side(x.(float64))
			}
			addedTile := make_tile_from_array(sides, false, false, false)
			p.gameReq <- GameRequest{player: p,
				reqType: playTile,
				roomId:  data["roomId"].(string),
				_tile:   TileInfo{_tile: &addedTile, x: int(data["x"].(float64)), y: int(data["y"].(float64))}}
		}

		if getEventType(data) == tempTilePlaced {
			sides := [4]tile_side{} // TODO: TERRIBLE HACK PLS FIX
			for i, x := range data["sides"].([]interface{}) {
				sides[i] = tile_side(x.(float64))
			}
			addedTile := make_tile_from_array(sides, false, false, false)
			p.gameReq <- GameRequest{player: p,
				reqType:       tempTile,
				roomId:        data["roomId"].(string),
				_tile:         TileInfo{_tile: &addedTile, x: int(data["x"].(float64)), y: int(data["y"].(float64))},
				validPosition: data["isValid"].(bool),
			}
		}

		if getEventType(data) == meepleAdded {
			//fmt.Println(data["color"])
			newMeeple := Meeple{
				x:        int(data["x"].(float64)),
				y:        int(data["y"].(float64)),
				color:    fmt.Sprint(data["color"]),
				isPriest: data["isPriest"].(bool),
			}
			p.gameReq <- GameRequest{player: p, reqType: placeMeeple, roomId: data["roomId"].(string), meeple: newMeeple}
		}

		if getEventType(data) == meepleRemoved {
			p.gameReq <- GameRequest{player: p, reqType: removeMeeple,
				roomId: data["roomId"].(string), meeple: Meeple{x: int(data["index"].(float64))}}
		}

		if getEventType(data) == meepleMoved {
			meeple := Meeple{
				x: int(data["x"].(float64)),
				y: int(data["y"].(float64)),
			}
			p.gameReq <- GameRequest{player: p, reqType: moveMeeple, roomId: data["roomId"].(string),
				index: int(data["index"].(float64)), meeple: meeple}
		}

	}
}
