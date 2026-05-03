package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	UserName string
}
type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var Clients = make(map[*websocket.Conn]*Client)
var mu sync.Mutex

func broadCast(msg Message) {
	mu.Lock()
	defer mu.Unlock()
	for conn := range Clients {
		err := conn.WriteJSON(msg)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			delete(Clients, conn)
		}
	}
}
func WsHanlder(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "User Name is requireed", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	mu.Lock()
	Clients[conn] = &Client{
		Conn:     conn,
		UserName: username,
	}
	mu.Unlock()
	broadCast(Message{
		From:    "System",
		Content: username + "Joined the chat",
	})
	defer func() {
		mu.Lock()
		delete(Clients, conn)
		mu.Unlock()
		conn.Close()
		broadCast(Message{
			From:    "System",
			Content: username + "Leaved the chat",
		})
	}()
	for {
		_, msgByte, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg Message
		err = json.Unmarshal(msgByte, &msg)

		if err != nil {
			continue
		}
		msg.From = username
		broadCast(msg)
	}
}
