package engine

import (
	"math/rand"
	"time"
)

type Shoe struct {
	Decks       int
	Cards       []Card
	Penetration float64
	Rand        *rand.Rand
}

func NewShoe(decks int, penetration float64, seed int64) *Shoe {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))
	s := &Shoe{Decks: decks, Penetration: penetration, Rand: rng}
	s.Shuffle()
	return s
}

func (s *Shoe) Shuffle() {
	if s.Rand == nil {
		s.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	s.Cards = s.Cards[:0]
	for d := 0; d < s.Decks; d++ {
		for r := RankAce; r <= RankKing; r++ {
			for i := 0; i < 4; i++ {
				s.Cards = append(s.Cards, Card{Rank: r})
			}
		}
	}
	s.Rand.Shuffle(len(s.Cards), func(i, j int) {
		s.Cards[i], s.Cards[j] = s.Cards[j], s.Cards[i]
	})
}

func (s *Shoe) Remaining() int {
	return len(s.Cards)
}

func (s *Shoe) ShouldShuffle() bool {
	total := 52 * s.Decks
	dealt := total - len(s.Cards)
	threshold := int(float64(total) * s.Penetration)
	return dealt >= threshold
}

func (s *Shoe) Draw() Card {
	if len(s.Cards) == 0 {
		s.Shuffle()
	}
	c := s.Cards[len(s.Cards)-1]
	s.Cards = s.Cards[:len(s.Cards)-1]
	return c
}
