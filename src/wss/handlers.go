package wss

import (
	"log"

	"github.com/LeReverandNox/GuessWhat/src/game"
	"github.com/fatih/structs"
	"golang.org/x/net/websocket"
)

func onConnection(ws *websocket.Conn) *game.Client {
	client := myGame.AddClient(ws)
	// Send everything the new client needs to know
	sendAllGameMessagesTo(client)
	sendAllGameClientsTo(client)
	sendAllRoomsTo(client)

	// Broadcast his arrival into the game to other clients
	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "incoming_client"
	updateMsg["client"] = client
	client.Socket.Broadcast(myGame, updateMsg)

	myGame.ListClients()
	myGame.ListRooms()
	myGame.ListMessages()

	return client
}

func onDisconnection(client *game.Client, err error) error {
	log.Printf("Socket closed because of : %v", err)
	myGame.RemoveClient(client)
	if room := myGame.GetCurrentClientRoom(client); room != nil {
		trueRoom := room.(*game.Room)
		trueRoom.RemoveClient(client)
		// Broadcast his departure from the channel to other clients
		sendRoomDepartureToAll(client, trueRoom)

		if trueRoom.IsEmpty() {
			myGame.RemoveRoom(trueRoom)
			// Tell everyone about the room suppression.
			sendRoomDeletionToAll(client, trueRoom)
		}
	}

	// Broadcast his departure from the game to other clients
	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "leaving_client"
	updateMsg["client"] = client
	client.Socket.Broadcast(myGame, updateMsg)

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
	room, isNew := myGame.GetRoom(roomName)

	cbMsg := make(map[string]interface{})
	cbMsg["action"] = "join_room_cb"
	cbMsg["room"] = room

	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "incoming_room_client"
	updateMsg["room"] = room
	updateMsg["client"] = client

	if err := room.AddClient(client); err != nil {
		cbMsg["success"] = false
		client.Socket.SendToSocket(client.Socket, cbMsg)
	} else {
		cbMsg["success"] = true
		client.Socket.SendToSocket(client.Socket, cbMsg)
		// If the channel just got created, broadcast it !
		if isNew {
			updateMsg := make(map[string]interface{})
			updateMsg["action"] = "incoming_room"
			updateMsg["room"] = room
			client.Socket.Broadcast(myGame, updateMsg)
		}
		// Send to the client all the infos about the joined room.
		sendAllRoomMessagesTo(client, room)
		sendAllRoomClientsTo(client, room)
		// Notify the others room clients the arrival of the client
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

func sendAllRoomsTo(client *game.Client) {
	rooms := make(map[string]interface{})
	rooms["action"] = "incoming_all_rooms"
	rooms["rooms"] = myGame.Rooms
	client.Socket.SendToSocket(client.Socket, rooms)
}

func sendRoomDepartureToAll(client *game.Client, room *game.Room) {
	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "leaving_room_client"
	updateMsg["client"] = client
	updateMsg["room"] = room
	client.Socket.SendToRoom(room, updateMsg)
}

func sendRoomDeletionToAll(client *game.Client, room *game.Room) {

	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "leaving_room"
	updateMsg["room"] = room
	client.Socket.SendToAll(myGame, updateMsg)
}
