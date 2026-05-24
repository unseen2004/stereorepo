package engine

import (
	"testing"

	"BlackJack/internal/rules"
)

type fixedProvider struct{}

func (f fixedProvider) BetRequest(state GameState, minBet, maxBet int) (int, error) {
	return minBet, nil
}

func (f fixedProvider) InsuranceRequest(state GameState) (bool, error) {
	return false, nil
}

func (f fixedProvider) ActionRequest(state GameState, handIndex int, actions []Action) (Action, error) {
	for _, a := range actions {
		if a == ActionStand {
			return ActionStand, nil
		}
	}
	return ActionStand, nil
}

func (f fixedProvider) Notify(state GameState, event Event) {}

func TestDealerHitsSoft17(t *testing.T) {
	r := rules.DefaultRules(rules.VariantAmerican)
	r.DealerHitsSoft17 = true
	shoe := &Shoe{Decks: 1}
	shoe.Cards = []Card{
		{Rank: RankFive},
		{Rank: RankTwo},
		{Rank: RankSeven},
		{Rank: RankAce},
		{Rank: RankSix},
	}
	player := &Player{ID: "P1", Name: "P1", Bankroll: 100, Provider: fixedProvider{}}
	round := NewRound(r, shoe, []*Player{player})
	round.Dealer = NewHand(0)
	round.Dealer.AddCard(Card{Rank: RankAce})
	round.Dealer.AddCard(Card{Rank: RankSix})
	if err := round.dealerPlay(); err != nil {
		t.Fatalf("dealerPlay error: %v", err)
	}
	if len(round.Dealer.Cards) < 3 {
		t.Fatalf("expected dealer to hit on soft 17")
	}
}

func TestSplitHand(t *testing.T) {
	r := rules.DefaultRules(rules.VariantEuropean)
	shoe := &Shoe{Decks: 1}
	shoe.Cards = []Card{
		{Rank: RankFive},
		{Rank: RankFour},
	}
	player := &Player{ID: "P1", Name: "P1", Bankroll: 100, Provider: fixedProvider{}}
	hand := NewHand(10)
	hand.Cards = []Card{{Rank: RankEight}, {Rank: RankEight}}
	player.Hands = []*Hand{hand}

	round := NewRound(r, shoe, []*Player{player})
	if err := round.splitHand(player, 0); err != nil {
		t.Fatalf("splitHand error: %v", err)
	}
	if len(player.Hands) != 2 {
		t.Fatalf("expected 2 hands after split")
	}
	if player.Bankroll != 90 {
		t.Fatalf("expected bankroll 90, got %d", player.Bankroll)
	}
}
