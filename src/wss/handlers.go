package wss

import (
	"errors"
	"log"

	"github.com/LeReverandNox/GuessWhat/src/game"
	"github.com/fatih/structs"
	"golang.org/x/net/websocket"
)

func onConnection(ws *websocket.Conn) (*game.Client, error) {
	socket := game.NewSocket(ws)
	nickname := ws.Config().Location.Query().Get("nickname")

	cbMsg := make(map[string]interface{})
	cbMsg["action"] = "connexion_cb"
	cbMsg["nickname"] = nickname

	if myGame.IsNicknameTaken(nickname) {
		cbMsg["client"] = nil
		cbMsg["success"] = false
		socket.SendToSocket(socket, cbMsg)
		return nil, errors.New("This nickname is already taken")
	}

	client := myGame.AddClient(socket, nickname)
	cbMsg["client"] = client
	cbMsg["success"] = true
	socket.SendToSocket(socket, cbMsg)

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
	myGame.ListWords()

	return client, nil
}

func onDisconnection(client *game.Client, err error) error {
	log.Printf("Socket closed because of : %v", err)
	myGame.RemoveClient(client)
	if room := myGame.GetCurrentClientRoom(client); room != nil {
		trueRoom := room.(*game.Room)
		isEmpty, _ := trueRoom.RemoveClient(client)
		// Broadcast his departure from the channel to other clients
		sendRoomDepartureToAll(client, trueRoom)

		if isEmpty {
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
	room, isNew := myGame.GetRoom(roomName, client)

	cbMsg := make(map[string]interface{})
	cbMsg["action"] = "join_room_cb"
	cbMsg["room"] = room

	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "incoming_room_client"
	updateMsg["room"] = room
	updateMsg["client"] = client

	if currRoom := myGame.GetCurrentClientRoom(client); currRoom != nil && currRoom != room {
		cbMsg["success"] = false
		cbMsg["reason"] = "You can only be in one room at a time."
		client.Socket.SendToSocket(client.Socket, cbMsg)
	} else if err := room.AddClient(client); err != nil {
		cbMsg["success"] = false
		cbMsg["reason"] = "You are already in this room."
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

func leaveRoomAction(client *game.Client, roomName string) {
	cbMsg := make(map[string]interface{})
	cbMsg["action"] = "leave_room_cb"

	if !myGame.IsRoomExisting(roomName) {
		cbMsg["room"] = roomName
		cbMsg["success"] = false
		cbMsg["reason"] = "This room doesn't exists."
		client.Socket.SendToSocket(client.Socket, cbMsg)
	} else {
		room, _ := myGame.GetRoom(roomName, client)
		cbMsg["room"] = room

		isEmpty, err := room.RemoveClient(client)
		if err != nil {
			cbMsg["success"] = false
			client.Socket.SendToSocket(client.Socket, cbMsg)
		} else {
			cbMsg["success"] = true
			client.Socket.SendToSocket(client.Socket, cbMsg)
			// Broadcast his departure from the channel to other clients
			sendRoomDepartureToAll(client, room)
			if isEmpty {
				myGame.RemoveRoom(room)
				// Tell everyone about the room suppression.
				sendRoomDeletionToAll(client, room)
			}
		}
	}
}

func canvasMouseDownAction(client *game.Client, msg map[string]string) {
	roomName := msg["room"]
	if myGame.IsRoomExisting(roomName) {
		room, _ := myGame.GetRoom(roomName, client)
		if room.IsStarted && room.Drawer.Nickname == client.Nickname {
			updateMsg := make(map[string]interface{})
			updateMsg["action"] = "canvas_mouse_down"
			updateMsg["client"] = client
			updateMsg["room"] = room
			updateMsg["x"] = msg["x"]
			updateMsg["y"] = msg["y"]
			updateMsg["color"] = msg["color"]
			updateMsg["thickness"] = msg["thickness"]
			client.Socket.SendToRoom(room, updateMsg)
		}
	}
}

func canvasMouseMoveAction(client *game.Client, msg map[string]string) {
	roomName := msg["room"]
	if myGame.IsRoomExisting(roomName) {
		room, _ := myGame.GetRoom(roomName, client)
		if room.IsStarted && room.Drawer.Nickname == client.Nickname {
			updateMsg := make(map[string]interface{})
			updateMsg["action"] = "canvas_mouse_move"
			updateMsg["client"] = client
			updateMsg["room"] = room
			updateMsg["x"] = msg["x"]
			updateMsg["y"] = msg["y"]
			updateMsg["color"] = msg["color"]
			updateMsg["thickness"] = msg["thickness"]
			client.Socket.SendToRoom(room, updateMsg)
		}
	}
}

func canvasMouseUpAction(client *game.Client, msg map[string]string) {
	roomName := msg["room"]
	if myGame.IsRoomExisting(roomName) {
		room, _ := myGame.GetRoom(roomName, client)
		if room.IsStarted && room.Drawer.Nickname == client.Nickname {
			updateMsg := make(map[string]interface{})
			updateMsg["action"] = "canvas_mouse_up"
			updateMsg["client"] = client
			updateMsg["room"] = room
			updateMsg["x"] = msg["x"]
			updateMsg["y"] = msg["y"]
			updateMsg["color"] = msg["color"]
			updateMsg["thickness"] = msg["thickness"]
			client.Socket.SendToRoom(room, updateMsg)
		}
	}
}

func startRoomAction(client *game.Client, roomName string) {
	cbMsg := make(map[string]interface{})
	cbMsg["action"] = "start_room_cb"

	if !myGame.IsRoomExisting(roomName) {
		cbMsg["room"] = roomName
		cbMsg["success"] = false
		cbMsg["reason"] = "This room doesn't exists."
		client.Socket.SendToSocket(client.Socket, cbMsg)
	} else {
		room, _ := myGame.GetRoom(roomName, client)

		cbMsg["room"] = room

		if !room.IsOwner(client) {
			cbMsg["success"] = false
			cbMsg["reason"] = "You are not the owner of this room."
			client.Socket.SendToSocket(client.Socket, cbMsg)
		} else if room.IsStarted {
			cbMsg["success"] = false
			cbMsg["reason"] = "This room is already started."
			client.Socket.SendToSocket(client.Socket, cbMsg)
		} else if room.GetNbClients() < 2 {
			cbMsg["success"] = false
			cbMsg["reason"] = "You have to be at least 2 players to start a room."
			client.Socket.SendToSocket(client.Socket, cbMsg)
		} else {
			room.Start()

			// BOUCLE DEBUT JEU
			// Pick and set random drawer
			drawer := room.PickRandomClient()
			room.SetDrawer(drawer)
			// Pick and set random word
			word := myGame.PickRandomWord()
			room.SetWord(word)
			// Start timer
			// BOUCLE FIN JEU

			// Send to room clients about it's state
			updateMsg := make(map[string]interface{})
			updateMsg["action"] = "room_start"
			updateMsg["room"] = room
			client.Socket.SendToRoom(room, updateMsg)

			cbMsg["success"] = true
			client.Socket.SendToSocket(client.Socket, cbMsg)
		}
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
