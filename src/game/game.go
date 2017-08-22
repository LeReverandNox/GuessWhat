package game

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/LeReverandNox/GuessWhat/src/tools"
)

type Game struct {
	Clients  []*Client
	Rooms    []*Room
	Messages []*Message
	Words    []*Word
}

// NewGame creates a new Game struct, and returns it
func NewGame() *Game {
	game := Game{}
	game.Clients = make([]*Client, 0)
	game.Rooms = make([]*Room, 0)
	game.Messages = make([]*Message, 0)
	game.Words = make([]*Word, 0)

	game.loadWordsFromFile("./assets/words.txt")

	return &game
}

// AddClient adds a client to the game
func (game *Game) AddClient(socket *Socket, nickname string) *Client {
	client := NewClient(socket, nickname)
	game.Clients = append(game.Clients, client)
	return client

}

// RemoveClient remove a client from the game
func (game *Game) RemoveClient(clientToDelete *Client) {
	for i, client := range game.Clients {
		if client.Socket == clientToDelete.Socket {
			game.Clients = append(game.Clients[:i], game.Clients[i+1:]...)
		}
	}
}

// RemoveRoom remove a room from the game
func (game *Game) RemoveRoom(roomToDelete *Room) {
	for i, room := range game.Rooms {
		if room.Name == roomToDelete.Name {
			game.Rooms = append(game.Rooms[:i], game.Rooms[i+1:]...)
		}
	}
}

// AddRoom adds a room to the server
func (game *Game) AddRoom(name string, owner *Client) (interface{}, error) {
	if !game.IsRoomExisting(name) {
		room := NewRoom(name, owner)
		game.Rooms = append(game.Rooms, room)
		return room, nil
	}
	return nil, errors.New("The channel " + name + " already exist")
}

// GetRoom returns the desired room. If not existing, creates it before.
func (game *Game) GetRoom(name string, client *Client) (*Room, bool) {
	if game.IsRoomExisting(name) {
		for _, room := range game.Rooms {
			if room.Name == name {
				return room, false
			}
		}
	}
	room, _ := game.AddRoom(name, client)
	return room.(*Room), true
}

// GetCurrentClientRoom returns the current room of a client.
func (game *Game) GetCurrentClientRoom(client *Client) interface{} {
	for _, room := range game.Rooms {
		if room.IsClientIn(client) {
			return room
		}
	}
	return nil
}

// AddMessage adds a message to the game
func (game *Game) AddMessage(sender *Client, content string) *Message {
	message := NewMessage(sender, content, time.Now())
	game.Messages = append(game.Messages, message)

	return message
}

// IsNicknameTaken return true if the given nickname already exists in the game
func (game *Game) IsNicknameTaken(nicknameToTest string) bool {
	for _, client := range game.Clients {
		if client.Nickname == nicknameToTest {
			return true
		}
	}
	return false
}

// AddWord adds a word to the game
func (game *Game) AddWord(wordStr string) *Word {
	trimmedString := tools.RemoveAllSpaces(wordStr)
	loweredString := tools.ToLowerCase(trimmedString)
	if len(loweredString) < 0 {
		return nil
	}
	if game.isWordIn(loweredString) {
		return nil
	}
	if !tools.IsAlphaHyphen(loweredString) {
		return nil
	}
	word := NewWord(loweredString)
	game.Words = append(game.Words, word)

	return word
}

func (game *Game) PickRandomWord() *Word {
	return game.Words[rand.Intn(len(game.Words))]
}

func (game *Game) IsRoomExisting(name string) bool {
	for _, room := range game.Rooms {
		if room.Name == name {
			return true
		}
	}
	return false
}

func (game *Game) loadWordsFromFile(path string) {
	lines, err := tools.ReadLines(path)
	if err != nil {
		log.Fatalf("Error when loading Game words : %s\n", err)
	}
	for _, line := range lines {
		game.AddWord(line)
	}
	if len(game.Words) < 1 {
		log.Fatalf("Error when loading Game words : The file must contain at least one (unique) word.")
	}
}

func (game *Game) isWordIn(wordToSearch string) bool {
	for _, word := range game.Words {
		if word.Value == wordToSearch {
			return true
		}
	}
	return false
}

// DEBUG METHODS

// ListMessages lists the messages of the game
func (game *Game) ListMessages() {
	log.Println("Voici les messages du serveur")
	for _, msg := range game.Messages {
		log.Println("")
		log.Printf("Content : %v\n", msg.Content)
		log.Printf("Date : %v\n", msg.Date)
		log.Printf("Sender : %v\n", msg.Sender.Nickname)
		log.Println("")
	}
	log.Println("...")
}

// ListRooms lists the rooms of the game
func (game *Game) ListRooms() {
	log.Println("Voici les rooms du serveur")
	for i, room := range game.Rooms {
		log.Println("")
		log.Printf("Room %v : %v\n", i, room.Name)
		log.Printf("Owner %v : %v\n", i, room.Owner.Nickname)
		room.ListClients()
		log.Println("")
	}
	log.Println("...")
}

// ListClients prints the list of clients in the game
func (game *Game) ListClients() {
	log.Println("Voici les clients du serveur")
	for i, client := range game.Clients {
		log.Println("")
		log.Printf("Client %v : %v\n", i, client.Nickname)
		log.Println("")
	}
	log.Println("...")
}

// ListWords prints the list of words in the game
func (game *Game) ListWords() {
	log.Println("Voici les words du serveur")
	for i, word := range game.Words {
		log.Println("")
		log.Printf("Word %v : %v\n", i, word.Value)
		log.Printf("Length of Word %v : %v\n", i, word.Length)
		log.Println("")
	}
	log.Println("...")
}
