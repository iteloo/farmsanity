package main

import (
	"encoding/json"
	"math/rand"

	"github.com/google/go-cmp/cmp"

	"testing"
)

type TestConnection struct {
	broadcastLog []string
}

func (c *TestConnection) broadcast(message Message) error {
	result, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	c.broadcastLog = append(c.broadcastLog, string(result))
	return nil
}

type TestUser struct {
	name       string
	messageLog []string
}

func (u *TestUser) Message(message Message) error {
	result, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	u.messageLog = append(u.messageLog, string(result))
	return nil
}

func CompareBroadcastLog(got, want TestConnection) string {
	return cmp.Diff(got.broadcastLog, want.broadcastLog)
}

func CompareMessageLog(got, want *TestUser) string {
	return cmp.Diff(got.messageLog, want.messageLog)
}

func TestChangeState(t *testing.T) {
	connection := TestConnection{}
	game := NewGame(&connection)
	game.ChangeState(TradeState)

	expected := TestConnection{}
	expected.broadcast(NewGameStateChangedMessage(TradeState))

	if diff := CompareBroadcastLog(connection, expected); diff != "" {
		t.Errorf("ChangeState(WaitingState): %v", diff)
	}

	if game.state.Name() != TradeState {
		t.Errorf("game.state = %v, want %v", game.state, TradeState)
	}
}

func TestAuctionStart(t *testing.T) {
	// Need to set the random seed to force deterministic behavior.
	rand.Seed(1)
	connection := TestConnection{}
	game := NewGame(&connection)
	game.ChangeState(AuctionState)

	rand.Seed(1)
	expected := TestConnection{}
	expected.broadcast(NewGameStateChangedMessage(AuctionState))
	expected.broadcast(NewAuctionSeedMessage(rand.Int()))

	if diff := CompareBroadcastLog(connection, expected); diff != "" {
		t.Errorf("ChangeState(WaitingState): %v", diff)
	}

	if game.state.Name() != AuctionState {
		t.Errorf("game.state.Name() = %v, want %v", game.state.Name, AuctionState)
	}
}

func TestAuctionPhases(t *testing.T) {
	// Need to set the random seed to force deterministic behavior.
	rand.Seed(1)
	connection := TestConnection{}
	game := NewGame(&connection)
	game.ChangeState(AuctionState)

	// Bid on a card.
	user := &TestUser{}
	game.RecieveMessage(user, NewBidMessage(10))

	// Wait until the player wins.
	game.Tick(AuctionBidTime + 1)

	// Wait until the second auction expires with no bids.
	game.Tick(2*AuctionBidTime + 2)

	// Wait until the third auction expires with no bids.
	game.Tick(3*AuctionBidTime + 3)

	rand.Seed(1)
	expected := TestConnection{}
	expected.broadcast(NewGameStateChangedMessage(AuctionState))
	expected.broadcast(NewAuctionSeedMessage(rand.Int()))
	expected.broadcast(NewAuctionSeedMessage(rand.Int()))
	expected.broadcast(NewAuctionSeedMessage(rand.Int()))
	expected.broadcast(NewGameStateChangedMessage(TradeState))

	if diff := CompareBroadcastLog(connection, expected); diff != "" {
		t.Errorf("ChangeState(WaitingState): %v", diff)
	}
}
