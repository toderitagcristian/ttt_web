package main

import (
	"fmt"
	"net/http"
)

// var PrivateKey []byte = generatePrivateKey()
var PrivateKey []byte = []byte("secret")

var usersDB = UsersCollection{}

var rooms = RoomCollection{}

var wsserver = WSServer{}

func main() {
	http.Handle("GET /ws", &wsserver)
	fmt.Println("Server will start")
	http.ListenAndServe(":3000", nil)
}
