package santase

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

	assert.PanicsWithValue(
		t, "9♥ is already in the hand",
		func() { hand.AddCard(NewCard(Nine, Hearts)) },
	)
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

func TestNewMove(t *testing.T) {
	// Check this compiles and runs without panics
	NewMove(NewCard(Nine, Clubs))
}

func TestNewMoveWithAnnouncement(t *testing.T) {
	// Check this compiles and runs without panics
	NewMoveWithAnnouncement(NewCard(Queen, Clubs))
}

func TestNewMoveWithAnnouncementImpossibleAnnouncement(t *testing.T) {
	assert.PanicsWithValue(
		t, "announcement moves are only possible with queens and kings",
		func() { NewMoveWithAnnouncement(NewCard(Nine, Clubs)) },
	)
}

func TestNewMoveWithTrumpCardSwitch(t *testing.T) {
	// Check this compiles and runs without panics
	NewMoveWithTrumpCardSwitch(NewCard(Ten, Hearts))
}

func TestNewMoveWithAnnouncementAndTrumpCardSwitch(t *testing.T) {
	// Check this compiles and runs without panics
	NewMoveWithAnnouncementAndTrumpCardSwitch(NewCard(Queen, Hearts))
}

func TestNewMoveWithAnnouncementAndTrumpCardSwitchImpossibleAnnouncement(t *testing.T) {
	assert.PanicsWithValue(
		t, "announcement moves are only possible with queens and kings",
		func() { NewMoveWithAnnouncementAndTrumpCardSwitch(NewCard(Ten, Hearts)) },
	)

}
