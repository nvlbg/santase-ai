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
func StrongerCard(a *Card, b *Card, trump Suit) *Card {
	if a.Suit == b.Suit {
		if a.Rank > b.Rank {
			return a
		}
		return b
	}

	if a.Suit == trump {
		return a
	}

	if b.Suit == trump {
		return b
	}

	return a
}

var pointsMap = map[Rank]int{
	Nine:  0,
	Jack:  2,
	Queen: 3,
	King:  4,
	Ten:   10,
	Ace:   11,
}

func Points(c *Card) int {
	if pts, ok := pointsMap[c.Rank]; ok {
		return pts
	}

	panic("invalid card")
}
