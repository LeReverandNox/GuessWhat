package wss

import (
	"log"

	"github.com/LeReverandNox/GuessWhat/src/game"
	"github.com/fatih/structs"
	"golang.org/x/net/websocket"
)

func onConnection(ws *websocket.Conn) *game.Client {
	client := myGame.AddClient(ws)
	sendAllGameMessagesTo(client)
	sendAllGameClientsTo(client)

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
	cbMsg := make(map[string]interface{})
	cbMsg["action"] = "set_nickname_cb"
	cbMsg["nickname"] = nickname

	if !myGame.IsNicknameTaken(nickname) {
		client.SetNickname(nickname)
		cbMsg["success"] = true
	} else {
		cbMsg["success"] = false
	}
	client.Socket.SendToSocket(client.Socket, cbMsg)
}

func sendMessageAction(client *game.Client, content string) {
	if room := myGame.GetCurrentClientRoom(client); room != nil {
		trueRoom := room.(*game.Room)
		msg := trueRoom.AddMessage(client, content)
		msgMap := structs.Map(msg)
		msgMap["action"] = "incoming_room_message"
		msgMap["channel"] = trueRoom.Name

		client.Socket.SendToRoom(trueRoom, msgMap)
	} else {
		msg := myGame.AddMessage(client, content)
		msgMap := structs.Map(msg)
		msgMap["action"] = "incoming_global_message"

		client.Socket.SendToAll(myGame, msgMap)
	}
}

func joinRoomAction(client *game.Client, roomName string) {
	cbMsg := make(map[string]interface{})
	cbMsg["action"] = "join_room_cb"
	cbMsg["room"] = roomName

	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "incoming_room_client"
	updateMsg["room"] = roomName
	updateMsg["client"] = client

	room := myGame.GetRoom(roomName)
	if err := room.AddClient(client); err != nil {
		cbMsg["success"] = false
		client.Socket.SendToSocket(client.Socket, cbMsg)
	} else {
		cbMsg["success"] = true
		client.Socket.SendToSocket(client.Socket, cbMsg)
		sendAllRoomMessagesTo(client, room)
		sendAllRoomClientsTo(client, room)
		client.Socket.BroadcastToRoom(room, updateMsg)
	}
}

func sendAllGameMessagesTo(client *game.Client) {
	messages := make(map[string]interface{})
	messages["action"] = "incoming_all_global_message"
	messages["messages"] = myGame.Messages
	client.Socket.SendToSocket(client.Socket, messages)
}

func sendAllGameClientsTo(client *game.Client) {
	clients := make(map[string]interface{})
	clients["action"] = "incoming_all_global_users"
	clients["clients"] = myGame.Clients
	client.Socket.SendToSocket(client.Socket, clients)
}

func sendAllRoomMessagesTo(client *game.Client, room *game.Room) {
	messages := make(map[string]interface{})
	messages["action"] = "incoming_all_room_message"
	messages["channel"] = room.Name
	messages["messages"] = room.Messages
	client.Socket.SendToSocket(client.Socket, messages)
}

func sendAllRoomClientsTo(client *game.Client, room *game.Room) {
	clients := make(map[string]interface{})
	clients["action"] = "incoming_all_room_clients"
	clients["channel"] = room.Name
	clients["clients"] = room.Clients
	client.Socket.SendToSocket(client.Socket, clients)
}
