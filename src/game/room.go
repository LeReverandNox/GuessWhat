package game

import (
	"errors"
	"log"
)

type Room struct {
	Name    string
	Clients []*Client
	Drawer  string
	Image   string
}

// NewRoom creates a new room and returns it
func NewRoom(name string) *Room {
	room := Room{}
	room.Name = name
	room.Clients = make([]*Client, 0)

	return &room
}

// RemoveClient removes a client from the room
func (room *Room) RemoveClient(client *Client) error {
	for i, client := range room.Clients {
		if room.IsClientIn(client) {
			room.Clients = append(room.Clients[:i], room.Clients[i+1:]...)
			return nil
		}
	}
	return errors.New("The client " + client.Nickname + " is already in the room " + room.Name)
}

// AddClient adds a client to the room
func (room *Room) AddClient(client *Client) error {
	if !room.IsClientIn(client) {
		room.Clients = append(room.Clients, client)
		return nil
	}
	return errors.New("The client " + client.Nickname + " is already in the room " + room.Name)
}

func (room *Room) IsClientIn(clientToSearch *Client) bool {
	for _, client := range room.Clients {
		if clientToSearch == client {
			return true
		}
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
