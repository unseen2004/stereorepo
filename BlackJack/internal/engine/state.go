package engine

import "BlackJack/internal/rules"

type EventType string

const (
	EventRoundStart  EventType = "round_start"
	EventInitialDeal EventType = "initial_deal"
	EventBetPlaced   EventType = "bet_placed"
	EventAction      EventType = "action"
	EventDealerPlay  EventType = "dealer_play"
	EventRoundEnd    EventType = "round_end"
)

type Event struct {
	Type    EventType
	Player  string
	Hand    int
	Action  Action
	Cards   []Card
	Results []RoundResult
	Message string
}

type GameState struct {
	Rules         rules.Rules
	DealerCards   []Card
	Players       []PlayerState
	ShoeRemaining int
	RoundOver     bool
}

type PlayerState struct {
	ID       string
	Name     string
	Bankroll int
	Hands    []HandState
}

type HandState struct {
	Cards       []Card
	Bet         int
	OriginalBet int
	IsSplitHand bool
	IsDoubled   bool
	IsFinished  bool
	IsSurrender bool
}
