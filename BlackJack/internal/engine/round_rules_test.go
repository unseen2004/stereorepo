package engine

import (
	"testing"

	"BlackJack/internal/rules"
)

type scriptedProvider struct {
	Bet       int
	Insurance bool
}

func (s scriptedProvider) BetRequest(state GameState, minBet, maxBet int) (int, error) {
	return s.Bet, nil
}

func (s scriptedProvider) InsuranceRequest(state GameState) (bool, error) {
	return s.Insurance, nil
}

func (s scriptedProvider) ActionRequest(state GameState, handIndex int, actions []Action) (Action, error) {
	for _, a := range actions {
		if a == ActionStand {
			return ActionStand, nil
		}
	}
	return ActionStand, nil
}

func (s scriptedProvider) Notify(state GameState, event Event) {}

func TestEuropeanDealerBlackjackPushesPlayerBlackjack(t *testing.T) {
	r := rules.DefaultRules(rules.VariantEuropean)
	shoe := &Shoe{Decks: 1, Penetration: 1}
	shoe.Cards = makeShoe([]Card{
		{Rank: RankKing},
		{Rank: RankKing},
		{Rank: RankAce},
		{Rank: RankAce},
	})

	player := &Player{ID: "P1", Name: "P1", Bankroll: 100, Provider: scriptedProvider{Bet: 10}}
	round := NewRound(r, shoe, []*Player{player})
	results, err := round.Play()
	if err != nil {
		t.Fatalf("Play error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Outcome != OutcomePush {
		t.Fatalf("expected push, got %s", results[0].Outcome)
	}
	if results[0].Payout != 10 {
		t.Fatalf("expected payout 10, got %d", results[0].Payout)
	}
}

func TestEuropeanOBOReturnsDoubleExtra(t *testing.T) {
	r := rules.DefaultRules(rules.VariantEuropean)
	r.EuropeanOBO = true
	shoe := &Shoe{Decks: 1, Penetration: 0}

	provider := scriptedProvider{Bet: 10}
	player := &Player{ID: "P1", Name: "P1", Bankroll: 100, Provider: provider}
	round := NewRound(r, shoe, []*Player{player})
	round.Dealer = NewHand(0)
	round.Dealer.AddCard(Card{Rank: RankAce})
	round.Dealer.AddCard(Card{Rank: RankKing})

	hand := NewHand(10)
	hand.AddCard(Card{Rank: RankTen})
	hand.AddCard(Card{Rank: RankTen})
	hand.Bet = 20
	hand.IsDoubled = true
	player.Hands = []*Hand{hand}
	player.Bankroll = 80

	results := round.resolveHands(true)
	if results[0].Payout != 10 {
		t.Fatalf("expected payout 10 for OBO, got %d", results[0].Payout)
	}
	if player.Bankroll != 90 {
		t.Fatalf("expected bankroll 90, got %d", player.Bankroll)
	}
}

func TestAmericanDealerBlackjackImmediateLoss(t *testing.T) {
	r := rules.DefaultRules(rules.VariantAmerican)
	shoe := &Shoe{Decks: 1, Penetration: 1}
	shoe.Cards = makeShoe([]Card{
		{Rank: RankKing},
		{Rank: RankTen},
		{Rank: RankAce},
		{Rank: RankNine},
	})

	player := &Player{ID: "P1", Name: "P1", Bankroll: 100, Provider: scriptedProvider{Bet: 10}}
	round := NewRound(r, shoe, []*Player{player})
	results, err := round.Play()
	if err != nil {
		t.Fatalf("Play error: %v", err)
	}
	if results[0].Outcome != OutcomeLose {
		t.Fatalf("expected loss, got %s", results[0].Outcome)
	}
	if results[0].Payout != 0 {
		t.Fatalf("expected payout 0, got %d", results[0].Payout)
	}
}

func TestInsurancePayoutOnDealerBlackjack(t *testing.T) {
	r := rules.DefaultRules(rules.VariantAmerican)
	shoe := &Shoe{Decks: 1, Penetration: 1}
	shoe.Cards = makeShoe([]Card{
		{Rank: RankTen},
		{Rank: RankNine},
		{Rank: RankAce},
		{Rank: RankTen},
	})

	player := &Player{ID: "P1", Name: "P1", Bankroll: 100, Provider: scriptedProvider{Bet: 10, Insurance: true}}
	round := NewRound(r, shoe, []*Player{player})
	results, err := round.Play()
	if err != nil {
		t.Fatalf("Play error: %v", err)
	}
	if results[0].Payout != 15 {
		t.Fatalf("expected insurance payout 15, got %d", results[0].Payout)
	}
	if player.Bankroll != 100 {
		t.Fatalf("expected bankroll 100, got %d", player.Bankroll)
	}
}

func makeShoe(drawOrder []Card) []Card {
	cards := make([]Card, 21)
	for i := range cards {
		cards[i] = Card{Rank: RankTwo}
	}
	return append(cards, drawOrder...)
}
