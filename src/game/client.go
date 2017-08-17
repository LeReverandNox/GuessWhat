package game

import (
	"github.com/LeReverandNox/GuessWhat/src/socket"
	"golang.org/x/net/websocket"
)

type Client struct {
	Socket   *socket.Socket
	Nickname string
}

func NewClient(ws *websocket.Conn) *Client {
	client := Client{}
	client.Socket = &socket.Socket{ws}
	return &client
}

func (game *Game) AddClient(ws *websocket.Conn) *Client {
	client := NewClient(ws)
	game.Clients = append(game.Clients, client)
	return client

}
