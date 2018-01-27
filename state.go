package main

import (
	"math/rand"
)

type GameState string

const (
	WaitingState    GameState = "waiting"
	ProductionState GameState = "production"
	AuctionState    GameState = "auction"
	TradeState      GameState = "trade"
)

const (
	AuctionBidTime   int64 = 100
	TradeTimeout     int64 = 100
	TradingStageTime int64 = 1000
)

type StateController interface {
	Name() GameState
	Begin()
	End()
	Timer(tick int64)
	RecieveMessage(User, Message)
}

type WaitingController struct {
	name GameState
}

func NewWaitingController(game *Game) *WaitingController {
	return &WaitingController{
		name: WaitingState,
	}
}
func (s *WaitingController) Name() GameState                  { return s.name }
func (s *WaitingController) Begin()                           {}
func (s *WaitingController) End()                             {}
func (s *WaitingController) Timer(tick int64)                 {}
func (s *WaitingController) RecieveMessage(u User, m Message) {}

type ProductionController struct {
	name GameState
}

func NewProductionController(game *Game) *ProductionController {
	return &ProductionController{
		name: ProductionState,
	}
}
func (s *ProductionController) Name() GameState                  { return s.name }
func (s *ProductionController) Begin()                           {}
func (s *ProductionController) End()                             {}
func (s *ProductionController) Timer(tick int64)                 {}
func (s *ProductionController) RecieveMessage(u User, m Message) {}

type AuctionController struct {
	name   GameState
	game   *Game
	bid    int
	step   int
	steps  int
	winner User
}

func NewAuctionController(game *Game) *AuctionController {
	return &AuctionController{
		name:  AuctionState,
		game:  game,
		steps: 3,
	}
}
func (s *AuctionController) Name() GameState { return s.name }
func (s *AuctionController) Begin() {
	s.issueCard()
}

func (s *AuctionController) issueCard() {
	// When the auction begins, we need to choose a random number and broadcast
	// it to the participants.
	seed := rand.Int()
	s.game.connection.broadcast(
		NewAuctionSeedMessage(seed),
	)

	// Set a timeout.
	s.game.SetTimeout(AuctionBidTime)
}

func (s *AuctionController) End() {}

// The timer is only used to determine when the auction is over. So when we get
// this call, the current auction is over.
func (s *AuctionController) Timer(tick int64) {
	if s.winner != nil {
		s.winner.Message(NewAuctionWonMessage())
	}

	// Reset the bid and winner.
	s.bid = 0
	s.winner = nil

	s.step += 1

	if s.step == s.steps {
		// We've reached the end of the auction process. So let's change
		// phases. The TradeState is next.
		s.game.ChangeState(TradeState)
	} else {
		// Issue the next card.
		s.issueCard()
	}
}

func (s *AuctionController) RecieveMessage(u User, m Message) {
	switch msg := m.(type) {
	case BidMessage:
		if msg.Amount > s.bid {
			s.bid = msg.Amount
			s.winner = u
			s.game.SetTimeout(AuctionBidTime)
		}
	}
}

type TradeController struct {
	name            GameState
	game            *Game
	stagedMaterials string
	stagedUser      User
	stagingTime     int64
}

func NewTradeController(game *Game) *TradeController {
	return &TradeController{
		name: TradeState,
		game: game,
	}
}

func (s *TradeController) Name() GameState { return s.name }

func (s *TradeController) Begin() {
	// The trading stage ends after a certain time.
	s.game.SetTimeout(TradingStageTime)
}

func (s *TradeController) Timer(tick int64) {
	// The stage is over, so begin the next stage.
	s.game.ChangeState(ProductionState)
}

func (s *TradeController) End() {}
func (s *TradeController) RecieveMessage(u User, m Message) {
	switch msg := m.(type) {
	case TradeMessage:
		isntSelfTrade := s.stagedUser != u
		withinTimeInterval := s.game.GetTime()-s.stagingTime < TradeTimeout
		if isntSelfTrade && s.stagedUser != nil && withinTimeInterval {
			// Execute the currently proposed trade.
			s.stagedUser.Message(NewTradeCompletedMessage(msg.Materials))
			u.Message(NewTradeCompletedMessage(s.stagedMaterials))

			// Reset the staged materials
			s.stagedUser = nil
			s.stagingTime = 0
			s.stagedMaterials = ""
		} else {
			s.stagedUser = u
			s.stagingTime = s.game.GetTime()
			s.stagedMaterials = msg.Materials
		}
	}
}

func NewStateController(game *Game, state GameState) StateController {
	switch state {

	case WaitingState:
		return NewWaitingController(game)
	case ProductionState:
		return NewProductionController(game)
	case AuctionState:
		return NewAuctionController(game)
	case TradeState:
		return NewTradeController(game)
	default:
		panic("Unknown state!")
	}
}
