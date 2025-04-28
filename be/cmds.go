package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
)

type LoginData struct {
	Nickname string `json:"nickname"`
}

func cmdLogin(ctx context.Context, client *WSClient, bytes *[]byte) {
	var data LoginData
	err := json.Unmarshal(*bytes, &data)
	if err != nil {
		reply_data := map[string]any{"msg": err}
		client.SendMessage(ctx, EV_LOGIN_ERROR, reply_data)
		return
	}

	var user *User

	if client.User != nil {
		user = client.User
	} else {
		newUser, err := usersDB.AddUser(data.Nickname)

		if err != nil {
			reply_data := map[string]any{"msg": err}
			client.SendMessage(ctx, EV_LOGIN_ERROR, reply_data)
			return
		}

		user = newUser
	}

	client.SetUser(user)
	client.LoginHandler(ctx)
}

func cmdRoomCreate(ctx context.Context, client *WSClient) {
	room, err := rooms.AddRoom(client)

	if err != nil {
		reply_data := map[string]any{"msg": err}
		client.SendMessage(ctx, EV_ROOM_CREATE_ERROR, reply_data)
		return
	}

	reply_data := map[string]any{"room": room}
	client.SendMessage(ctx, EV_ROOM_CREATE_SUCCESS, reply_data)
}

func cmdRoomJoin(ctx context.Context, client *WSClient, bytes *[]byte) {
	var data Room

	err := json.Unmarshal(*bytes, &data)

	if err != nil {
		reply_data := map[string]any{"msg": err}
		client.SendMessage(ctx, EV_ROOM_JOIN_ERROR, reply_data)
		return
	}

	room, err := rooms.JoinRoom(client, data.ID)

	if err != nil {
		reply_data := map[string]any{"msg": err}
		client.SendMessage(ctx, EV_ROOM_JOIN_ERROR, reply_data)
		return
	}

	reply_data := map[string]any{"room": room}
	client.SendMessage(ctx, EV_ROOM_JOIN_SUCCESS, reply_data)

	// send room_update to host
	room.HostRoomUpdate(ctx)
}

func cmdBoardChoice(ctx context.Context, client *WSClient, bytes *[]byte) {
	var data BoardChoice

	err := json.Unmarshal(*bytes, &data)

	if err != nil {
		reply_data := map[string]any{"msg": err}
		client.SendMessage(ctx, EV_BOARD_CHOICE_ERROR, reply_data)
		return
	}

	room, err := rooms.GetRoomByUser(client.User)

	if err != nil {
		reply_data := map[string]any{"msg": err}
		client.SendMessage(ctx, EV_BOARD_CHOICE_ERROR, reply_data)
		return
	}

	var user *User
	var playerValue string

	room.mu.Lock()

	if room.Turn == PLAYER_TYPE_P1 {
		user = room.P1.User
		playerValue = "X"
	} else {
		user = room.P2.User
		playerValue = "0"
	}

	if user != client.User {
		reply_data := map[string]any{"msg": "Not your turn"}
		client.SendMessage(ctx, EV_BOARD_CHOICE_ERROR, reply_data)
		return
	}

	prevValue := room.Board[data.Coords[0]][data.Coords[1]]

	if prevValue != "" {
		reply_data := map[string]any{"msg": "Cannot choose this block"}
		client.SendMessage(ctx, EV_BOARD_CHOICE_ERROR, reply_data)
		return
	}

	room.Board[data.Coords[0]][data.Coords[1]] = playerValue

	// Check win conditions
	hasWinner := false

	win := [][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {0, 2}},
		{{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}

	for _, v := range win {
		v1 := room.Board[v[0][0]][v[0][1]]
		v2 := room.Board[v[1][0]][v[1][1]]
		v3 := room.Board[v[1][0]][v[1][1]]

		hasEmptyString := slices.Contains([]string{v1, v2, v3}, "")

		if hasEmptyString {
			continue
		}

		c1 := room.Board[v[0][0]][v[0][1]] == room.Board[v[1][0]][v[1][1]]
		c2 := room.Board[v[0][0]][v[0][1]] == room.Board[v[2][0]][v[2][1]]

		if c1 && c2 {
			hasWinner = true
			break
		}
	}

	if hasWinner {
		fmt.Println("winner", user.Nickname, playerValue)
		room.Turn = ""
	} else {
		// Check no more moves
		hasDraw := true

	main:
		for _, v := range room.Board {
			for _, k := range v {
				if k == "" {
					hasDraw = false
					break main
				}
			}
		}

		if hasDraw {
			fmt.Println("draw", user.Nickname, playerValue)
			room.Turn = ""
		} else {
			if room.Turn == PLAYER_TYPE_P1 {
				room.Turn = PLAYER_TYPE_P2
			} else {
				room.Turn = PLAYER_TYPE_P1
			}
		}
	}

	room.Broadcast(ctx)

	room.mu.Unlock()
}
