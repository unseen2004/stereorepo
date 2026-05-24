package engine

import "testing"

func TestHandTotals(t *testing.T) {
	h := NewHand(10)
	h.AddCard(Card{Rank: RankAce})
	h.AddCard(Card{Rank: RankSix})
	if total, soft := h.SoftTotal(); total != 17 || !soft {
		t.Fatalf("expected soft 17, got %d soft=%v", total, soft)
	}
	h.AddCard(Card{Rank: RankNine})
	if total, soft := h.SoftTotal(); total != 16 || soft {
		t.Fatalf("expected hard 16, got %d soft=%v", total, soft)
	}
}

func TestBlackjack(t *testing.T) {
	h := NewHand(10)
	h.AddCard(Card{Rank: RankAce})
	h.AddCard(Card{Rank: RankKing})
	if !h.IsBlackjack() {
		t.Fatalf("expected blackjack")
	}
}
