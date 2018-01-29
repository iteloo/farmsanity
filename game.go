package main

import (
	"time"
)

type User interface {
	Message(message Message) error
}

type GameConnection interface {
	Broadcast(message Message) error
}

type Game struct {
	name        string
	connection  GameConnection
	state       StateController
	nextTimeout time.Duration
	tick        time.Duration
}

func NewGame(name string, connection GameConnection) *Game {
	game := Game{
		name:       name,
		connection: connection,
		state:      nil,
	}
	game.state = NewStateController(&game, WaitingState)
	game.state.Begin()
	return &game
}

func (g *Game) SetTimeout(duration time.Duration) {
	g.nextTimeout = g.tick + duration
}

func (g *Game) GetTime() time.Duration {
	return g.tick
}

func (g *Game) Tick(time time.Duration) {
	g.tick = time

	// If a timer is currently set, notify the state controller.
	if g.nextTimeout != 0 && time > g.nextTimeout {
		g.nextTimeout = 0
		g.state.Timer(time)
	}
}

func (g *Game) RecieveMessage(user User, message Message) {
	switch message.(type) {
	case JoinMessage:
		user.Message(NewWelcomeMessage(g.name, string(g.state.Name())))
	}
	g.state.RecieveMessage(user, message)
}

func (g *Game) ChangeState(newState GameState) {
	g.state.End()

	// Clean up any timers that are currently running
	g.nextTimeout = 0

	g.connection.Broadcast(NewGameStateChangedMessage(newState))
	g.state = NewStateController(g, newState)
	g.state.Begin()
}
