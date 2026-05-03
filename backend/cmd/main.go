package main

import (
	"net/http"

	websocket "github.com/Blue-Onion/WebGocket/handler/webSocket"
)

func main() {
	http.HandleFunc("/ws", websocket.WsHanlder)
	http.ListenAndServe(":3480", nil)
}
