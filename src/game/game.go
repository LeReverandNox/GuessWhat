package game

import (
	"errors"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type Game struct {
	Clients  []*Client
	Rooms    []*Room
	Messages []*Message
}

// NewGame creates a new Game struct, and returns it
func NewGame() *Game {
	game := Game{}
	game.Clients = make([]*Client, 0)
	game.Rooms = make([]*Room, 0)
	game.Messages = make([]*Message, 0)

	return &game
}

// AddClient adds a client to the game
func (game *Game) AddClient(ws *websocket.Conn) *Client {
	client := NewClient(ws)
	game.Clients = append(game.Clients, client)
	return client

}

// RemoveClient remove a client from the game
func (game *Game) RemoveClient(clientToDelete *Client) {
	for i, client := range game.Clients {
		if client.Socket == clientToDelete.Socket {
			game.Clients = append(game.Clients[:i], game.Clients[i+1:]...)
		}
	}
}

// AddRoom adds a room to the server
func (game *Game) AddRoom(name string) (interface{}, error) {
	if !game.isRoomExisting(name) {
		room := NewRoom(name)
		game.Rooms = append(game.Rooms, room)
		return room, nil
	}
	return nil, errors.New("The channel " + name + " already exist")
}

// GetRoom returns the desired room. If not existing, creates it before.
func (game *Game) GetRoom(name string) (*Room, bool) {
	if game.isRoomExisting(name) {
		for _, room := range game.Rooms {
			if room.Name == name {
				return room, false
			}
		}
	}
	room, _ := game.AddRoom(name)
	return room.(*Room), true
}

// GetCurrentClientRoom returns the current room of a client.
func (game *Game) GetCurrentClientRoom(client *Client) interface{} {
	for _, room := range game.Rooms {
		if room.IsClientIn(client) {
			return room
		}
	}
	return nil
}

// AddMessage adds a message to the game
func (game *Game) AddMessage(sender *Client, content string) *Message {
	message := NewMessage(sender, content, time.Now())
	game.Messages = append(game.Messages, message)

	return message
}

// IsNicknameTaken return true if the given nickname already exists in the game
func (game *Game) IsNicknameTaken(nicknameToTest string) bool {
	for _, client := range game.Clients {
		if client.Nickname == nicknameToTest {
			return true
		}
	}
	return false
}

// Privates methods

func (game *Game) isRoomExisting(name string) bool {
	for _, room := range game.Rooms {
		if room.Name == name {
			return true
		}
	}
	return false
}

// DEBUG METHODS

// ListMessages lists the messages of the game
func (game *Game) ListMessages() {
	log.Println("Voici les messages du serveur")
	for _, msg := range game.Messages {
		log.Println("")
		log.Printf("Content : %v\n", msg.Content)
		log.Printf("Date : %v\n", msg.Date)
		log.Printf("Sender : %v\n", msg.Sender.Nickname)
		log.Println("")
	}
	log.Println("...")
}

// ListRooms lists the rooms of the game
func (game *Game) ListRooms() {
	log.Println("Voici les rooms du serveur")
	for i, room := range game.Rooms {
		log.Println("")
		log.Printf("Room %v : %v\n", i, room.Name)
		room.ListClients()
		log.Println("")
	}
	log.Println("...")
}

// ListClients prints the list of clients in the game
func (game *Game) ListClients() {
	log.Println("Voici les clients du serveur")
	for i, client := range game.Clients {
		log.Println("")
		log.Printf("Client %v : %v\n", i, client.Nickname)
		log.Println("")
	}
	log.Println("...")
}
