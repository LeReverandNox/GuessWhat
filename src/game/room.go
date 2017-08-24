package game

import (
	"errors"
	"log"
	"time"

	"github.com/LeReverandNox/GuessWhat/src/tools"
)

type Room struct {
	Name           string
	Messages       []*Message
	Clients        []*Client
	NeedingDrawing []*Client
	Drawer         *Client
	Owner          *Client
	Image          string
	Word           *Word
	IsStarted      bool
	TotalRounds    int
	ActualRound    int
	IsRoundGoing   bool
	RoundDuration  int
	roundEnd       time.Time
	roundTicker    *time.Ticker
}

// NewRoom creates a new room and returns it
func NewRoom(name string, owner *Client) *Room {
	room := Room{}
	room.Name = name
	room.Clients = make([]*Client, 0)
	room.Messages = make([]*Message, 0)
	room.Owner = owner
	room.IsStarted = false
	room.TotalRounds = 10
	room.ActualRound = 0
	room.IsRoundGoing = false
	room.RoundDuration = 90
	return &room
}

// RemoveClient removes a client from the room
func (room *Room) RemoveClient(clientToRemove *Client) (bool, error) {
	if room.IsClientIn(clientToRemove) {
		for i, client := range room.Clients {
			if client.Nickname == clientToRemove.Nickname {
				room.Clients = append(room.Clients[:i], room.Clients[i+1:]...)
				return room.IsEmpty(), nil
			}
		}
	}
	return room.IsEmpty(), errors.New("The client " + clientToRemove.Nickname + " is not in the room " + room.Name)
}

// AddClient adds a client to the room
func (room *Room) AddClient(client *Client) error {
	if !room.IsClientIn(client) {
		room.Clients = append(room.Clients, client)
		return nil
	}
	return errors.New("The client " + client.Nickname + " is already in the room " + room.Name)
}

// IsClientIn returns true if the given client is in.
func (room *Room) IsClientIn(clientToSearch *Client) bool {
	for _, client := range room.Clients {
		if clientToSearch == client {
			return true
		}
	}
	return false
}

// AddMessage adds a message to the game
func (room *Room) AddMessage(sender *Client, content string) *Message {
	message := NewMessage(sender, content, time.Now())
	room.Messages = append(room.Messages, message)

	return message
}

// IsEmpty returns true if the room is empty.
func (room *Room) IsEmpty() bool {
	if len(room.Clients) == 0 {
		return true
	}
	return false
}

func (room *Room) IsOwner(client *Client) bool {
	if client.Nickname == room.Owner.Nickname {
		return true
	}
	return false
}

func (room *Room) Start() {
	room.IsStarted = true
}

func (room *Room) GetNbClients() int {
	return len(room.Clients)
}

func (room *Room) SetWord(word *Word) {
	room.Word = word
}

func (room *Room) SetDrawer(drawer *Client) {
	room.Drawer = drawer
}

func (room *Room) PickRandomClient() *Client {
	return room.Clients[tools.RandomInt(len(room.Clients))]
}

func (room *Room) SetImage(base64 string) {
	room.Image = base64
}

func (room *Room) AddDrawingNeeder(client *Client) {
	room.NeedingDrawing = append(room.NeedingDrawing, client)
}

func (room *Room) CleanDrawingNeeders() {
	room.NeedingDrawing = room.NeedingDrawing[:0]
}

func (room *Room) ResetImage() {
	room.Image = ""
}

func (room *Room) IncrementRound() {
	room.ActualRound++
}

func (room *Room) ResetRounds() {
	room.ActualRound = 0
}

func (room *Room) ResetMessages() {
	room.Messages = room.Messages[:0]
}

func (room *Room) StartRound() {
	room.IsRoundGoing = true
}

func (room *Room) StopRound() {
	room.IsRoundGoing = false
}

func (room *Room) SetTicker(ticker *time.Ticker) *time.Ticker {
	room.roundTicker = ticker
	return room.roundTicker
}

func (room *Room) GetTicker() *time.Ticker {
	return room.roundTicker
}

func (room *Room) StopTicker() *time.Ticker {
	room.roundTicker.Stop()
	return room.roundTicker
}

func (room *Room) SetRoundEnd(time time.Time) time.Time {
	room.roundEnd = time
	return room.roundEnd
}

func (room *Room) GetRoundEnd() time.Time {
	return room.roundEnd
}

func (room *Room) SetRoundDuration(duration int) int {
	room.RoundDuration = duration
	return room.RoundDuration
}

func (room *Room) GetRoundDuration() int {
	return room.RoundDuration
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
