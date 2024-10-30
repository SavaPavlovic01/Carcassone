package main

import "fmt"

type GameRequestType int

const (
	startGame GameRequestType = iota
	playTile
)

type TileInfo struct {
	_tile *tile
	x     int
	y     int
}

type GameRequest struct {
	reqType GameRequestType
	roomId  string
	player  *Player
	_tile   TileInfo
}

type GameManager struct {
	roomManager *RoomManager
	gameReqs    chan GameRequest
}

func newGameManager(rm *RoomManager) GameManager {
	return GameManager{
		roomManager: rm,
		gameReqs:    make(chan GameRequest),
	}
}

func (gm GameManager) start() {
	for {
		select {
		case req := <-gm.gameReqs:
			if req.reqType == startGame {
				room, err := gm.roomManager.getRoom(req.roomId)
				if err != nil {
					req.player.sendString("room doesnt exist")
					continue
				}
				if room.gameStarted {
					req.player.sendString("Game already started")
					continue
				}
				if room.roomOwner != req.player.id {
					req.player.sendString("you are not the owner")
				}
				room.gameStarted = true
				room.pingRoom("game started")
			}

			if req.reqType == playTile {
				room, err := gm.roomManager.getRoom(req.roomId)
				if err != nil {
					req.player.sendString("room doesnt exist")
					continue
				}
				if !room.gameStarted {
					req.player.sendString("game not started")
					continue
				}

				req._tile._tile.print()
				err = room._game.place_tile(*req._tile._tile, req._tile.x, req._tile.y)
				if err != nil {
					req.player.sendString("invalid location")
				}
				fmt.Printf("%+v", room._game.Board)
				//room.pingRoom("new tile placed")
				resp := struct {
					MsgType   int          `json:"msgType"`
					TileX     int          `json:"tileX"`
					TileY     int          `json:"tileY"`
					TileSides [4]tile_side `json:"tileSides"`
					Crest     bool         `json:"crest"`
					EndsRoad  bool         `json:"endsRoad"`
					Cathedral bool         `json:"cathedral"`
				}{
					MsgType:   int(tileAdded),
					TileX:     req._tile.x,
					TileY:     req._tile.y,
					TileSides: req._tile._tile.Sides,
					Crest:     false,
					EndsRoad:  false,
					Cathedral: false,
				}
				room.pingRoomStruct(resp)
			}

		}
	}
}
