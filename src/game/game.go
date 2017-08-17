package game

import "errors"

type Game struct {
	Clients []*Client
	Rooms   []*Room
}

func NewGame() *Game {
	game := Game{}
	game.Clients = make([]*Client, 0)
	game.Rooms = make([]*Room, 0)

	return &game
}

// AddRoom adds a room to the server
func (game *Game) AddRoom(name string) error {
	if !game.isRoomExisting(name) {
		room := NewRoom(name)
		game.Rooms = append(game.Rooms, room)
		return nil
	}
	return errors.New("The channel " + name + " already exist")
}

func (game *Game) isRoomExisting(name string) bool {
	for _, room := range game.Rooms {
		if room.Name == name {
			return true
		}
	}
	return false
}
