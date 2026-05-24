package ml

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"BlackJack/internal/engine"
)

type Agent struct {
	Name   string
	Reader *bufio.Reader
	Writer *bufio.Writer
}

type Observation struct {
	DealerCards []engine.Card      `json:"dealer_cards"`
	Hands       []engine.HandState `json:"hands"`
	Bankroll    int                `json:"bankroll"`
	Actions     []engine.Action    `json:"actions"`
	HandIndex   int                `json:"hand_index"`
}

func NewAgent(name string) *Agent {
	return &Agent{
		Name:   name,
		Reader: bufio.NewReader(os.Stdin),
		Writer: bufio.NewWriter(os.Stdout),
	}
}

func (a *Agent) BetRequest(state engine.GameState, minBet, maxBet int) (int, error) {
	obs := map[string]int{"min_bet": minBet, "max_bet": maxBet}
	if err := writeJSON(a.Writer, obs); err != nil {
		return 0, err
	}
	line, err := readLine(a.Reader)
	if err != nil {
		return 0, err
	}
	bet, err := parseInt(line, "bet")
	if err != nil {
		return 0, err
	}
	return bet, nil
}

func (a *Agent) InsuranceRequest(state engine.GameState) (bool, error) {
	obs := map[string]string{"insurance": "offer"}
	if err := writeJSON(a.Writer, obs); err != nil {
		return false, err
	}
	line, err := readLine(a.Reader)
	if err != nil {
		return false, err
	}
	val, err := parseString(line, "insurance")
	if err != nil {
		return false, err
	}
	val = strings.ToLower(val)
	return val == "y" || val == "yes" || val == "true", nil
}

func (a *Agent) ActionRequest(state engine.GameState, handIndex int, actions []engine.Action) (engine.Action, error) {
	player := findPlayer(state, a.Name)
	obs := Observation{
		DealerCards: state.DealerCards,
		Hands:       player.Hands,
		Bankroll:    player.Bankroll,
		Actions:     actions,
		HandIndex:   handIndex,
	}
	if err := writeJSON(a.Writer, obs); err != nil {
		return "", err
	}
	line, err := readLine(a.Reader)
	if err != nil {
		return "", err
	}
	action, err := parseString(line, "action")
	if err != nil {
		return "", err
	}
	return engine.Action(action), nil
}

func (a *Agent) Notify(state engine.GameState, event engine.Event) {}

func writeJSON(w *bufio.Writer, v any) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		return err
	}
	return w.Flush()
}

func findPlayer(state engine.GameState, name string) engine.PlayerState {
	for _, p := range state.Players {
		if p.Name == name {
			return p
		}
	}
	return engine.PlayerState{}
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func parseString(line, key string) (string, error) {
	if strings.HasPrefix(line, "{") {
		var payload map[string]any
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			return "", err
		}
		if v, ok := payload[key]; ok {
			if s, ok := v.(string); ok {
				return s, nil
			}
		}
		return "", fmt.Errorf("missing %s", key)
	}
	return line, nil
}

func parseInt(line, key string) (int, error) {
	if strings.HasPrefix(line, "{") {
		var payload map[string]any
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			return 0, err
		}
		if v, ok := payload[key]; ok {
			if f, ok := v.(float64); ok {
				return int(f), nil
			}
		}
		return 0, fmt.Errorf("missing %s", key)
	}
	var val int
	_, err := fmt.Sscan(line, &val)
	return val, err
}
