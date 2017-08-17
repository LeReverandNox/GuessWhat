package wss

import (
	"log"

	"github.com/LeReverandNox/GuessWhat/src/game"
	"golang.org/x/net/websocket"
)

func onConnection(ws *websocket.Conn) *game.Client {
	client := myGame.AddClient(ws)
	return client
}

func onDisconnection(client *game.Client, err error) error {
	log.Printf("Socket closed because of : %v", err)
	myGame.RemoveClient(client)
	return nil
}

func setNicknameAction(client *game.Client, nickname string) {
	client.SetNickname(nickname)

	msg := make(map[string]string)
	msg["message"] = "Je s'apelle " + client.Nickname
	client.Socket.SendToAll(myGame, msg)
}
