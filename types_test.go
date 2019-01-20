package santase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHand(t *testing.T) {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))
	hand.AddCard(NewCard(Jack, Hearts))
	hand.AddCard(NewCard(Queen, Hearts))
	hand.AddCard(NewCard(King, Hearts))
	hand.AddCard(NewCard(Ten, Hearts))
	hand.AddCard(NewCard(Ace, Hearts))
}

func TestNewHandOverfilling(t *testing.T) {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))
	hand.AddCard(NewCard(Jack, Hearts))
	hand.AddCard(NewCard(Queen, Hearts))
	hand.AddCard(NewCard(King, Hearts))
	hand.AddCard(NewCard(Ten, Hearts))
	hand.AddCard(NewCard(Ace, Hearts))

	assert.PanicsWithValue(
		t, "hand has 6 cards already",
		func() { hand.AddCard(NewCard(Ace, Spades)) },
	)
}

func TestNewHandAddingSameCardTwice(t *testing.T) {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))
	hand.AddCard(NewCard(Nine, Hearts))
}

func TestHasCard(t *testing.T) {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))

	assert.False(t, hand.HasCard(NewCard(Ten, Hearts)))
	assert.True(t, hand.HasCard(NewCard(Nine, Hearts)))
}

func TestRemoveCard(t *testing.T) {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))

	assert.True(t, hand.HasCard(NewCard(Nine, Hearts)))
	hand.RemoveCard(NewCard(Nine, Hearts))
	assert.False(t, hand.HasCard(NewCard(Nine, Hearts)))
}
