package bot

import (
	"strings"

	"BlackJack/internal/engine"
)

type BasicStrategyBot struct {
	Name string
}

func (b *BasicStrategyBot) BetRequest(state engine.GameState, minBet, maxBet int) (int, error) {
	return minBet, nil
}

func (b *BasicStrategyBot) InsuranceRequest(state engine.GameState) (bool, error) {
	return false, nil
}

func (b *BasicStrategyBot) ActionRequest(state engine.GameState, handIndex int, actions []engine.Action) (engine.Action, error) {
	player := findPlayerState(state, b.Name)
	if player == nil || handIndex >= len(player.Hands) {
		return engine.ActionStand, nil
	}
	dealerUp := dealerUpcard(state)
	hand := player.Hands[handIndex]
	rules := strategyRules{DealerHitsSoft17: state.Rules.DealerHitsSoft17, AllowDoubleAfterSplit: state.Rules.AllowDoubleAfterSplit}
	return decideBasicAction(hand, dealerUp, rules, actions), nil
}

func (b *BasicStrategyBot) Notify(state engine.GameState, event engine.Event) {}

func findPlayerState(state engine.GameState, name string) *engine.PlayerState {
	for i := range state.Players {
		if state.Players[i].Name == name {
			return &state.Players[i]
		}
	}
	return nil
}

func dealerUpcard(state engine.GameState) engine.Card {
	if len(state.DealerCards) == 0 {
		return engine.Card{Rank: engine.RankTen}
	}
	return state.DealerCards[0]
}

func decideBasicAction(hand engine.HandState, dealerUp engine.Card, rules strategyRules, actions []engine.Action) engine.Action {
	dealer := upcardValue(dealerUp)
	if canSurrender(actions) {
		hardTotal, soft := handTotal(hand.Cards)
		if !soft && (hardTotal == 16 && dealer >= 9 || hardTotal == 15 && dealer == 10) {
			return engine.ActionSurrender
		}
	}

	if canSplit(actions) && isPair(hand.Cards) {
		pairRank := hand.Cards[0].Rank
		if action := pairStrategy(pairRank, dealer, rules); action == engine.ActionSplit {
			return engine.ActionSplit
		}
	}

	total, soft := handTotal(hand.Cards)
	if soft {
		return chooseAction(softStrategy(total, dealer, rules), actions, engine.ActionHit)
	}
	return chooseAction(hardStrategy(total, dealer, rules), actions, engine.ActionHit)
}

type strategyRules struct {
	DealerHitsSoft17      bool
	AllowDoubleAfterSplit bool
}

func upcardValue(card engine.Card) int {
	switch card.Rank {
	case engine.RankAce:
		return 11
	case engine.RankTen, engine.RankJack, engine.RankQueen, engine.RankKing:
		return 10
	default:
		return int(card.Rank)
	}
}

func handTotal(cards []engine.Card) (int, bool) {
	total := 0
	aces := 0
	for _, c := range cards {
		if c.Rank == engine.RankAce {
			aces++
		}
		total += c.Value()
	}
	if aces > 0 && total+10 <= 21 {
		return total + 10, true
	}
	return total, false
}

func isPair(cards []engine.Card) bool {
	return len(cards) == 2 && cards[0].Rank == cards[1].Rank
}

func pairStrategy(rank engine.Rank, dealer int, rules strategyRules) engine.Action {
	switch rank {
	case engine.RankAce, engine.RankEight:
		return engine.ActionSplit
	case engine.RankTwo, engine.RankThree:
		if dealer >= 2 && dealer <= 7 {
			return engine.ActionSplit
		}
	case engine.RankFour:
		if rules.AllowDoubleAfterSplit && dealer >= 5 && dealer <= 6 {
			return engine.ActionSplit
		}
	case engine.RankFive:
		return engine.ActionHit
	case engine.RankSix:
		if dealer >= 2 && dealer <= 6 {
			return engine.ActionSplit
		}
	case engine.RankSeven:
		if dealer >= 2 && dealer <= 7 {
			return engine.ActionSplit
		}
	case engine.RankNine:
		if (dealer >= 2 && dealer <= 6) || dealer == 8 || dealer == 9 {
			return engine.ActionSplit
		}
	case engine.RankTen:
		return engine.ActionStand
	}
	return engine.ActionHit
}

func softStrategy(total int, dealer int, rules strategyRules) engine.Action {
	switch total {
	case 13, 14:
		if dealer >= 5 && dealer <= 6 {
			return engine.ActionDouble
		}
	case 15, 16:
		if dealer >= 4 && dealer <= 6 {
			return engine.ActionDouble
		}
	case 17:
		if dealer >= 3 && dealer <= 6 {
			return engine.ActionDouble
		}
	case 18:
		if rules.DealerHitsSoft17 {
			if dealer >= 2 && dealer <= 6 {
				return engine.ActionDouble
			}
		} else if dealer >= 3 && dealer <= 6 {
			return engine.ActionDouble
		}
		if dealer == 2 || dealer == 7 || dealer == 8 {
			return engine.ActionStand
		}
		return engine.ActionHit
	default:
		if total >= 19 {
			return engine.ActionStand
		}
	}
	return engine.ActionHit
}

func hardStrategy(total int, dealer int, rules strategyRules) engine.Action {
	switch {
	case total <= 8:
		return engine.ActionHit
	case total == 9:
		if rules.DealerHitsSoft17 {
			if dealer >= 3 && dealer <= 6 {
				return engine.ActionDouble
			}
		} else if dealer >= 2 && dealer <= 6 {
			return engine.ActionDouble
		}
		return engine.ActionHit
	case total == 10:
		if dealer >= 2 && dealer <= 9 {
			return engine.ActionDouble
		}
		return engine.ActionHit
	case total == 11:
		if dealer >= 2 && dealer <= 10 {
			return engine.ActionDouble
		}
		return engine.ActionHit
	case total == 12:
		if dealer >= 4 && dealer <= 6 {
			return engine.ActionStand
		}
		return engine.ActionHit
	case total >= 13 && total <= 16:
		if dealer >= 2 && dealer <= 6 {
			return engine.ActionStand
		}
		return engine.ActionHit
	default:
		return engine.ActionStand
	}
}

func chooseAction(desired engine.Action, actions []engine.Action, fallback engine.Action) engine.Action {
	if containsAction(actions, desired) {
		return desired
	}
	if desired == engine.ActionDouble && containsAction(actions, engine.ActionHit) {
		return engine.ActionHit
	}
	if desired == engine.ActionSplit {
		return chooseAction(fallback, actions, engine.ActionStand)
	}
	if containsAction(actions, fallback) {
		return fallback
	}
	return actions[0]
}

func containsAction(actions []engine.Action, target engine.Action) bool {
	for _, a := range actions {
		if strings.EqualFold(string(a), string(target)) {
			return true
		}
	}
	return false
}

func canSplit(actions []engine.Action) bool {
	return containsAction(actions, engine.ActionSplit)
}

func canSurrender(actions []engine.Action) bool {
	return containsAction(actions, engine.ActionSurrender)
}
