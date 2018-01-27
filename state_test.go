package main

import "testing"

func TestAuctionBidding(t *testing.T) {
	connection := TestConnection{}
	game := NewGame(&connection)
	ctrl := NewAuctionController(game)

	u1 := &TestUser{}
	u2 := &TestUser{}
	ctrl.RecieveMessage(u1, NewBidMessage(10))
	ctrl.RecieveMessage(u2, NewBidMessage(5))

	if ctrl.bid != 10 {
		t.Errorf("Expected ctrl.bid = 10, got %v", ctrl.bid)
	}
	if ctrl.winner != u1 {
		t.Errorf("Expected ctrl.winner = u, got %v", ctrl.winner)
	}

	// Can't win by bidding the same amount.
	ctrl.RecieveMessage(u2, NewBidMessage(10))
	if ctrl.winner != u1 {
		t.Errorf("Expected ctrl.winner = u, got %v", ctrl.winner)
	}

	// Outbidding will switch winner.
	ctrl.RecieveMessage(u2, NewBidMessage(12))
	if ctrl.winner != u2 {
		t.Errorf("Expected ctrl.winner = u, got %v", ctrl.winner)
	}
}

func TestAuctionTimeout(t *testing.T) {
	connection := TestConnection{}
	game := NewGame(&connection)
	ctrl := NewAuctionController(game)
	game.state = ctrl

	user := &TestUser{}
	loser := &TestUser{}
	ctrl.RecieveMessage(user, NewBidMessage(10))
	ctrl.RecieveMessage(loser, NewBidMessage(5))

	// Wait for the auction to end.
	game.Tick(2 * AuctionBidTime)

	// Expect the winner to get a winning message.
	want := &TestUser{}
	want.Message(NewAuctionWonMessage())

	if diff := CompareMessageLog(user, want); diff != "" {
		t.Errorf("AuctionWonMessage: %q, %q, diff: %v", user.messageLog, want.messageLog, diff)
	}
}
