package game

import "golang.org/x/net/websocket"

type Client struct {
	Socket   *Socket
	Nickname string
	ID       int
}

var id int

func NewClient(ws *websocket.Conn) *Client {
	defer func() {
		id++
	}()

	client := Client{}
	client.Socket = &Socket{ws}
	client.ID = id
	return &client
}

func (client *Client) SetNickname(nickname string) error {
	client.Nickname = nickname
	return nil
}
