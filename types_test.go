package santase

import (
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
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("code did not panic")
		}
	}()

	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))
	hand.AddCard(NewCard(Jack, Hearts))
	hand.AddCard(NewCard(Queen, Hearts))
	hand.AddCard(NewCard(King, Hearts))
	hand.AddCard(NewCard(Ten, Hearts))
	hand.AddCard(NewCard(Ace, Hearts))

	// adding 7th card panics
	hand.AddCard(NewCard(Ace, Spades))
}

func TestNewHandAddingSameCardTwice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("code did not panic")
		}
	}()

	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))
	hand.AddCard(NewCard(Nine, Hearts))
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
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("code did not panic")
		}
	}()

	NewMoveWithAnnouncement(NewCard(Nine, Clubs))
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
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("code did not panic")
		}
	}()

	NewMoveWithAnnouncementAndTrumpCardSwitch(NewCard(Ten, Hearts))
}
