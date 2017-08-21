package game

import (
	"errors"
	"log"

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

// ListClients prints the list of clients in the game
func (game *Game) ListClients() {
	log.Printf("Voici les clients du serveur")
	for i, client := range game.Clients {
		log.Printf("Client %v : %v", i, client.Nickname)
	}
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

func (game *Game) isRoomExisting(name string) bool {
	for _, room := range game.Rooms {
		if room.Name == name {
			return true
		}
	}
	return false
}

// GetRoom returns the desired room. If not existing, creates it before.
func (game *Game) GetRoom(name string) *Room {
	if game.isRoomExisting(name) {
		for _, room := range game.Rooms {
			if room.Name == name {
				return room
			}
		}
	}
	room, _ := game.AddRoom(name)
	return room.(*Room)
}
