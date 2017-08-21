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
	return client
}

func onDisconnection(client *game.Client, err error) error {
	log.Printf("Socket closed because of : %v", err)
	myGame.RemoveClient(client)
	return nil
}

func setNicknameAction(client *game.Client, nickname string) {
	client.SetNickname(nickname)

	msg := make(map[string]interface{})
	// DEBUG
	msg["message"] = "Je s'apelle " + client.Nickname
	client.Socket.SendToAll(myGame, msg)
	// DEBUG
}

func sendMessageAction(client *game.Client, content string) {
	msg := myGame.AddMessage(client.Nickname, content)
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
	log.Printf("%v veut rejoindre %v", client.Nickname, roomName)
	msg := make(map[string]interface{})

	room := myGame.GetRoom(roomName)
	if err := room.AddClient(client); err != nil {
		msg["join_room_cb"] = false
	} else {
		msg["join_room_cb"] = true
	}
	client.Socket.SendToSocket(client.Socket, msg)
	myGame.ListRooms()
}
