package wss

import (
	"errors"
	"log"
	"time"

	"github.com/LeReverandNox/GuessWhat/src/tools"

	"github.com/LeReverandNox/GuessWhat/src/game"
	"github.com/fatih/structs"
	"golang.org/x/net/websocket"
)

func onConnection(ws *websocket.Conn) (*game.Client, error) {
	socket := game.NewSocket(ws)
	nickname := ws.Config().Location.Query().Get("nickname")

	cbMsg := make(map[string]interface{})
	cbMsg["action"] = "connexion_cb"
	sanitizedNickname := tools.Sanitize(nickname)
	cbMsg["nickname"] = sanitizedNickname

	if len(sanitizedNickname) < 1 {
		cbMsg["client"] = nil
		cbMsg["reason"] = "This nickname is too short."
		cbMsg["success"] = false
		socket.SendToSocket(socket, cbMsg)
		return nil, errors.New("This nickname is too short.")
	}

	if myGame.IsNicknameTaken(sanitizedNickname) {
		cbMsg["client"] = nil
		cbMsg["success"] = false
		cbMsg["reason"] = "This nickname is already taken."
		socket.SendToSocket(socket, cbMsg)
		return nil, errors.New("This nickname is already taken.")
	}

	client := myGame.AddClient(socket, sanitizedNickname)
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
	content = tools.Sanitize(content)
	if len(content) > 0 {
		if roomInt := myGame.GetCurrentClientRoom(client); roomInt != nil {
			room := roomInt.(*game.Room)
			msg := room.AddMessage(client, content)

			hasWon := false
			if room.IsRoundGoing && client.Nickname != room.Drawer.Nickname {
				hasWon = parseForAnswer(client, room, msg)
			}

			if !hasWon {
				// Send the message to the room
				msgMap := structs.Map(msg)
				msgMap["action"] = "incoming_room_message"
				msgMap["channel"] = room.Name
				client.Socket.SendToRoom(room, msgMap)
			} else {
				// Send a 'youFindTheWord' event
				// Send a 'XXX find the word' event

				if room.GetNbWinners() == room.GetNbClients()-1 {
					endRound(client, room, "EVERYONE_WINS")
				}
			}
		} else {
			msg := myGame.AddMessage(client, content)
			msgMap := structs.Map(msg)
			msgMap["action"] = "incoming_global_message"

			client.Socket.SendToAll(myGame, msgMap)
		}
	}
}

func joinRoomAction(client *game.Client, roomName string) {
	roomName = tools.Sanitize(roomName)
	if len(roomName) > 0 {
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

			// Send the current image to the client if the room is started
			if room.IsRoundGoing {
				room.AddDrawingNeeder(client)
				askDrawerForImage(room)
			}
			// Notify the others room clients the arrival of the client
			client.Socket.BroadcastToRoom(room, updateMsg)
		}
	}
}

func leaveRoomAction(client *game.Client, roomName string) {
	roomName = tools.Sanitize(roomName)
	if len(roomName) > 0 {
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
				cbMsg["reason"] = "You are not in this room."
				client.Socket.SendToSocket(client.Socket, cbMsg)
			} else {
				// Reset the client points
				client.ResetScore()

				cbMsg["success"] = true
				cbMsg["me"] = client
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
}

func canvasMouseDownAction(client *game.Client, msg map[string]string) {
	roomName := msg["room"]
	if myGame.IsRoomExisting(roomName) {
		room, _ := myGame.GetRoom(roomName, client)
		if room.IsRoundGoing && room.Drawer.Nickname == client.Nickname {
			updateMsg := make(map[string]interface{})
			updateMsg["action"] = "canvas_mouse_down"
			updateMsg["client"] = client
			updateMsg["room"] = room
			updateMsg["toX"] = msg["toX"]
			updateMsg["toY"] = msg["toY"]
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
		if room.IsRoundGoing && room.Drawer.Nickname == client.Nickname {
			updateMsg := make(map[string]interface{})
			updateMsg["action"] = "canvas_mouse_move"
			updateMsg["client"] = client
			updateMsg["room"] = room
			updateMsg["fromX"] = msg["fromX"]
			updateMsg["fromY"] = msg["fromY"]
			updateMsg["toX"] = msg["toX"]
			updateMsg["toY"] = msg["toY"]
			updateMsg["color"] = msg["color"]
			updateMsg["thickness"] = msg["thickness"]
			client.Socket.SendToRoom(room, updateMsg)
		}
	}
}

func startRoomAction(client *game.Client, roomName string) {
	roomName = tools.Sanitize(roomName)
	if len(roomName) > 0 {
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

				startRound(client, room)

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
}

func sendImageAction(client *game.Client, msg map[string]string) {
	roomName := msg["room"]
	if myGame.IsRoomExisting(roomName) {
		room, _ := myGame.GetRoom(roomName, client)
		if room.IsRoundGoing && room.Drawer.Nickname == client.Nickname {
			room.SetImage(msg["image"])
			for _, clientToSendTo := range room.NeedingDrawing {
				sendRoomImageTo(clientToSendTo, room)
			}
			room.CleanDrawingNeeders()
		}
	}
}

func cleanCanvasAction(client *game.Client, roomName string) {
	if myGame.IsRoomExisting(roomName) {
		room, _ := myGame.GetRoom(roomName, client)
		if room.IsRoundGoing && room.Drawer.Nickname == client.Nickname {
			room.ResetImage()

			updateMsg := make(map[string]interface{})
			updateMsg["action"] = "clean_canvas"
			updateMsg["room"] = room
			client.Socket.SendToRoom(room, updateMsg)
		}
	}
}

// Non actions

func parseForAnswer(proposer *game.Client, room *game.Room, message *game.Message) bool {
	dist := tools.Distance(message.Content, room.Word.Value)
	if dist == 0 {
		room.AddWinner(proposer)
		return true
	} else if dist <= 2 {
		updateMsg := make(map[string]interface{})
		updateMsg["action"] = "close_word"
		updateMsg["room"] = room
		updateMsg["proposed_word"] = message.Content
		proposer.Socket.SendToSocket(proposer.Socket, updateMsg)
		return false
	}
	return false
}

func startRound(client *game.Client, room *game.Room) {
	// Pick and set random drawer
	drawer := room.PickRandomClient()
	room.SetDrawer(drawer)
	// Pick and set random word
	word := myGame.PickRandomWord()
	room.SetWord(word)
	room.ResetImage()
	room.CleanWinners()
	room.IncrementRound()
	room.StartRound()
	handleRoundTimer(client, room)

	// Send to room clients about it's state
	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "new_round_start"
	updateMsg["room"] = room
	updateMsg["drawer"] = drawer
	updateMsg["word"] = room.Word.Value
	client.Socket.SendToRoom(room, updateMsg)
}

func handleRoundTimer(client *game.Client, room *game.Room) {
	roundTicker := room.SetTicker(time.NewTicker(time.Second * 1))
	roundEnd := room.SetRoundEnd(time.Now().Local().Add(time.Second * time.Duration(room.RoundDuration+1)))

	go func() {
		i := 1
		defer func() {
			room.StopTicker()
			endRound(client, room, "TIMESUP")
		}()

		for t := range roundTicker.C {
			if t.After(roundEnd) {
				break
			}

			// TODO: implement letter revelation according to iteration
			// log.Printf("Iteration : %v", i)
			i++
			room.PassedSeconds = i
		}
	}()
}

func revealWordLetter(client *game.Client, room *game.Room) {
	log.Printf("On va reveler une lettre")
}

func endRound(client *game.Client, room *game.Room, reason string) {
	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "round_end"
	updateMsg["room"] = room

	room.StopRound()

	switch reason {
	case "EVERYONE_WINS":
		updateMsg["clients"] = room.Clients
		updateMsg["room"] = room
		updateMsg["reason"] = "EVERYONE_WINS"

		room.ComputeClientsPoints()
		client.Socket.SendToRoom(room, updateMsg)
	}

	// Wait a moment, so clients can see score and stuff.
	time.Sleep(5 * time.Second)

	if room.GetNbClients() >= 2 && room.ActualRound < room.TotalRounds {
		startRound(client, room)
	} else {

		// stop room
	}
}

func askDrawerForImage(room *game.Room) {
	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "ask_for_image"
	updateMsg["room"] = room
	room.Drawer.Socket.SendToSocket(room.Drawer.Socket, updateMsg)
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

func sendRoomImageTo(client *game.Client, room *game.Room) {
	updateMsg := make(map[string]interface{})
	updateMsg["action"] = "incoming_room_image"
	updateMsg["room"] = room
	client.Socket.SendToSocket(client.Socket, updateMsg)
}
