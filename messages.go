package main

type MessageAction string

const (
	// Server broadcast messages
	GameStateChangedAction MessageAction = "game_state_changed"
	AuctionSeedAction      MessageAction = "auction_seed"

	// Server-to-client messages
	AuctionWonAction MessageAction = "auction_won"

	// Client messages
	BidAction          MessageAction = "bid"
	ReadyAction        MessageAction = "ready"
	JoinAction         MessageAction = "join"
	SaleAction         MessageAction = "join"
	TradeAction        MessageAction = "trade"
	ActivateCardAction MessageAction = "activate_card"
)

type Message interface{}

type MessageInterface struct {
	message interface{}
}

func NewMessage(message interface{}) Message {
	return MessageInterface{message}
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
