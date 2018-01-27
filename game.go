package main

type User interface {
	Message(message Message) error
}

type GameConnection interface {
	broadcast(message Message) error
}

type Game struct {
	connection  GameConnection
	state       StateController
	nextTimeout int64
	tick        int64
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

func (g *Game) SetTimeout(duration int64) {
	g.nextTimeout = g.tick + duration
}

func (g *Game) GetTime() int64 {
	return g.tick
}

func (g *Game) Tick(time int64) {
	g.tick = time

	// If a timer is currently set, notify the state controller.
	if g.nextTimeout != 0 && time > g.nextTimeout {
		g.nextTimeout = 0
		g.state.Timer(time)
	}
}

func (g *Game) RecieveMessage(user User, message Message) {
	g.state.RecieveMessage(user, message)
}

func (g *Game) ChangeState(newState GameState) {
	g.state.End()

	// Clean up any timers that are currently running
	g.nextTimeout = 0

	g.connection.broadcast(NewGameStateChangedMessage(newState))
	g.state = NewStateController(g, newState)
	g.state.Begin()
}
