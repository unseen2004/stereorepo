package netplay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"BlackJack/internal/engine"
)

type Client struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Conn:   conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}

func (c *Client) Run() error {
	dec := json.NewDecoder(c.Reader)
	enc := json.NewEncoder(c.Writer)
	stdin := bufio.NewReader(os.Stdin)

	for {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			return err
		}

		switch msg.Type {
		case "bet":
			var payload struct {
				Min int `json:"min"`
				Max int `json:"max"`
			}
			_ = json.Unmarshal(msg.Payload, &payload)
			fmt.Printf("Bet (%d-%d): ", payload.Min, payload.Max)
			var bet int
			fmt.Fscanln(stdin, &bet)
			resp, _ := NewMessage("bet", map[string]int{"bet": bet})
			if err := enc.Encode(resp); err != nil {
				return err
			}
			c.Writer.Flush()
		case "insurance":
			fmt.Print("Insurance? (y/n): ")
			var v string
			fmt.Fscanln(stdin, &v)
			buy := v == "y" || v == "yes"
			resp, _ := NewMessage("insurance", map[string]bool{"buy": buy})
			if err := enc.Encode(resp); err != nil {
				return err
			}
			c.Writer.Flush()
		case "action":
			var payload struct {
				Hand    int      `json:"hand"`
				Actions []string `json:"actions"`
			}
			_ = json.Unmarshal(msg.Payload, &payload)
			fmt.Printf("Action for hand %d (%v): ", payload.Hand+1, payload.Actions)
			var act string
			fmt.Fscanln(stdin, &act)
			resp, _ := NewMessage("action", map[string]string{"action": act})
			if err := enc.Encode(resp); err != nil {
				return err
			}
			c.Writer.Flush()
		case "event":
			fmt.Println("Event:", string(msg.Payload))
		default:
			fmt.Println("Message:", string(msg.Payload))
		}
	}
}

func cardsToString(cards []engine.Card) string {
	parts := make([]string, 0, len(cards))
	for _, c := range cards {
		parts = append(parts, c.String())
	}
	return strings.Join(parts, ",")
}
