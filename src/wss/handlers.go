package wss

import (
	"log"

	"github.com/LeReverandNox/GuessWhat/src/game"
	"github.com/fatih/structs"
	"golang.org/x/net/websocket"
)

func onConnection(ws *websocket.Conn) *game.Client {
	client := myGame.AddClient(ws)
	sendAllMessages(client)
	myGame.ListClients()
	myGame.ListRooms()
	myGame.ListMessages()

	return client
}

func onDisconnection(client *game.Client, err error) error {
	log.Printf("Socket closed because of : %v", err)
	myGame.RemoveClient(client)
	if room := myGame.GetCurrentClientRoom(client); room != nil {

		room.(*game.Room).RemoveClient(client)
	}
	myGame.ListClients()
	myGame.ListRooms()
	myGame.ListMessages()

	return nil
}

func setNicknameAction(client *game.Client, nickname string) {
	msg := make(map[string]interface{})
	msg["action"] = "set_nickname_cb"
	msg["nickname"] = nickname

	if !myGame.IsNicknameTaken(nickname) {
		client.SetNickname(nickname)
		msg["success"] = true
	} else {
		msg["success"] = false
	}
	client.Socket.SendToSocket(client.Socket, msg)
}

func sendMessageAction(client *game.Client, content string) {
	msg := myGame.AddMessage(client, content)
	msgMap := structs.Map(msg)
	client.Socket.SendToAll(myGame, msgMap)
}

func sendAllMessages(client *game.Client) {
	for _, msg := range myGame.Messages {
		msgMap := structs.Map(msg)
		client.Socket.SendToSocket(client.Socket, msgMap)
	}
}

func joinRoomAction(client *game.Client, roomName string) {
	msg := make(map[string]interface{})
	msg["action"] = "join_room_cb"
	msg["room"] = roomName

	room := myGame.GetRoom(roomName)
	if err := room.AddClient(client); err != nil {
		msg["success"] = false
	} else {
		msg["success"] = true
	}
	client.Socket.SendToSocket(client.Socket, msg)
}
