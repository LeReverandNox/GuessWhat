package game

import (
	"log"

	"golang.org/x/net/websocket"
)

type Client struct {
	Socket   *Socket
	Nickname string
}

func NewClient(ws *websocket.Conn) *Client {
	client := Client{}
	client.Socket = &Socket{ws}
	return &client
}

func (game *Game) AddClient(ws *websocket.Conn) *Client {
	client := NewClient(ws)
	game.Clients = append(game.Clients, client)
	return client

}

func (client *Client) SetNickname(nickname string) error {
	client.Nickname = nickname
	return nil
}

func (game *Game) ListClients() {
	log.Printf("Voici les clients du serveur")
	for i, client := range game.Clients {
		log.Printf("Client %v : %v", i, client.Nickname)
	}
}

func (game *Game) RemoveClient(clientToDelete *Client) {
	for i, client := range game.Clients {
		if client.Socket == clientToDelete.Socket {
			game.Clients = append(game.Clients[:i], game.Clients[i+1:]...)
		}
	}
}
