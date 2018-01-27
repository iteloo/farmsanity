package main

import (
	"math/rand"
	"time"
)

type GameState string

const (
	WaitingState    GameState = "waiting"
	ProductionState GameState = "production"
	AuctionState    GameState = "auction"
	TradeState      GameState = "trade"
)

const (
	AuctionBidTime   time.Duration = 30 * time.Second
	TradeTimeout     time.Duration = 100 * time.Millisecond
	TradingStageTime time.Duration = 120 * time.Second

	MinPlayers int = 2
)

type StateController interface {
	Name() GameState
	Begin()
	End()
	Timer(tick time.Duration)
	RecieveMessage(User, Message)
}

type WaitingController struct {
	game  *Game
	name  GameState
	ready map[User]bool
}

func NewWaitingController(game *Game) *WaitingController {
	return &WaitingController{
		game:  game,
		name:  WaitingState,
		ready: map[User]bool{},
	}
}
func (s *WaitingController) Name() GameState          { return s.name }
func (s *WaitingController) Begin()                   {}
func (s *WaitingController) End()                     {}
func (s *WaitingController) Timer(tick time.Duration) {}
func (s *WaitingController) RecieveMessage(u User, m Message) {
	switch m.(type) {
	case ReadyMessage:
		s.ready[u] = true
		s.proceedIfReady()
	case JoinMessage:
		s.ready[u] = false
	case LeaveMessage:
		delete(s.ready, u)
		s.proceedIfReady()
	}
}

func (s *WaitingController) proceedIfReady() {
	count := 0
	for _, v := range s.ready {
		if !v {
			return
		}
		count += 1
	}

	if count >= MinPlayers {
		s.game.ChangeState(ProductionState)
	}
}

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
func (s *ProductionController) Timer(tick time.Duration)         {}
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
	s.game.connection.Broadcast(
		NewAuctionSeedMessage(seed),
	)

	// Set a timeout.
	s.game.SetTimeout(AuctionBidTime)
}

func (s *AuctionController) End() {}

// The timer is only used to determine when the auction is over. So when we get
// this call, the current auction is over.
func (s *AuctionController) Timer(tick time.Duration) {
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
	stagingTime     time.Duration
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

func (s *TradeController) Timer(tick time.Duration) {
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
