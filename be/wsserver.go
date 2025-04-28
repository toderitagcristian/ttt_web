package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	EV_ERROR           = "error"
	EV_CONNECT_SUCCESS = "connect_success"

	EV_LOGIN         = "login"
	EV_LOGIN_SUCCESS = "login_success"
	EV_LOGIN_ERROR   = "login_error"

	EV_ROOM_CREATE         = "room_create"
	EV_ROOM_CREATE_SUCCESS = "room_create_success"
	EV_ROOM_CREATE_ERROR   = "room_create_error"

	EV_ROOM_JOIN         = "room_join"
	EV_ROOM_JOIN_SUCCESS = "room_join_success"
	EV_ROOM_JOIN_ERROR   = "room_join_error"

	EV_ROOM_UPDATE = "room_update"

	EV_BOARD_CHOICE       = "board_choice"
	EV_BOARD_CHOICE_ERROR = "board_choice_error"
)

type WsMessage struct {
	Event string         `json:"event"`
	Data  map[string]any `json:"data"`
}

type WSServer struct {
	ClientsCollection WSClientsCollection
}

func (wsserver *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		fmt.Println("webscocket accept err", err)
	}
	defer c.CloseNow()

	fmt.Println("[IN] connect")

	ctx := context.Background()

	client, err := wsserver.ClientsCollection.AddClient(c)

	if err != nil {
		reply_data := map[string]any{"msg": err}
		client.SendMessage(ctx, EV_ERROR, reply_data)
		return
	}

	client.SendMessage(ctx, EV_CONNECT_SUCCESS, nil)

	// Can be used for JWT token and auth
	token := r.Header.Get("Sec-Websocket-Protocol")

	if token != "" {
		fmt.Println("[IN]", EV_LOGIN, client.UUID)
		err = client.User.ValidateToken(token)

		if err != nil {
			reply_data := map[string]any{"msg": err}
			client.SendMessage(ctx, EV_LOGIN_ERROR, reply_data)
		} else {
			user, err := usersDB.GetUserByToken(token)

			if err != nil {
				reply_data := map[string]any{"msg": err}
				client.SendMessage(ctx, EV_LOGIN_ERROR, reply_data)
			} else {
				client.SetUser(user)
				client.LoginHandler(ctx)
			}
		}
	}

	for {
		var msg WsMessage

		err := wsjson.Read(ctx, c, &msg)
		if err != nil {
			client.mu.RLock()
			fmt.Println("event: disconnect", client.UUID)

			if client.User != nil {
				fmt.Println("event: disconnect", client.User.Nickname)

				// room, _ := rooms.GetRoomByUser(client.User)

				// if room != nil {
				// 	room.mu.RLock()
				// 	if room.P1.User == client.User {
				// 		room.P1.mu.Lock()
				// 		room.P1.Connected = false
				// 		room.P1.mu.Lock()

				// 		room.JoinerRoomUpdate(ctx)
				// 	} else {
				// 		room.P1.mu.Lock()
				// 		room.P2.Connected = false
				// 		room.P1.mu.Lock()

				// 		room.HostRoomUpdate(ctx)
				// 	}
				// 	room.mu.RUnlock()
				// }
			}
			client.mu.RUnlock()

			wsserver.ClientsCollection.RemoveClient(client)
			return
		}

		go func() {
			fmt.Println("[IN]", msg.Event)
			bytes, err := json.Marshal(&msg.Data)
			if err != nil {
				fmt.Println("json Marshal err: ", err)
			}

			switch msg.Event {
			case EV_LOGIN:
				cmdLogin(ctx, client, &bytes)
			case EV_ROOM_CREATE:
				cmdRoomCreate(ctx, client)
			case EV_ROOM_JOIN:
				cmdRoomJoin(ctx, client, &bytes)
			case EV_BOARD_CHOICE:
				cmdBoardChoice(ctx, client, &bytes)
			default:
				reply_data := map[string]any{"msg": "Event unknown"}
				client.SendMessage(ctx, EV_ERROR, reply_data)
			}
		}()
	}
}
