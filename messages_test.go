package main

import "testing"

func TestMessageDecoding(t *testing.T) {
	data := []byte(`{"action": "trade", "materials": "a hammer"}`)
	msg, err := DecodeMessage(data)
	if err != nil {
		t.Errorf("DecodeMessage(...) returned err: %v", err)
	}
	trade := msg.(TradeMessage)

	want := "a hammer"
	if trade.Materials != want {
		t.Errorf("trade.Materials = %q, want %q", trade.Materials, want)
	}
}

func TestBidMessageDecoding(t *testing.T) {
	data := []byte(`{"action": "bid", "amount": 12}`)
	msg, err := DecodeMessage(data)
	if err != nil {
		t.Errorf("DecodeMessage(...) returned err: %v", err)
	}

	bid := msg.(BidMessage)

	want := 12
	if bid.Amount != want {
		t.Errorf("bid.Amount = %q, want %q", bid.Amount, want)
	}
}
