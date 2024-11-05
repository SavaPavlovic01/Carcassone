package main

import "fmt"

type GameRequestType int

const (
	startGame GameRequestType = iota
	playTile
	placeMeeple
	removeMeeple
	moveMeeple
	tempTile
)

type TileInfo struct {
	_tile *tile
	x     int
	y     int
}

type GameRequest struct {
	reqType       GameRequestType
	roomId        string
	player        *Player
	_tile         TileInfo
	validPosition bool
	meeple        Meeple
	index         int
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
				type response struct {
					MsgType   int          `json:"msgType"`
					MyTurn    bool         `json:"myTurn"`
					TileSides [4]tile_side `json:"tileSides"`
					Crest     bool         `json:"crest"`
					EndsRoad  bool         `json:"endsRoad"`
					Cathedral bool         `json:"cathedral"`
				}
				drawnTile := room._game.draw_card()
				resp := response{
					MsgType:   int(gameStarted),
					MyTurn:    false,
					TileSides: drawnTile.Sides,
					Crest:     drawnTile.Crest,
					EndsRoad:  drawnTile.Ends_road,
					Cathedral: drawnTile.Cathedral,
				}
				room.playerPingRoom(resp, room.playerIds[room.curTurn])
				resp.MyTurn = true
				room.players[room.playerIds[room.curTurn]].sendStruct(resp)
				//room.nextTurn()
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

				if room.playerIds[room.curTurn] != req.player.id {
					req.player.sendString("Not your turn")
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

				// saljemo svima koja je izvucena sledece posto je gotov potez
				type response struct {
					MsgType   int          `json:"msgType"`
					MyTurn    bool         `json:"myTurn"`
					TileSides [4]tile_side `json:"tileSides"`
					Crest     bool         `json:"crest"`
					EndsRoad  bool         `json:"endsRoad"`
					Cathedral bool         `json:"cathedral"`
				}
				drawnTile := room._game.draw_card()
				secondResp := response{
					MsgType:   int(gameStarted),
					MyTurn:    false,
					TileSides: drawnTile.Sides,
					Crest:     drawnTile.Crest,
					EndsRoad:  drawnTile.Ends_road,
					Cathedral: drawnTile.Cathedral,
				}
				room.nextTurn()
				room.playerPingRoom(secondResp, room.playerIds[room.curTurn])
				secondResp.MyTurn = true
				room.players[room.playerIds[room.curTurn]].sendStruct(secondResp)

			}

			if req.reqType == tempTile {
				room, err := gm.roomManager.getRoom(req.roomId)
				if err != nil {
					req.player.sendString("room doesnt exist")
					continue
				}
				if !room.gameStarted {
					req.player.sendString("game not started")
					continue
				}

				if room.playerIds[room.curTurn] != req.player.id {
					req.player.sendString("Not your turn")
					continue
				}

				resp := struct {
					MsgType   int          `json:"msgType"`
					TileX     int          `json:"tileX"`
					TileY     int          `json:"tileY"`
					TileSides [4]tile_side `json:"tileSides"`
					Crest     bool         `json:"crest"`
					EndsRoad  bool         `json:"endsRoad"`
					Cathedral bool         `json:"cathedral"`
					IsValid   bool         `json:"isValid"`
				}{
					MsgType:   int(tempTilePlaced),
					TileX:     req._tile.x,
					TileY:     req._tile.y,
					TileSides: req._tile._tile.Sides,
					Crest:     false,
					EndsRoad:  false,
					Cathedral: false,
					IsValid:   req.validPosition,
				}
				room.pingRoomStruct(resp)
			}

			if req.reqType == placeMeeple {
				room, err := gm.roomManager.getRoom(req.roomId)
				if err != nil {
					req.player.sendString("room doesnt exist")
					continue
				}
				if !room.gameStarted {
					req.player.sendString("game not started")
					continue
				}

				if room.playerIds[room.curTurn] != req.player.id {
					req.player.sendString("Not your turn")
					continue
				}

				room._game.addMeeple(req.meeple.x, req.meeple.y, req.meeple.color, req.meeple.isPriest)
				resp := struct {
					MsgType  int    `json:"msgType"`
					X        int    `json:"x"`
					Y        int    `json:"y"`
					Color    string `json:"color"`
					IsPriest bool   `json:"isPriest"`
				}{X: req.meeple.x, Y: req.meeple.y, Color: req.meeple.color,
					IsPriest: req.meeple.isPriest, MsgType: int(meepleAdded)}
				room.playerPingRoom(resp, req.player.id)
			}

			if req.reqType == removeMeeple {
				room, err := gm.roomManager.getRoom(req.roomId)
				if err != nil {
					req.player.sendString("room doesnt exist")
					continue
				}
				if !room.gameStarted {
					req.player.sendString("game not started")
					continue
				}

				if room.playerIds[room.curTurn] != req.player.id {
					req.player.sendString("Not your turn")
					continue
				}

				room._game.removeMeeple(req.meeple.x)
				resp := struct {
					MsgType int `json:"msgType"`
					Index   int `json:"index"`
				}{MsgType: int(meepleRemoved), Index: req.meeple.x}
				room.playerPingRoom(resp, req.player.id)
			}

			if req.reqType == moveMeeple {
				room, err := gm.roomManager.getRoom(req.roomId)
				if err != nil {
					req.player.sendString("room doesnt exist")
					continue
				}
				if !room.gameStarted {
					req.player.sendString("game not started")
					continue
				}

				if room.playerIds[room.curTurn] != req.player.id {
					req.player.sendString("Not your turn")
					continue
				}
				room._game.moveMeeple(req.index, req.meeple.x, req.meeple.y)
				resp := struct {
					MsgType int `json:"msgType"`
					Index   int `json:"index"`
					X       int `json:"x"`
					Y       int `json:"y"`
				}{MsgType: int(meepleMoved), Index: req.index, X: req.meeple.x, Y: req.meeple.y}
				room.playerPingRoom(resp, req.player.id)
			}

		}
	}
}
