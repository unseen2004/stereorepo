package netplay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"BlackJack/internal/engine"
)

type Server struct {
	Listener net.Listener
}

func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{Listener: ln}, nil
}

func (s *Server) Close() error {
	return s.Listener.Close()
}

func (s *Server) WaitForPlayers(n int, timeout time.Duration) ([]*RemoteProvider, error) {
	providers := make([]*RemoteProvider, 0, n)
	for len(providers) < n {
		if timeout > 0 {
			if tcp, ok := s.Listener.(*net.TCPListener); ok {
				_ = tcp.SetDeadline(time.Now().Add(timeout))
			}
		}
		conn, err := s.Listener.Accept()
		if err != nil {
			return nil, err
		}
		providers = append(providers, NewRemoteProvider(conn))
	}
	return providers, nil
}

type RemoteProvider struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
	Name   string
}

func NewRemoteProvider(conn net.Conn) *RemoteProvider {
	return &RemoteProvider{
		Conn:   conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}
}

func (rp *RemoteProvider) Close() error {
	return rp.Conn.Close()
}

func (rp *RemoteProvider) BetRequest(state engine.GameState, minBet, maxBet int) (int, error) {
	player := findPlayer(state, rp.Name)
	bankroll := 0
	if player != nil {
		bankroll = player.Bankroll
	}
	msg, err := NewMessage("bet", map[string]int{"min": minBet, "max": maxBet, "bankroll": bankroll})
	if err != nil {
		return 0, err
	}
	if err := rp.send(msg); err != nil {
		return 0, err
	}
	var resp struct {
		Bet int `json:"bet"`
	}
	if err := rp.recv("bet", &resp); err != nil {
		return 0, err
	}
	return resp.Bet, nil
}

func (rp *RemoteProvider) InsuranceRequest(state engine.GameState) (bool, error) {
	msg, err := NewMessage("insurance", map[string]any{"offer": true, "dealer": state.DealerCards})
	if err != nil {
		return false, err
	}
	if err := rp.send(msg); err != nil {
		return false, err
	}
	var resp struct {
		Buy bool `json:"buy"`
	}
	if err := rp.recv("insurance", &resp); err != nil {
		return false, err
	}
	return resp.Buy, nil
}

func (rp *RemoteProvider) ActionRequest(state engine.GameState, handIndex int, actions []engine.Action) (engine.Action, error) {
	player := findPlayer(state, rp.Name)
	var hand engine.HandState
	if player != nil && handIndex >= 0 && handIndex < len(player.Hands) {
		hand = player.Hands[handIndex]
	}
	msg, err := NewMessage("action", map[string]any{
		"hand":         handIndex,
		"actions":      actions,
		"dealer_cards": state.DealerCards,
		"player_cards": hand.Cards,
		"bankroll":     player.Bankroll,
	})
	if err != nil {
		return "", err
	}
	if err := rp.send(msg); err != nil {
		return "", err
	}
	var resp struct {
		Action engine.Action `json:"action"`
	}
	if err := rp.recv("action", &resp); err != nil {
		return "", err
	}
	return resp.Action, nil
}

func (rp *RemoteProvider) Notify(state engine.GameState, event engine.Event) {
	msg, err := NewMessage("event", event)
	if err != nil {
		return
	}
	_ = rp.send(msg)
}

func (rp *RemoteProvider) send(msg Message) error {
	enc := json.NewEncoder(rp.Writer)
	if err := enc.Encode(msg); err != nil {
		return err
	}
	return rp.Writer.Flush()
}

func findPlayer(state engine.GameState, name string) *engine.PlayerState {
	for i := range state.Players {
		if state.Players[i].Name == name {
			return &state.Players[i]
		}
	}
	return nil
}

func (rp *RemoteProvider) recv(expected string, out any) error {
	dec := json.NewDecoder(rp.Reader)
	var msg Message
	if err := dec.Decode(&msg); err != nil {
		return err
	}
	if msg.Type != expected {
		return fmt.Errorf("unexpected message: %s", msg.Type)
	}
	return json.Unmarshal(msg.Payload, out)
}
