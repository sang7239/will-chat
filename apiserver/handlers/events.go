package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/will-slack/apiserver/events"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//WebSocketUpgradeHandler handles websocket upgrade requests
func (ctx *Context) WebSocketUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	// state := &SessionState{}
	// if _, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state); err != nil {
	// 	http.Error(w, "error retrieving current user", http.StatusUnauthorized)
	// 	return
	// }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error adding client: %s", err.Error())
		return
	}
	ctx.Notifier.AddClient(conn)
}

func (ctx *Context) notify(dType string, data interface{}) {
	// create a new event so we can add it to the notifications queue
	event := &events.Event{
		Type: dType,
		Data: data,
	}
	ctx.Notifier.Notify(event)
}
