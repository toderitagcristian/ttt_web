package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type WSClient struct {
	mu   sync.RWMutex
	Conn *websocket.Conn
	UUID string
	User *User
}

func (client *WSClient) SetUser(user *User) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.User = user
}

func (client *WSClient) SendMessage(ctx context.Context, event string, data map[string]any) {
	reply := WsMessage{Event: event, Data: data}
	err := wsjson.Write(ctx, client.Conn, reply)

	if err != nil {
		fmt.Println("wsjson.Write err: ", event, err, data)
		return
	}

	user_nickname := "no_user"
	if client.User != nil {
		user_nickname = client.User.Nickname
	}

	fmt.Println("[OUT]", event, user_nickname, data, client.UUID)
}

func (client *WSClient) LoginHandler(ctx context.Context) {
	reply_data := map[string]any{"user": client.User}
	client.SendMessage(ctx, EV_LOGIN_SUCCESS, reply_data)

	// send room_update if applicable
	client.RoomUpdateHandler(ctx)
}

func (client *WSClient) RoomUpdateHandler(ctx context.Context) {
	room, err := rooms.GetRoomByUser(client.User)

	// if room.P1.User == client.User {
	// 	room.P1.Connected = true
	// } else {
	// 	room.P2.Connected = true
	// }

	if err != nil {
		fmt.Println("RoomUpdateHandler", err)
		return
	}

	reply_data := map[string]any{"room": room}
	client.SendMessage(ctx, EV_ROOM_UPDATE, reply_data)
}

type WSClientsCollection struct {
	mu      sync.RWMutex
	Clients []*WSClient
}

func (wscc *WSClientsCollection) AddClient(conn *websocket.Conn) (*WSClient, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	client := WSClient{
		Conn: conn,
		UUID: uuid.String(),
	}

	wscc.mu.Lock()
	wscc.Clients = append(wscc.Clients, &client)
	wscc.mu.Unlock()

	return &client, nil
}

func (wscc *WSClientsCollection) RemoveClient(client *WSClient) {
	wscc.mu.RLock()
	i := slices.Index(wscc.Clients, client)
	wscc.mu.RUnlock()

	wscc.mu.Lock()
	wscc.Clients = slices.Delete(wscc.Clients, i, i+1)
	wscc.mu.Unlock()
}

func (wscc *WSClientsCollection) GetClientByUser(user *User) (*WSClient, error) {
	wscc.mu.RLock()
	idx := slices.IndexFunc(wscc.Clients, func(c *WSClient) bool {
		return c.User == user
	})
	wscc.mu.RUnlock()

	if idx == -1 {
		return nil, errors.New("WsClient not found by user")
	}

	return wscc.Clients[idx], nil
}
