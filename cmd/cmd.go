package cmd

import (
	"github.com/gorilla/websocket"
	"net/http"
	"go-windows-monitor/utils/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1,
	WriteBufferSize: 1,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PtyHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Websocket upgrade failed: %s\n", err)
	}
	defer conn.Close()

	wp := wsPty{ws: conn}

	wp.Start()

	go wp.writePump()
	wp.readPump()
}
