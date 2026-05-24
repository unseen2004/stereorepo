package engine

type Rank int

const (
	RankAce   Rank = 1
	RankTwo   Rank = 2
	RankThree Rank = 3
	RankFour  Rank = 4
	RankFive  Rank = 5
	RankSix   Rank = 6
	RankSeven Rank = 7
	RankEight Rank = 8
	RankNine  Rank = 9
	RankTen   Rank = 10
	RankJack  Rank = 11
	RankQueen Rank = 12
	RankKing  Rank = 13
)

type Card struct {
	Rank Rank
}

func (c Card) Value() int {
	switch c.Rank {
	case RankJack, RankQueen, RankKing:
		return 10
	case RankAce:
		return 1
	default:
		return int(c.Rank)
	}
}

func (c Card) String() string {
	switch c.Rank {
	case RankAce:
		return "A"
	case RankJack:
		return "J"
	case RankQueen:
		return "Q"
	case RankKing:
		return "K"
	default:
		return intToString(int(c.Rank))
	}
}

func intToString(v int) string {
	if v == 10 {
		return "10"
	}
	return string(rune('0' + v))
}
