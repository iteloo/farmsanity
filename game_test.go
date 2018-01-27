package main

import (
	"encoding/json"

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

func CompareBroadcastLog(got, want TestConnection) string {
	return cmp.Diff(got.broadcastLog, want.broadcastLog)
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
