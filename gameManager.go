package main

import (
	"errors"
)

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
	RoomManager *RoomManager
	gameReqs    chan GameRequest
}

type Validator func(GameManager, GameRequest, *Room) error

type EventHandler func(*Room, GameRequest) error

func roomExists(gm GameManager, req GameRequest, room *Room) error {
	if room == nil {
		req.player.sendString("room doesnt exist")
		return errors.New("room doesnt exist")
	}
	return nil
}

func gameNotStarted(gm GameManager, req GameRequest, room *Room) error {
	if room.gameStarted {
		req.player.sendString("game already started")
		return errors.New("game already started")
	}
	return nil
}

func isGameStarted(gm GameManager, req GameRequest, room *Room) error {
	if !room.gameStarted {
		req.player.sendString("game not started")
		return errors.New("game not started")
	}
	return nil
}

func isRoomOwner(gm GameManager, req GameRequest, room *Room) error {
	if room.roomOwner != req.player.id {
		req.player.sendString("you are not the owner")
		return errors.New("you are not the owner")
	}
	return nil
}

func isPlayerTurn(gm GameManager, req GameRequest, room *Room) error {
	if room.playerIds[room.curTurn] != req.player.id {
		req.player.sendString("Not your turn")
		return errors.New("not your turn")
	}
	return nil
}

func (gm GameManager) validate(req GameRequest, handler EventHandler, room *Room, validators ...Validator) error {
	for _, validator := range validators {
		err := validator(gm, req, room)
		if err != nil {
			return err
		}
	}
	return handler(room, req)
}

type TileMsg struct {
	MsgType   int          `json:"msgType"`
	TileSides [4]tile_side `json:"tileSides"`
	Crest     bool         `json:"crest"`
	EndsRoad  bool         `json:"endsRoad"`
	Cathedral bool         `json:"cathedral"`
}

func sendDrawnTile(room *Room) error {
	drawnTile := room._game.draw_card()

	type response struct {
		TileMsg
		MyTurn bool `json:"myTurn"`
	}

	tile := TileMsg{
		MsgType:   int(gameStarted),
		TileSides: drawnTile.Sides,
		Crest:     drawnTile.Crest,
		EndsRoad:  drawnTile.Ends_road,
		Cathedral: drawnTile.Cathedral,
	}

	resp := response{
		TileMsg: tile,
		MyTurn:  false,
	}
	room.nextTurn()
	room.playerPingRoom(resp, room.playerIds[room.curTurn])
	resp.MyTurn = true
	room.players[room.playerIds[room.curTurn]].sendStruct(resp)
	return nil
}

func handleStartGame(room *Room, req GameRequest) error {
	room.gameStarted = true
	return sendDrawnTile(room)
}

func handlePlayTile(room *Room, req GameRequest) error {
	err := room._game.place_tile(*req._tile._tile, req._tile.x, req._tile.y)
	if err != nil {
		req.player.sendString("invalid location")
	}

	tile := TileMsg{
		MsgType:   int(tileAdded),
		TileSides: req._tile._tile.Sides,
		Crest:     false,
		EndsRoad:  false,
		Cathedral: false,
	}

	resp := struct {
		TileMsg
		TileX int `json:"tileX"`
		TileY int `json:"tileY"`
	}{
		TileMsg: tile,
		TileX:   req._tile.x,
		TileY:   req._tile.y,
	}

	room.pingRoomStruct(resp)

	return sendDrawnTile(room)

}

func handleTempTile(room *Room, req GameRequest) error {
	tile := TileMsg{
		MsgType:   int(tempTilePlaced),
		TileSides: req._tile._tile.Sides,
		Crest:     false,
		EndsRoad:  false,
		Cathedral: false,
	}
	resp := struct {
		TileMsg
		IsValid bool `json:"isValid"`
		TileX   int  `json:"tileX"`
		TileY   int  `json:"tileY"`
	}{
		TileMsg: tile,
		IsValid: req.validPosition,
		TileX:   req._tile.x,
		TileY:   req._tile.y,
	}

	room.pingRoomStruct(resp)
	return nil
}

func handlePlaceMeeple(room *Room, req GameRequest) error {
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
	return nil
}

func handleRemoveMeeple(room *Room, req GameRequest) error {
	room._game.removeMeeple(req.meeple.x)
	resp := struct {
		MsgType int `json:"msgType"`
		Index   int `json:"index"`
	}{MsgType: int(meepleRemoved), Index: req.meeple.x}
	room.playerPingRoom(resp, req.player.id)
	return nil
}

func handleMoveMeeple(room *Room, req GameRequest) error {
	room._game.moveMeeple(req.index, req.meeple.x, req.meeple.y)
	resp := struct {
		MsgType int `json:"msgType"`
		Index   int `json:"index"`
		X       int `json:"x"`
		Y       int `json:"y"`
	}{MsgType: int(meepleMoved), Index: req.index, X: req.meeple.x, Y: req.meeple.y}
	room.playerPingRoom(resp, req.player.id)
	return nil
}

func newGameManager(rm *RoomManager) GameManager {
	return GameManager{
		RoomManager: rm,
		gameReqs:    make(chan GameRequest),
	}
}

func (gm GameManager) start() {
	for {
		select {
		case req := <-gm.gameReqs:
			room, _ := gm.RoomManager.getRoom(req.roomId)
			switch req.reqType {
			case startGame:
				gm.validate(req, handleStartGame, room, roomExists, gameNotStarted, isRoomOwner)
			case playTile:
				gm.validate(req, handlePlayTile, room, roomExists, isGameStarted, isPlayerTurn)
			case tempTile:
				gm.validate(req, handleTempTile, room, roomExists, isGameStarted, isPlayerTurn)
			case placeMeeple:
				gm.validate(req, handlePlaceMeeple, room, roomExists, isGameStarted, isPlayerTurn)
			case removeMeeple:
				gm.validate(req, handleRemoveMeeple, room, roomExists, isGameStarted, isPlayerTurn)
			case moveMeeple:
				gm.validate(req, handleMoveMeeple, room, roomExists, isGameStarted, isPlayerTurn)
			}
		}
	}
}
