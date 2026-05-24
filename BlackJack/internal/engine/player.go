package engine

type Player struct {
	ID       string
	Name     string
	Bankroll int
	Hands    []*Hand
	Provider DecisionProvider
}

type DecisionProvider interface {
	BetRequest(state GameState, minBet, maxBet int) (int, error)
	InsuranceRequest(state GameState) (bool, error)
	ActionRequest(state GameState, handIndex int, actions []Action) (Action, error)
	Notify(state GameState, event Event)
}
