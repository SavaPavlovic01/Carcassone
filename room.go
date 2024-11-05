package main

import "errors"

type Room struct {
	_game       game
	roomOwner   string
	gameStarted bool
	players     map[string]*Player
	playerIds   []string
	curTurn     int
}

func (r *Room) addPlayer(player *Player) error {
	_, inGame := r.players[player.id]
	if inGame {
		return errors.New("already in game")
	}

	r.players[player.id] = player
	r.playerIds = append(r.playerIds, player.id)
	return nil
}

func (r *Room) pingRoom(msg string) {
	for _, player := range r.players {
		player.sendString(msg)
	}
}

func (r *Room) pingRoomStruct(msg interface{}) {
	for _, player := range r.players {
		player.sendStruct(msg)
	}
}

func (r *Room) playerPingRoom(msg interface{}, playerId string) {
	for _, player := range r.players {
		if player.id == playerId {
			continue
		}
		player.sendStruct(msg)
	}
}

func (r *Room) nextTurn() int {
	r.curTurn++
	if r.curTurn == len(r.players) {
		r.curTurn = 0
	}
	return r.curTurn
}

func newRoom(owner *Player) *Room {
	return &Room{
		_game:       new_game(),
		roomOwner:   owner.id,
		gameStarted: false,
		players:     map[string]*Player{owner.id: owner},
		playerIds:   []string{owner.id},
	}
}
