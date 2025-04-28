package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
)

const (
	PLAYER_TYPE_P1 = "p1"
	PLAYER_TYPE_P2 = "p2"
)

// json marshall should omit user token
type Player struct {
	// mu    sync.RWMutex
	User *User `json:"user"`
	// Ready bool  `json:"ready"`
	// Connected bool  `json:"connected"`
	Type string `json:"type"`
}

type Board [3][3]string

type Room struct {
	mu    sync.RWMutex
	ID    string  `json:"id"`
	P1    *Player `json:"p1"`
	P2    *Player `json:"p2"`
	Board *Board  `json:"board"`
	Turn  string  `json:"turn"`
}

type BoardChoice struct {
	Coords []int `json:"coords"`
}

func (r *Room) HostRoomUpdate(ctx context.Context) {
	client, err := wsserver.ClientsCollection.GetClientByUser(r.P1.User)

	if err != nil {
		fmt.Println("room broadcast client for p1 not found")
	} else {
		reply_data := map[string]any{"room": r}
		client.SendMessage(ctx, EV_ROOM_UPDATE, reply_data)
	}
}

func (r *Room) JoinerRoomUpdate(ctx context.Context) {
	if r.P2 == nil {
		fmt.Println("room broadcast client for p2 not found")
		return
	}

	client, err := wsserver.ClientsCollection.GetClientByUser(r.P2.User)

	if err != nil {
		fmt.Println("room broadcast client for p2 not found")
	} else {
		reply_data := map[string]any{"room": r}
		client.SendMessage(ctx, EV_ROOM_UPDATE, reply_data)
	}
}

func (r *Room) Broadcast(ctx context.Context) {
	r.HostRoomUpdate(ctx)
	r.JoinerRoomUpdate(ctx)
}

type RoomCollection struct {
	mu    sync.RWMutex
	rooms []*Room
}

func (rc *RoomCollection) AddRoom(client *WSClient) (*Room, error) {
	game := Board{}
	p1 := Player{
		User: client.User,
		Type: PLAYER_TYPE_P1,
		// Connected: true,
	}
	room := Room{
		ID:    RandStringBytesMaskImprSrcUnsafe(4),
		P1:    &p1,
		Board: &game,
		Turn:  PLAYER_TYPE_P1,
	}

	rc.mu.Lock()
	rc.rooms = append(rc.rooms, &room)
	rc.mu.Unlock()

	return &room, nil
}

func (rc *RoomCollection) JoinRoom(client *WSClient, roomID string) (*Room, error) {
	rc.mu.RLock()
	idx := slices.IndexFunc(rc.rooms, func(room *Room) bool {
		return room.ID == roomID
	})
	rc.mu.RUnlock()

	if idx == -1 {
		return nil, errors.New("Room not found")
	}

	rc.mu.RLock()
	room := rc.rooms[idx]
	rc.mu.RUnlock()

	p2 := Player{
		User: client.User,
		Type: PLAYER_TYPE_P2,
		// Ready: false,
		// Connected: true,
	}

	room.mu.Lock()
	room.P2 = &p2
	room.mu.Unlock()

	return room, nil
}

func (rc *RoomCollection) RemoveRoom(room *Room) {
	rc.mu.RLock()
	idx := slices.Index(rc.rooms, room)
	rc.mu.RUnlock()

	rc.mu.Lock()
	rc.rooms = slices.Delete(rc.rooms, idx, idx+1)
	rc.mu.Unlock()
}

func (rc *RoomCollection) GetRoomByUser(user *User) (*Room, error) {
	rc.mu.RLock()
	idx := slices.IndexFunc(rc.rooms, func(r *Room) bool {
		if r.P1 != nil {
			if r.P1.User == user {
				return true
			}
		}

		if r.P2 != nil {
			if r.P2.User == user {
				return true
			}
		}

		return false
	})
	rc.mu.RUnlock()

	if idx == -1 {
		return nil, errors.New("Room not found")
	}

	rc.mu.RLock()
	room := rc.rooms[idx]
	rc.mu.RUnlock()

	return room, nil
}
