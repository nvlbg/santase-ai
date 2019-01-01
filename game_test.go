package santase

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGame(t *testing.T) {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Hearts))
	hand.AddCard(NewCard(Jack, Hearts))
	hand.AddCard(NewCard(Queen, Hearts))
	hand.AddCard(NewCard(King, Hearts))
	hand.AddCard(NewCard(Ten, Hearts))
	hand.AddCard(NewCard(Ace, Hearts))

	trumpCard := NewCard(Ace, Spades)

	// check that it compiles and runs without panics
	CreateGame(hand, trumpCard, false)
}

func TestNewGameIncompleteHand(t *testing.T) {
	hand := NewHand()
	trumpCard := NewCard(Ace, Spades)
	assert.PanicsWithValue(
		t, "player's hand is not complete",
		func() { CreateGame(hand, trumpCard, false) },
	)
}

func createInitialHand() Hand {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Diamonds))
	hand.AddCard(NewCard(King, Spades))
	hand.AddCard(NewCard(Queen, Diamonds))
	hand.AddCard(NewCard(Nine, Spades))
	hand.AddCard(NewCard(Ace, Spades))
	hand.AddCard(NewCard(Ten, Hearts))
	return hand
}

func TestUpdateOpponentMove(t *testing.T) {
	hand := createInitialHand()
	trumpCard := NewCard(Ten, Clubs)
	game := CreateGame(hand, trumpCard, true)

	opponentMove := NewMove(NewCard(Ace, Diamonds))
	game.UpdateOpponentMove(opponentMove)
}

func TestUpdateOpponentMoveWrongTurn(t *testing.T) {
	hand := createInitialHand()
	trumpCard := NewCard(Ten, Clubs)
	game := CreateGame(hand, trumpCard, false)

	opponentMove := NewMove(NewCard(Ace, Diamonds))

	assert.PanicsWithValue(
		t, "not opponent's turn",
		func() { game.UpdateOpponentMove(opponentMove) },
	)
}

func TestUpdateOpponentMoveWithCardInOurHand(t *testing.T) {
	hand := createInitialHand()
	trumpCard := NewCard(Ten, Clubs)
	game := CreateGame(hand, trumpCard, true)

	// played card is in ai's hand
	opponentMove := NewMove(NewCard(Nine, Diamonds))

	assert.PanicsWithValue(
		t, "card is in ai's hand",
		func() { game.UpdateOpponentMove(opponentMove) },
	)
}
