package wss

import (
	"log"

	"github.com/LeReverandNox/GuessWhat/src/game"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var myGame = game.NewGame()

func parseMessage(client *game.Client, msg map[string]string) {
	log.Print("On recoi un msg a traité ", msg)
	switch msg["action"] {
	case "set_nickname":
		setNicknameAction(client, msg["nickname"])
	case "send_message":
		sendMessageAction(client, msg["content"])
	case "join_room":
		joinRoomAction(client, msg["room"])
	case "leave_room":
		leaveRoomAction(client, msg["room"])
	}

	myGame.ListClients()
	myGame.ListRooms()
	myGame.ListMessages()
}

func socketHandler(ws *websocket.Conn) {
	client, connectionError := onConnection(ws)
	if connectionError != nil {
		ws.Close()
	}

	for {
		var msg map[string]string
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			if connectionError == nil {
				onDisconnection(client, err)
			}
			break
		}
		parseMessage(client, msg)
	}
}

// StartServer launches the WebSocket server
func StartServer(router *mux.Router) error {
	router.Handle("/ws", websocket.Handler(socketHandler))
	return nil
}
