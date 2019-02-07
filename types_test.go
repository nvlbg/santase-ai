package santase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHand(t *testing.T) {
	NewHand(
		NewCard(Nine, Hearts),
		NewCard(Jack, Hearts),
		NewCard(Queen, Hearts),
		NewCard(King, Hearts),
		NewCard(Ten, Hearts),
		NewCard(Ace, Hearts),
	)
}

func TestNewHandOverfilling(t *testing.T) {
	assert.PanicsWithValue(
		t, "too many cards given",
		func() {
			NewHand(
				NewCard(Nine, Hearts),
				NewCard(Jack, Hearts),
				NewCard(Queen, Hearts),
				NewCard(King, Hearts),
				NewCard(Ten, Hearts),
				NewCard(Ace, Hearts),
				NewCard(Ace, Spades),
			)
		},
	)
}

func TestAddCard(t *testing.T) {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))
	hand.AddCard(NewCard(Jack, Hearts))
	hand.AddCard(NewCard(Queen, Hearts))
	hand.AddCard(NewCard(King, Hearts))
	hand.AddCard(NewCard(Ten, Hearts))
	hand.AddCard(NewCard(Ace, Hearts))
}

func TestAddCardOverfilling(t *testing.T) {
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

func TestToSlice(t *testing.T) {
	hand := NewHand(
		NewCard(Nine, Hearts),
		NewCard(Jack, Hearts),
		NewCard(Queen, Hearts),
		NewCard(King, Hearts),
		NewCard(Ten, Hearts),
		NewCard(Ace, Hearts),
	)

	slice := hand.ToSlice()
	assert.Equal(t, 6, len(slice))
	assert.Contains(t, slice, NewCard(Nine, Hearts))
	assert.Contains(t, slice, NewCard(Jack, Hearts))
	assert.Contains(t, slice, NewCard(Queen, Hearts))
	assert.Contains(t, slice, NewCard(King, Hearts))
	assert.Contains(t, slice, NewCard(Ten, Hearts))
	assert.Contains(t, slice, NewCard(Ace, Hearts))
}

func TestString(t *testing.T) {
	hand := NewHand(
		NewCard(Nine, Hearts),
		NewCard(Jack, Hearts),
		NewCard(Queen, Hearts),
		NewCard(King, Hearts),
		NewCard(Ten, Hearts),
		NewCard(Ace, Hearts),
	)

	assert.Equal(t, "{ 9♥ J♥ Q♥ K♥ 10♥ A♥ }", hand.String())
}
