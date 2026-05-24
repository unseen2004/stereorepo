package bot

import "BlackJack/internal/engine"

type BasicBot struct {
	Name string
}

func (b *BasicBot) BetRequest(state engine.GameState, minBet, maxBet int) (int, error) {
	return minBet, nil
}

func (b *BasicBot) InsuranceRequest(state engine.GameState) (bool, error) {
	return false, nil
}

func (b *BasicBot) ActionRequest(state engine.GameState, handIndex int, actions []engine.Action) (engine.Action, error) {
	player := findPlayerState(state, b.Name)
	if player == nil || handIndex >= len(player.Hands) {
		return engine.ActionStand, nil
	}
	dealerUp := dealerUpcard(state)
	hand := player.Hands[handIndex]
	rules := strategyRules{DealerHitsSoft17: state.Rules.DealerHitsSoft17, AllowDoubleAfterSplit: state.Rules.AllowDoubleAfterSplit}
	return decideBasicAction(hand, dealerUp, rules, actions), nil
}

func (b *BasicBot) Notify(state engine.GameState, event engine.Event) {}

func findHand(state engine.GameState, name string, handIndex int) *engine.HandState {
	for _, p := range state.Players {
		if p.Name == name {
			if handIndex >= 0 && handIndex < len(p.Hands) {
				return &p.Hands[handIndex]
			}
		}
	}
	return nil
}
