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

func (socket *Socket) Broadcast(game *Game, data map[string]interface{}) error {
	for _, client := range game.Clients {
		if client.Socket != socket {
			if err := websocket.JSON.Send(client.Socket.Socket, data); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (socket *Socket) SendToAll(game *Game, data map[string]interface{}) error {
	for _, client := range game.Clients {
		if err := websocket.JSON.Send(client.Socket.Socket, data); err != nil {
			log.Println(err)
			return err

		}
	}
	return nil
}

func (socket *Socket) SendToSocket(receiver *Socket, data map[string]interface{}) error {
	if err := websocket.JSON.Send(receiver.Socket, data); err != nil {
		log.Println(err)
		return err

	}
	return nil
}
