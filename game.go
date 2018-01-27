package main

type GameConnection interface {
	broadcast(message Message) error
}

type Game struct {
	connection GameConnection
	state      StateController
}

func NewGame(connection GameConnection) *Game {
	game := Game{
		connection: connection,
		state:      nil,
	}
	game.state = NewStateController(&game, WaitingState)
	game.state.Begin()
	return &game
}

func (g *Game) ChangeState(newState GameState) {
	g.state.End()
	g.connection.broadcast(NewGameStateChangedMessage(newState))
	g.state = NewStateController(g, newState)
	g.state.Begin()
}
