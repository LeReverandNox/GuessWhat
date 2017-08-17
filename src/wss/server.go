package wss

import (
	"fmt"

	"github.com/LeReverandNox/GuessWhat/src/game"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var myGame = game.NewGame()

func parseMessage() {

}

func onConnection(ws *websocket.Conn) *game.Client {
	client := myGame.AddClient(ws)
	return client
}

func socketHandler(ws *websocket.Conn) {
	client := onConnection(ws)

	for {
		var msg map[string]string
		// receive a message using the codec
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			fmt.Println(err)
			// removeConn(ws)
			break
		}

		// if err := websocket.Message.Receive(ws, &msg); err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(msg)
		fmt.Println(msg["type"], msg["message"])
		// log.Printf("Received message %v : %v", incomingMsg.Type, incomingMsg.Message)

		// m4 := Message{Message: "Le serveur viens de send un msg", Type: "message"}
		// sendToAll(m4)
		// m3 := Message{Message: incomingMsg.Message, Type: "message"}
		// brodcast(ws, m3)
	}
}

// StartServer launches the WebSocket server
func StartServer(router *mux.Router) error {
	router.Handle("/ws", websocket.Handler(socketHandler))
	return nil
}
