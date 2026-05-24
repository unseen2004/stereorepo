package bot

import (
	"math"

	"BlackJack/internal/engine"
)

type HiLoBot struct {
	Name         string
	RunningCount int
}

func (b *HiLoBot) BetRequest(state engine.GameState, minBet, maxBet int) (int, error) {
	decks := math.Max(1, float64(state.ShoeRemaining)/52.0)
	trueCount := int(math.Floor(float64(b.RunningCount) / decks))
	betUnits := 1
	if trueCount > 1 {
		betUnits = trueCount
	}
	bet := minBet * betUnits
	if bet > maxBet {
		bet = maxBet
	}
	return bet, nil
}

func (b *HiLoBot) InsuranceRequest(state engine.GameState) (bool, error) {
	return b.RunningCount >= 3, nil
}

func (b *HiLoBot) ActionRequest(state engine.GameState, handIndex int, actions []engine.Action) (engine.Action, error) {
	player := findPlayerState(state, b.Name)
	if player == nil || handIndex >= len(player.Hands) {
		return engine.ActionStand, nil
	}
	dealerUp := dealerUpcard(state)
	hand := player.Hands[handIndex]
	rules := strategyRules{DealerHitsSoft17: state.Rules.DealerHitsSoft17, AllowDoubleAfterSplit: state.Rules.AllowDoubleAfterSplit}
	return decideBasicAction(hand, dealerUp, rules, actions), nil
}

func (b *HiLoBot) Notify(state engine.GameState, event engine.Event) {
	for _, c := range event.Cards {
		b.RunningCount += hiLoValue(c)
	}
}

func hiLoValue(card engine.Card) int {
	switch card.Rank {
	case engine.RankTwo, engine.RankThree, engine.RankFour, engine.RankFive, engine.RankSix:
		return 1
	case engine.RankSeven, engine.RankEight, engine.RankNine:
		return 0
	default:
		return -1
	}
}
