package game

import "time"
import "log"

type Message struct {
	Nickname string
	Content  string
	Date     time.Time
}

func NewMessage(nickname string, content string, date time.Time) *Message {
	message := Message{}
	message.Nickname = nickname
	message.Content = content
	message.Date = date

	return &message
}

func (game *Game) AddMessage(nickname string, content string) *Message {
	message := NewMessage(nickname, content, time.Now())
	game.Messages = append(game.Messages, message)

	return message
}

func (game *Game) ListMessages() {
	log.Print("Voici les messages du serveur")
	for _, msg := range game.Messages {
		log.Print(msg)
	}
}
