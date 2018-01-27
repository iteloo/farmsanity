package main

import (
	"encoding/json"
	"fmt"
)

type MessageAction string

const (
	// Server broadcast messages
	GameStateChangedAction MessageAction = "game_state_changed"
	AuctionSeedAction      MessageAction = "auction_seed"

	// Server-to-client messages
	AuctionWonAction     MessageAction = "auction_won"
	TradeCompletedAction MessageAction = "trade_completed"

	// Client messages
	BidAction   MessageAction = "bid"
	ReadyAction MessageAction = "ready"
	JoinAction  MessageAction = "join"
	LeaveAction MessageAction = "leave"
	TradeAction MessageAction = "trade"
)

type Message interface{}

type MessageInterface struct {
	message interface{}
}

type BasicMessage struct {
	Action string `json:"action"`
}

// Messages broadcast by the server.

type GameStateChangedMessage struct {
	Action   string `json:"action"`
	NewState string `json:"new_state"`
}

func NewGameStateChangedMessage(newState GameState) Message {
	return GameStateChangedMessage{
		Action:   string(GameStateChangedAction),
		NewState: string(newState),
	}
}

type AuctionSeedMessage struct {
	Action string `json:"action"`
	Seed   int    `json:"seed"`
}

func NewAuctionSeedMessage(seed int) Message {
	return AuctionSeedMessage{
		Action: string(AuctionSeedAction),
		Seed:   seed,
	}
}

type TradeCompletedMessage struct {
	Action    string `json:"action"`
	Materials string `json:"materials"`
}

func NewTradeCompletedMessage(materials string) Message {
	return TradeCompletedMessage{string(TradeCompletedAction), materials}
}

// Client messages

type BidMessage struct {
	Action string `json:"action"`
	Amount int    `json:"amount"`
}

func NewBidMessage(amount int) Message {
	return BidMessage{
		Action: string(BidAction),
		Amount: amount,
	}
}

type AuctionWonMessage struct {
	Action string `json:"action"`
}

func NewAuctionWonMessage() Message {
	return AuctionWonMessage{string(AuctionWonAction)}
}

type ReadyMessage struct {
	Action string `json:"action"`
}

func NewReadyMessage() Message {
	return ReadyMessage{string(ReadyAction)}
}

type JoinMessage struct {
	Action string `json:"action"`
}

func NewJoinMessage() Message {
	return JoinMessage{string(JoinAction)}
}

type LeaveMessage struct {
	Action string `json:"action"`
}

func NewLeaveMessage() Message {
	return LeaveMessage{string(LeaveAction)}
}

type TradeMessage struct {
	Action    string `json:"action"`
	Materials string `json:"materials"`
}

func NewTradeMessage(materials string) Message {
	return TradeMessage{
		Action:    string(TradeAction),
		Materials: materials,
	}
}

func DecodeMessage(data []byte) (Message, error) {
	msg := BasicMessage{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("Unable to decode message: %q", data)
	}

	// Now that we know the type of the message (based on the action)
	// we can decode it properly.
	var message Message
	var err error
	switch msg.Action {
	case string(GameStateChangedAction):
		m := GameStateChangedMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	case string(AuctionSeedAction):
		m := AuctionSeedMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	case string(AuctionWonAction):
		m := AuctionWonMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	case string(TradeCompletedAction):
		m := TradeCompletedMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	case string(BidAction):
		m := BidMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	case string(ReadyAction):
		m := ReadyMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	case string(JoinAction):
		m := JoinMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	case string(LeaveAction):
		m := LeaveMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	case string(TradeAction):
		m := TradeMessage{}
		err = json.Unmarshal(data, &m)
		message = m
	default:
		err = fmt.Errorf("Unknown action: %v", msg.Action)
	}

	return message, err
}
