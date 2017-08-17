package game

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
