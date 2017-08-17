package game

import "errors"

type Room struct {
	Name    string
	Clients []*Client
	Drawer  string
	Image   string
}

func NewRoom(name string) *Room {
	room := Room{}
	room.Name = name
	room.Clients = make([]*Client, 0)

	return &room
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
