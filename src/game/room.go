package game

import (
	"errors"
	"log"
	"time"
)

type Room struct {
	Name     string
	Messages []*Message
	Clients  []*Client
	Drawer   *Client
	Owner    *Client
	Image    string
}

// NewRoom creates a new room and returns it
func NewRoom(name string, owner *Client) *Room {
	room := Room{}
	room.Name = name
	room.Clients = make([]*Client, 0)
	room.Messages = make([]*Message, 0)
	room.Owner = owner
	return &room
}

// RemoveClient removes a client from the room
func (room *Room) RemoveClient(clientToRemove *Client) (bool, error) {
	if room.IsClientIn(clientToRemove) {
		for i, client := range room.Clients {
			if client.Nickname == clientToRemove.Nickname {
				room.Clients = append(room.Clients[:i], room.Clients[i+1:]...)
				return room.IsEmpty(), nil
			}
		}
	}
	return room.IsEmpty(), errors.New("The client " + clientToRemove.Nickname + " is not in the room " + room.Name)
}

// AddClient adds a client to the room
func (room *Room) AddClient(client *Client) error {
	if !room.IsClientIn(client) {
		room.Clients = append(room.Clients, client)
		return nil
	}
	return errors.New("The client " + client.Nickname + " is already in the room " + room.Name)
}

// IsClientIn returns true if the given client is in.
func (room *Room) IsClientIn(clientToSearch *Client) bool {
	for _, client := range room.Clients {
		if clientToSearch == client {
			return true
		}
	}
	return false
}

// AddMessage adds a message to the game
func (room *Room) AddMessage(sender *Client, content string) *Message {
	message := NewMessage(sender, content, time.Now())
	room.Messages = append(room.Messages, message)

	return message
}

// IsEmpty returns true if the room is empty.
func (room *Room) IsEmpty() bool {
	if len(room.Clients) == 0 {
		return true
	}
	return false
}

// ListClients lists the clients of the room
func (room *Room) ListClients() {
	log.Println("Voici les clients de la channel")
	for i, client := range room.Clients {
		log.Println("")
		log.Printf("Client %v : %v\n", i, client.Nickname)
		log.Println("")
	}
	log.Println("...")
}
