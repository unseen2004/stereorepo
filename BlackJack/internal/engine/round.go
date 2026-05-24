package engine

import (
	"errors"

	"BlackJack/internal/rules"
)

var ErrNoActiveHands = errors.New("no active hands")

type Outcome string

const (
	OutcomeWin       Outcome = "win"
	OutcomeLose      Outcome = "lose"
	OutcomePush      Outcome = "push"
	OutcomeBlackjack Outcome = "blackjack"
	OutcomeSurrender Outcome = "surrender"
)

type RoundResult struct {
	PlayerID string
	Hand     int
	Outcome  Outcome
	Payout   int
}

type Round struct {
	Rules         rules.Rules
	Shoe          *Shoe
	Dealer        *Hand
	Players       []*Player
	InsuranceBets map[string]int
	SplitCount    map[string]int
}

func NewRound(r rules.Rules, shoe *Shoe, players []*Player) *Round {
	return &Round{
		Rules:         r,
		Shoe:          shoe,
		Dealer:        NewHand(0),
		Players:       players,
		InsuranceBets: make(map[string]int),
		SplitCount:    make(map[string]int),
	}
}

func (r *Round) Play() ([]RoundResult, error) {
	if r.Shoe.Remaining() < 20 || r.Shoe.ShouldShuffle() {
		r.Shoe.Shuffle()
	}

	for _, p := range r.Players {
		p.Hands = nil
	}
	*r.Dealer = *NewHand(0)

	for _, p := range r.Players {
		p.Provider.Notify(r.buildState(false), Event{Type: EventRoundStart})
	}

	for _, p := range r.Players {
		if p.Bankroll < r.Rules.MinBet {
			continue
		}
		state := r.buildState(false)
		bet, err := p.Provider.BetRequest(state, r.Rules.MinBet, r.Rules.MaxBet)
		if err != nil {
			return nil, err
		}
		if bet < r.Rules.MinBet || bet > r.Rules.MaxBet || bet > p.Bankroll {
			continue
		}
		p.Bankroll -= bet
		p.Hands = append(p.Hands, NewHand(bet))
		p.Provider.Notify(r.buildState(false), Event{Type: EventBetPlaced, Player: p.ID, Hand: 0})
	}

	if !r.hasActiveHands() {
		return nil, ErrNoActiveHands
	}

	if err := r.initialDeal(); err != nil {
		return nil, err
	}
	r.notifyInitialDeal()

	dealerHasBJ := r.dealerHasBlackjack()
	if r.dealerUpcardIsAce() && r.Rules.AllowInsurance {
		r.collectInsurance()
	}
	if r.Rules.Variant == rules.VariantAmerican {
		if dealerHasBJ {
			results := r.resolveImmediateDealerBlackjack()
			return results, nil
		}
	}

	for _, p := range r.Players {
		if len(p.Hands) == 0 {
			continue
		}
		for handIndex := 0; handIndex < len(p.Hands); handIndex++ {
			if err := r.playHand(p, handIndex); err != nil {
				return nil, err
			}
		}
	}

	if err := r.dealerPlay(); err != nil {
		return nil, err
	}
	r.notifyDealerPlay()
	dealerHasBJ = r.Dealer.IsBlackjack()

	results := r.resolveHands(dealerHasBJ)
	for _, p := range r.Players {
		p.Provider.Notify(r.buildState(true), Event{Type: EventRoundEnd, Results: results})
	}
	return results, nil
}

func (r *Round) hasActiveHands() bool {
	for _, p := range r.Players {
		if len(p.Hands) > 0 {
			return true
		}
	}
	return false
}

func (r *Round) initialDeal() error {
	for i := 0; i < 2; i++ {
		for _, p := range r.Players {
			if len(p.Hands) == 0 {
				continue
			}
			p.Hands[0].AddCard(r.Shoe.Draw())
		}
		if r.Rules.Variant == rules.VariantAmerican || i == 0 {
			r.Dealer.AddCard(r.Shoe.Draw())
		}
	}
	return nil
}

func (r *Round) notifyInitialDeal() {
	for _, p := range r.Players {
		cards := make([]Card, 0)
		for _, h := range p.Hands {
			cards = append(cards, h.Cards...)
		}
		cards = append(cards, r.visibleDealerCards()...)
		p.Provider.Notify(r.buildState(false), Event{Type: EventInitialDeal, Cards: cards})
	}
}

func (r *Round) dealerUpcardIsAce() bool {
	if len(r.Dealer.Cards) == 0 {
		return false
	}
	return r.Dealer.Cards[0].Rank == RankAce
}

func (r *Round) collectInsurance() {
	for _, p := range r.Players {
		if len(p.Hands) == 0 {
			continue
		}
		if p.Bankroll < p.Hands[0].Bet/2 {
			continue
		}
		state := r.buildState(false)
		buy, err := p.Provider.InsuranceRequest(state)
		if err != nil || !buy {
			continue
		}
		insurance := p.Hands[0].Bet / 2
		p.Bankroll -= insurance
		r.InsuranceBets[p.ID] = insurance
	}
}

func (r *Round) dealerHasBlackjack() bool {
	return r.Dealer.IsBlackjack()
}

func (r *Round) resolveImmediateDealerBlackjack() []RoundResult {
	results := make([]RoundResult, 0)
	for _, p := range r.Players {
		for i, h := range p.Hands {
			outcome := OutcomeLose
			payout := 0
			if h.IsBlackjack() {
				outcome = OutcomePush
				payout = h.Bet
			}
			if insurance, ok := r.InsuranceBets[p.ID]; ok {
				payout += insurance * 3
			}
			p.Bankroll += payout
			results = append(results, RoundResult{PlayerID: p.ID, Hand: i, Outcome: outcome, Payout: payout})
		}
	}
	return results
}

func (r *Round) playHand(p *Player, handIndex int) error {
	h := p.Hands[handIndex]
	if h.IsBlackjack() {
		h.IsFinished = true
		return nil
	}

	if h.IsSplitHand && len(h.Cards) == 1 {
		h.AddCard(r.Shoe.Draw())
		if h.Cards[0].Rank == RankAce && r.Rules.OneCardOnSplitAces {
			h.IsFinished = true
			return nil
		}
	}

	for !h.IsFinished {
		if h.IsBust() {
			h.IsFinished = true
			break
		}
		actions := r.legalActions(p, handIndex)
		state := r.buildState(false)
		act, err := p.Provider.ActionRequest(state, handIndex, actions)
		if err != nil {
			return err
		}
		if !isActionAllowed(act, actions) {
			p.Provider.Notify(r.buildState(false), Event{Type: EventAction, Player: p.ID, Hand: handIndex, Action: act, Message: "invalid action"})
			continue
		}
		switch act {
		case ActionHit:
			h.AddCard(r.Shoe.Draw())
			p.Provider.Notify(r.buildState(false), Event{Type: EventAction, Player: p.ID, Hand: handIndex, Action: act, Cards: h.Cards[len(h.Cards)-1:]})
		case ActionStand:
			h.IsFinished = true
		case ActionDouble:
			if !containsAction(actions, ActionDouble) {
				continue
			}
			if p.Bankroll < h.Bet {
				continue
			}
			p.Bankroll -= h.Bet
			h.Bet *= 2
			h.IsDoubled = true
			h.AddCard(r.Shoe.Draw())
			p.Provider.Notify(r.buildState(false), Event{Type: EventAction, Player: p.ID, Hand: handIndex, Action: act, Cards: h.Cards[len(h.Cards)-1:]})
			h.IsFinished = true
		case ActionSplit:
			if !containsAction(actions, ActionSplit) {
				continue
			}
			if err := r.splitHand(p, handIndex); err != nil {
				return err
			}
		case ActionSurrender:
			if !containsAction(actions, ActionSurrender) {
				continue
			}
			h.IsSurrender = true
			h.IsFinished = true
		default:
			p.Provider.Notify(r.buildState(false), Event{Type: EventAction, Player: p.ID, Hand: handIndex, Action: act, Message: "unknown action"})
			continue
		}
		if act == ActionStand || act == ActionSplit || act == ActionSurrender {
			p.Provider.Notify(r.buildState(false), Event{Type: EventAction, Player: p.ID, Hand: handIndex, Action: act})
		}
	}

	return nil
}

func isActionAllowed(act Action, actions []Action) bool {
	for _, a := range actions {
		if a == act {
			return true
		}
	}
	return false
}

func (r *Round) splitHand(p *Player, handIndex int) error {
	h := p.Hands[handIndex]
	if len(h.Cards) != 2 {
		return errors.New("split requires two cards")
	}
	if p.Bankroll < h.Bet {
		return errors.New("insufficient bankroll to split")
	}
	p.Bankroll -= h.Bet
	newHand := NewHand(h.Bet)
	newHand.IsSplitHand = true
	h.IsSplitHand = true

	newHand.AddCard(h.Cards[1])
	h.Cards = h.Cards[:1]

	newHand.AddCard(r.Shoe.Draw())
	h.AddCard(r.Shoe.Draw())

	newCards := append([]Card(nil), h.Cards[len(h.Cards)-1])
	newCards = append(newCards, newHand.Cards[len(newHand.Cards)-1])
	p.Provider.Notify(r.buildState(false), Event{Type: EventAction, Player: p.ID, Hand: handIndex, Action: ActionSplit, Cards: newCards})

	if h.Cards[0].Rank == RankAce && r.Rules.OneCardOnSplitAces {
		h.IsFinished = true
		newHand.IsFinished = true
	}

	p.Hands = append(p.Hands, newHand)
	r.SplitCount[p.ID]++
	return nil
}

func (r *Round) legalActions(p *Player, handIndex int) []Action {
	h := p.Hands[handIndex]
	actions := []Action{ActionHit, ActionStand}

	if r.Rules.AllowSurrender && len(h.Cards) == 2 && !h.IsSplitHand {
		actions = append(actions, ActionSurrender)
	}

	if len(h.Cards) == 2 && (!h.IsSplitHand || r.Rules.AllowDoubleAfterSplit) {
		actions = append(actions, ActionDouble)
	}

	if r.canSplit(p, h) {
		actions = append(actions, ActionSplit)
	}

	return actions
}

func (r *Round) canSplit(p *Player, h *Hand) bool {
	if len(h.Cards) != 2 {
		return false
	}
	if !r.Rules.AllowSplitAces && (h.Cards[0].Rank == RankAce || h.Cards[1].Rank == RankAce) {
		return false
	}
	if h.Cards[0].Rank != h.Cards[1].Rank {
		return false
	}
	if !r.Rules.AllowResplit && h.IsSplitHand {
		return false
	}
	if r.SplitCount[p.ID] >= r.Rules.MaxSplits {
		return false
	}
	if p.Bankroll < h.Bet {
		return false
	}
	return true
}

func (r *Round) dealerPlay() error {
	if r.Rules.Variant == rules.VariantEuropean && len(r.Dealer.Cards) == 1 {
		r.Dealer.AddCard(r.Shoe.Draw())
	}

	for {
		total, soft := r.Dealer.SoftTotal()
		if total > 21 {
			break
		}
		if total < 17 {
			r.Dealer.AddCard(r.Shoe.Draw())
			continue
		}
		if total == 17 && soft && r.Rules.DealerHitsSoft17 {
			r.Dealer.AddCard(r.Shoe.Draw())
			continue
		}
		break
	}
	return nil
}

func (r *Round) notifyDealerPlay() {
	for _, p := range r.Players {
		p.Provider.Notify(r.buildState(true), Event{Type: EventDealerPlay, Cards: r.Dealer.Cards})
	}
}

func (r *Round) visibleDealerCards() []Card {
	if r.Rules.Variant == rules.VariantAmerican && len(r.Dealer.Cards) == 2 {
		return r.Dealer.Cards[:1]
	}
	return append([]Card(nil), r.Dealer.Cards...)
}

func (r *Round) resolveHands(dealerBlackjack bool) []RoundResult {
	results := make([]RoundResult, 0)
	dealerTotal := r.Dealer.BestTotal()
	dealerBust := dealerTotal > 21

	for _, p := range r.Players {
		for i, h := range p.Hands {
			outcome := OutcomeLose
			payout := 0

			if h.IsSurrender {
				outcome = OutcomeSurrender
				payout = h.Bet / 2
			} else if h.IsBust() {
				outcome = OutcomeLose
				payout = 0
			} else if dealerBlackjack && r.Rules.Variant == rules.VariantEuropean {
				if h.IsBlackjack() && !h.IsSplitHand {
					outcome = OutcomePush
					payout = h.Bet
				} else if r.Rules.EuropeanOBO {
					outcome = OutcomeLose
					payout = h.Bet - h.OriginalBet
				} else {
					outcome = OutcomeLose
					payout = 0
				}
			} else {
				playerTotal := h.BestTotal()
				if h.IsBlackjack() && !h.IsSplitHand {
					if dealerTotal == 21 && len(r.Dealer.Cards) == 2 {
						outcome = OutcomePush
						payout = h.Bet
					} else {
						outcome = OutcomeBlackjack
						payout = h.Bet + int(float64(h.Bet)*r.Rules.BlackjackPayout)
					}
				} else if dealerBust {
					outcome = OutcomeWin
					payout = h.Bet * 2
				} else if playerTotal > dealerTotal {
					outcome = OutcomeWin
					payout = h.Bet * 2
				} else if playerTotal == dealerTotal {
					outcome = OutcomePush
					payout = h.Bet
				} else {
					outcome = OutcomeLose
					payout = 0
				}
			}

			p.Bankroll += payout
			results = append(results, RoundResult{PlayerID: p.ID, Hand: i, Outcome: outcome, Payout: payout})
		}
	}
	return results
}

func (r *Round) buildState(revealDealer bool) GameState {
	players := make([]PlayerState, 0, len(r.Players))
	for _, p := range r.Players {
		hands := make([]HandState, 0, len(p.Hands))
		for _, h := range p.Hands {
			hands = append(hands, HandState{
				Cards:       append([]Card(nil), h.Cards...),
				Bet:         h.Bet,
				IsSplitHand: h.IsSplitHand,
				IsDoubled:   h.IsDoubled,
				IsFinished:  h.IsFinished,
				IsSurrender: h.IsSurrender,
			})
		}
		players = append(players, PlayerState{
			ID:       p.ID,
			Name:     p.Name,
			Bankroll: p.Bankroll,
			Hands:    hands,
		})
	}

	dealerCards := append([]Card(nil), r.Dealer.Cards...)
	if !revealDealer && r.Rules.Variant == rules.VariantAmerican && len(dealerCards) == 2 {
		dealerCards = dealerCards[:1]
	}

	return GameState{
		Rules:         r.Rules,
		DealerCards:   dealerCards,
		Players:       players,
		ShoeRemaining: r.Shoe.Remaining(),
		RoundOver:     false,
	}
}

func containsAction(actions []Action, target Action) bool {
	for _, a := range actions {
		if a == target {
			return true
		}
	}
	return false
}
