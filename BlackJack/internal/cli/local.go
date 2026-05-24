package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"BlackJack/internal/engine"
)

type LocalPlayer struct {
	Name   string
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func NewLocalPlayer(name string) *LocalPlayer {
	return &LocalPlayer{
		Name:   name,
		Reader: bufio.NewReader(os.Stdin),
		Writer: bufio.NewWriter(os.Stdout),
	}
}

func (lp *LocalPlayer) BetRequest(state engine.GameState, minBet, maxBet int) (int, error) {
	bankroll := playerBankroll(state, lp.Name)
	for {
		fmt.Fprintf(lp.Writer, "%s bankroll=%d bet (%d-%d): ", lp.Name, bankroll, minBet, maxBet)
		lp.Writer.Flush()
		var bet int
		if _, err := fmt.Fscanln(lp.Reader, &bet); err != nil {
			return 0, err
		}
		if bet >= minBet && bet <= maxBet && bet <= bankroll {
			return bet, nil
		}
		fmt.Fprintln(lp.Writer, "invalid bet")
		lp.Writer.Flush()
	}
}

func (lp *LocalPlayer) InsuranceRequest(state engine.GameState) (bool, error) {
	for {
		fmt.Fprintf(lp.Writer, "%s insurance? (y/n): ", lp.Name)
		lp.Writer.Flush()
		line, err := lp.Reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "y" || line == "yes" {
			return true, nil
		}
		if line == "n" || line == "no" {
			return false, nil
		}
		fmt.Fprintln(lp.Writer, "invalid input")
		lp.Writer.Flush()
	}
}

func (lp *LocalPlayer) ActionRequest(state engine.GameState, handIndex int, actions []engine.Action) (engine.Action, error) {
	player := findPlayer(state, lp.Name)
	dealer := cardsToString(state.DealerCards)
	if player != nil && handIndex < len(player.Hands) {
		hand := player.Hands[handIndex]
		fmt.Fprintf(lp.Writer, "Dealer: %s\n", dealer)
		fmt.Fprintf(lp.Writer, "%s hand %d: %s\n", lp.Name, handIndex+1, cardsToString(hand.Cards))
	}
	for {
		fmt.Fprintf(lp.Writer, "%s hand %d action %v: ", lp.Name, handIndex+1, actions)
		lp.Writer.Flush()
		line, err := lp.Reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.ToLower(strings.TrimSpace(line))
		act := engine.Action(line)
		if containsAction(actions, act) {
			return act, nil
		}
		fmt.Fprintln(lp.Writer, "invalid action")
		lp.Writer.Flush()
	}
}

func (lp *LocalPlayer) Notify(state engine.GameState, event engine.Event) {
	switch event.Type {
	case engine.EventRoundStart:
		fmt.Fprintf(lp.Writer, "Round start\n")
	case engine.EventInitialDeal:
		fmt.Fprintf(lp.Writer, "Initial deal\n")
	case engine.EventDealerPlay:
		fmt.Fprintf(lp.Writer, "Dealer plays: %s\n", cardsToString(event.Cards))
	case engine.EventRoundEnd:
		fmt.Fprintf(lp.Writer, "Round end\n")
	}
	if event.Message != "" {
		fmt.Fprintf(lp.Writer, "%s\n", event.Message)
	}
	lp.Writer.Flush()
}

func cardsToString(cards []engine.Card) string {
	parts := make([]string, 0, len(cards))
	for _, c := range cards {
		parts = append(parts, c.String())
	}
	return strings.Join(parts, ",")
}

func findPlayer(state engine.GameState, name string) *engine.PlayerState {
	for i := range state.Players {
		if state.Players[i].Name == name {
			return &state.Players[i]
		}
	}
	return nil
}

func playerBankroll(state engine.GameState, name string) int {
	if p := findPlayer(state, name); p != nil {
		return p.Bankroll
	}
	return 0
}

func containsAction(actions []engine.Action, target engine.Action) bool {
	for _, a := range actions {
		if a == target {
			return true
		}
	}
	return false
}
