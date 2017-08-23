package game

type Client struct {
	Socket   *Socket
	Nickname string
	ID       int
	Score    int
}

var id int

func NewClient(socket *Socket, nickname string) *Client {
	defer func() {
		id++
	}()

	client := Client{}
	client.Socket = socket
	client.Nickname = nickname
	client.ID = id
	return &client
}

func (client *Client) SetNickname(nickname string) error {
	client.Nickname = nickname
	return nil
}

func (client *Client) AddToScore(amount int) {
	client.Score += amount
}

func (client *Client) ResetScore() {
	client.Score = 0
}
