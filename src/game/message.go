package game

import "time"

type Message struct {
	Sender  *Client
	Content string
	Date    time.Time
}

func NewMessage(sender *Client, content string, date time.Time) *Message {
	message := Message{}
	message.Sender = sender
	message.Content = content
	message.Date = date

	return &message
}
