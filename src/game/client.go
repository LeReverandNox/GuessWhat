package game

import "golang.org/x/net/websocket"

type Client struct {
	Socket   *Socket
	Nickname string
}

func NewClient(ws *websocket.Conn) *Client {
	client := Client{}
	client.Socket = &Socket{ws}
	return &client
}

func (client *Client) SetNickname(nickname string) error {
	client.Nickname = nickname
	return nil
}
