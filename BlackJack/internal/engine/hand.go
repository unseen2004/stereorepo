package engine

type Hand struct {
	Cards       []Card
	Bet         int
	OriginalBet int
	IsSplitHand bool
	IsDoubled   bool
	IsFinished  bool
	IsSurrender bool
}

func NewHand(bet int) *Hand {
	return &Hand{Bet: bet, OriginalBet: bet}
}

func (h *Hand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
}

func (h *Hand) SoftTotal() (total int, soft bool) {
	aces := 0
	for _, c := range h.Cards {
		if c.Rank == RankAce {
			aces++
		}
		total += c.Value()
	}

	if aces > 0 && total+10 <= 21 {
		return total + 10, true
	}

	return total, false
}

func (h *Hand) BestTotal() int {
	t, _ := h.SoftTotal()
	return t
}

func (h *Hand) IsBlackjack() bool {
	if len(h.Cards) != 2 {
		return false
	}
	return h.BestTotal() == 21
}

func (h *Hand) IsBust() bool {
	return h.BestTotal() > 21
}
