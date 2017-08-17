package wss

import (
	"github.com/LeReverandNox/GuessWhat/src/game"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var myGame = game.NewGame()

func socket(ws *websocket.Conn) {
	for {
		// // allocate our container struct
		// var incomingMsg Message

		// // receive a message using the codec
		// if err := websocket.JSON.Receive(ws, &incomingMsg); err != nil {
		// 	removeConn(ws)
		// 	break
		// }

		// log.Printf("Received message %v : %v", incomingMsg.Type, incomingMsg.Message)

		// m4 := Message{Message: "Le serveur viens de send un msg", Type: "message"}
		// sendToAll(m4)
		// m3 := Message{Message: incomingMsg.Message, Type: "message"}
		// brodcast(ws, m3)
	}
}

// StartServer launches the WebSocket server
func StartServer(router *mux.Router) error {
	router.Handle("/ws", websocket.Handler(socket))
	return nil
}
