package game

type Winner struct {
	Client  *Client
	WinTime int
}

func NewWinner(client *Client, winTime int) *Winner {
	winner := &Winner{}
	winner.Client = client
	winner.WinTime = winTime

	return winner
}
