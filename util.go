package santase

func getHiddenCards(hand Hand, trumpCard Card) Pile {
	remaining := NewPile()
	for _, card := range AllCards {
		isInHand := hand.HasCard(card)
		isTrumpCard := card == trumpCard
		if !isInHand && !isTrumpCard {
			remaining.AddCard(card)
		}
	}
	return remaining
}

// StrongerCard returns the stronger of two cards played in a game with a particular trump.
func StrongerCard(first *Card, second *Card, trump Suit) *Card {
	if first.Suit == second.Suit {
		if first.Rank > second.Rank {
			return first
		}
		return second
	}

	if first.Suit == trump {
		return first
	}

	if second.Suit == trump {
		return second
	}

	return first
}

var pointsMap = map[Rank]int{
	Nine:  0,
	Jack:  2,
	Queen: 3,
	King:  4,
	Ten:   10,
	Ace:   11,
}

// Points returns the value of a card depending on its rank.
// The ranks are valued as follows:
//
//	Nine  = 0 points
// 	Jack  = 2 points
// 	Queen = 3 points
// 	King  = 4 points
// 	Ten   = 10 points
// 	Ace   = 11 points
func Points(c *Card) int {
	if pts, ok := pointsMap[c.Rank]; ok {
		return pts
	}

	panic("invalid card")
}
