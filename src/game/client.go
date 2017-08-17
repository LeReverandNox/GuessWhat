package game

import "github.com/LeReverandNox/GuessWhat/src/socket"

type Client struct {
	Socket   *socket.Socket
	Nickname string
}
