package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	TickInterval time.Duration = 300 * time.Millisecond
)

type Player struct {
	Name       string
	Connection *websocket.Conn
}

func (p *Player) Message(message Message) error {
	return p.Connection.WriteJSON(message)
}

func GenerateGameName() string {
	return "test"
}

type Event struct {
	Message Message
	Player  *Player
}

func NewEvent(player *Player, message Message) Event {
	return Event{
		Player:  player,
		Message: message,
	}
}

type GameServer struct {
	players          []Player
	Name             string
	game             *Game
	incomingMessages chan Event
}

func (s *GameServer) Broadcast(message Message) error {
	log.Printf("Broadcast: %v", message)
	for _, p := range s.players {
		err := p.Message(message)
		if err != nil {
			log.Printf("Write failed during broadcast: %v\n", err)
		}
	}
	return nil
}

type PlayerMessage struct {
	Player Player
}

type TickMessage struct {
	Tick time.Duration
}

func (s *GameServer) AddPlayer(player Player) {
	s.incomingMessages <- NewEvent(&player, NewJoinMessage())
}

func (s *GameServer) HandleCommunication(player Player) {
	// Send a join message as we arrive.
	s.incomingMessages <- NewEvent(&player, NewJoinMessage)

	for {
		t, data, err := player.Connection.ReadMessage()
		if err != nil {
			log.Printf("Websocket[name=%v] read error: %v", player.Name, err)
			return
		}

		if t != websocket.TextMessage {
			log.Printf("Websocket[name=%v] sent binary message", player.Name)
		}

		msg, err := DecodeMessage(data)
		if err != nil {
			log.Printf("Websocket[name=%v] sent invalid message: %v", player.Name, err)
		}
		s.incomingMessages <- NewEvent(&player, msg)
	}
}

func (s *GameServer) HandleMessages() {
	for {
		event := <-s.incomingMessages
		switch msg := event.Message.(type) {
		case TickMessage:
			s.game.Tick(msg.Tick)
		case JoinMessage:
			new := true
			for _, x := range s.players {
				if *event.Player == x {
					new = false
					break
				}
			}

			if new {
				// On the first pass, set up the player and begin handling
				// their messages for them.
				s.players = append(s.players, *event.Player)
				go s.HandleCommunication(*event.Player)
			} else {
				// On subsequent passes, we just want to send the message
				// through to the game controller.
				s.game.RecieveMessage(event.Player, event.Message)
			}
		default:
			s.game.RecieveMessage(event.Player, event.Message)
		}
	}
}

func (s *GameServer) RunClock() {
	ticks := 0 * time.Second
	for {
		time.Sleep(TickInterval)
		ticks += TickInterval
		s.incomingMessages <- NewEvent(nil, TickMessage{ticks})
	}
}

func NewGameServer(name string) *GameServer {
	g := GameServer{
		Name:             name,
		game:             nil,
		incomingMessages: make(chan Event),
	}
	g.game = NewGame(&g)

	go g.HandleMessages()
	go g.RunClock()

	return &g
}
