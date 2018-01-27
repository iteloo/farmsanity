package main

type GameState string

const (
	WaitingState    GameState = "waiting"
	ProductionState GameState = "production"
	AuctionState    GameState = "auction"
	TradeState      GameState = "trade"
)

type StateController interface {
	Name() GameState
	Begin()
	End()
}

type WaitingController struct {
	name GameState
}

func NewWaitingController(game *Game) *WaitingController {
	return &WaitingController{
		name: WaitingState,
	}
}
func (s *WaitingController) Name() GameState { return s.name }
func (s *WaitingController) Begin()          {}
func (s *WaitingController) End()            {}

type ProductionController struct {
	name GameState
}

func NewProductionController(game *Game) *ProductionController {
	return &ProductionController{
		name: ProductionState,
	}
}
func (s *ProductionController) Name() GameState { return s.name }
func (s *ProductionController) Begin()          {}
func (s *ProductionController) End()            {}

type AuctionController struct {
	name GameState
	game *Game
}

func NewAuctionController(game *Game) *AuctionController {
	return &AuctionController{
		name: AuctionState,
		game: game,
	}
}
func (s *AuctionController) Name() GameState { return s.name }
func (s *AuctionController) Begin()          {}
func (s *AuctionController) End()            {}

type TradeController struct {
	name GameState
}

func NewTradeController(game *Game) *TradeController {
	return &TradeController{
		name: TradeState,
	}
}
func (s *TradeController) Name() GameState { return s.name }
func (s *TradeController) Begin()          {}
func (s *TradeController) End()            {}

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
