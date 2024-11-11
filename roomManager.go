package main

import (
	"errors"
	"strconv"
)

type RoomRequestType int

const (
	makeNewRoom RoomRequestType = iota
	playerJoinRoom
)

type RoomRequest struct {
	player  *Player
	reqType RoomRequestType
	room    string
}

type RoomManager struct {
	rooms       map[string]*Room
	makeRoomReq chan RoomRequest
}

func (rm RoomManager) addRoom(owner *Player) error {
	if _, exists := rm.rooms[owner.id]; exists {
		return errors.New("room already exists")
	}

	rm.rooms[owner.id] = newRoom(owner)
	return nil
}

func (rm RoomManager) getRoom(roomId string) (*Room, error) {
	room, exists := rm.rooms[roomId]
	if !exists {
		return nil, errors.New("room doesnt exist")
	}
	return room, nil
}

func (rm RoomManager) joinRoom(player *Player, roomId string) error {
	room, exists := rm.rooms[roomId]
	if !exists {
		player.sendString("room doesnt exist")
		return errors.New("room doesnt exist")
	}
	err := room.addPlayer(player)
	if err != nil {
		player.sendString("already in room")
		return err
	}
	return nil
}

// TODO: refactor
func (rm RoomManager) start() {
	for {
		select {
		case req := <-rm.makeRoomReq:
			if req.reqType == makeNewRoom {
				//fmt.Printf("%+v\n", req.player)
				err := rm.addRoom(req.player)
				if err != nil {
					req.player.sendString("room already exists")
				} else {
					req.player.sendString("1")
					req.player.inRoom = true
				}
			}
			if req.reqType == playerJoinRoom {
				err := rm.joinRoom(req.player, req.room)
				if err != nil {
					continue
				}
				req.player.inRoom = true
				req.player.sendString("OK")
				//fmt.Printf("%+v\n", rm.rooms)
				rm.rooms[req.room].pingRoom(strconv.Itoa(len(rm.rooms[req.room].players)))
			}

		}
	}
}
