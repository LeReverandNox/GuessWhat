package game

import (
	"log"

	"golang.org/x/net/websocket"
)

type Socket struct {
	Socket *websocket.Conn
}

// var sockets []*websocket.Conn

// func removeConn(ws *websocket.Conn) error {
// 	for index, socket := range sockets {
// 		if socket == ws {
// 			sockets = append(sockets[:index], sockets[index+1:]...)
// 		}
// 	}

// 	return nil
// }

// func brodcast(emitter *websocket.Conn, data Message) error {
// 	for _, socket := range sockets {
// 		if socket != emitter {
// 			if err := websocket.JSON.Send(socket, data); err != nil {
// 				log.Println(err)
// 				return err
// 			}

// 		}
// 	}
// 	return nil
// }

func (socket *Socket) SendToAll(game *Game, data map[string]string) error {
	for _, client := range game.Clients {
		if err := websocket.JSON.Send(client.Socket.Socket, data); err != nil {
			log.Println(err)
			return err

		}
	}
	return nil
}

// func sendTo(receiver *websocket.Conn, data Message) error {
// 	if err := websocket.JSON.Send(receiver, data); err != nil {
// 		log.Println(err)
// 		return err

// 	}
// 	return nil
// }
