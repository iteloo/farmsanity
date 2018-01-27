package main

type GameConnection interface {
	broadcast(message Message) error
}

type GameState string

const (
	WaitingState    GameState = "waiting"
	ProductionState GameState = "production"
	AuctionState    GameState = "auction"
	TradeState      GameState = "trade"
)

type Game struct {
	connection GameConnection
	state      GameState
}

func NewGame(connection GameConnection) Game {
	return Game{
		connection: connection,
		state:      WaitingState,
	}
}

func (g *Game) ChangeState(newState GameState) {
	g.state = newState
	g.connection.broadcast(NewGameStateChangedMessage(newState))
}
