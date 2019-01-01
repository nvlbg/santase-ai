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

func createSampleHand() Hand {
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Diamonds))
	hand.AddCard(NewCard(King, Spades))
	hand.AddCard(NewCard(Queen, Diamonds))
	hand.AddCard(NewCard(Nine, Spades))
	hand.AddCard(NewCard(Ace, Spades))
	hand.AddCard(NewCard(Ten, Hearts))
	return hand
}

func createSampleGame() Game {
	hand := createSampleHand()
	trumpCard := NewCard(Ten, Clubs)
	return CreateGame(hand, trumpCard, true)
}

func TestUpdateOpponentMove(t *testing.T) {
	game := createSampleGame()

	card := NewCard(Ace, Diamonds)
	opponentMove := NewMove(card)
	game.UpdateOpponentMove(opponentMove)

	assert.Equal(t, card, *game.cardPlayed)
	assert.False(t, game.isOpponentMove)
}

func TestUpdateOpponentMoveInferingOpponentCards(t *testing.T) {
	t.Run("when announcing", func(t *testing.T) {
		game := createSampleGame()

		// simulate if one hand has been played already
		game.seenCards[NewCard(Queen, Diamonds)] = struct{}{}
		game.seenCards[NewCard(Ace, Diamonds)] = struct{}{}
		game.hand.RemoveCard(NewCard(Queen, Diamonds))
		game.opponentScore = 14
		game.UpdateDrawnCard(NewCard(Jack, Hearts))

		assert.False(t, game.knownOpponentCards.HasCard(NewCard(King, Hearts)))
		opponentMove := NewMoveWithAnnouncement(NewCard(Queen, Hearts))
		game.UpdateOpponentMove(opponentMove)
		assert.True(t, game.knownOpponentCards.HasCard(NewCard(King, Hearts)))
	})

	t.Run("after switching trump card", func(t *testing.T) {
		game := createSampleGame()

		// simulate if one hand has been played already
		game.seenCards[NewCard(Queen, Diamonds)] = struct{}{}
		game.seenCards[NewCard(Ace, Diamonds)] = struct{}{}
		game.hand.RemoveCard(NewCard(Queen, Diamonds))
		game.opponentScore = 14
		game.UpdateDrawnCard(NewCard(Jack, Hearts))

		originalTrumpCard := game.trumpCard
		assert.False(t, game.knownOpponentCards.HasCard(originalTrumpCard))
		opponentMove := NewMoveWithTrumpCardSwitch(NewCard(Queen, Hearts))
		game.UpdateOpponentMove(opponentMove)
		assert.True(t, game.knownOpponentCards.HasCard(originalTrumpCard))
	})

	t.Run("after switching trump card and announcing", func(t *testing.T) {
		game := createSampleGame()

		// simulate if one hand has been played already
		game.seenCards[NewCard(Queen, Diamonds)] = struct{}{}
		game.seenCards[NewCard(Ace, Diamonds)] = struct{}{}
		game.hand.RemoveCard(NewCard(Queen, Diamonds))
		game.opponentScore = 14
		game.UpdateDrawnCard(NewCard(Jack, Hearts))

		originalTrumpCard := game.trumpCard
		assert.False(t, game.knownOpponentCards.HasCard(originalTrumpCard))
		assert.False(t, game.knownOpponentCards.HasCard(NewCard(King, Hearts)))
		opponentMove := NewMoveWithAnnouncementAndTrumpCardSwitch(NewCard(Queen, Hearts))
		game.UpdateOpponentMove(opponentMove)
		assert.True(t, game.knownOpponentCards.HasCard(originalTrumpCard))
		assert.True(t, game.knownOpponentCards.HasCard(NewCard(King, Hearts)))
	})

	t.Run("after drawing all cards", func(t *testing.T) {
		// TODO
	})
}

func TestUpdateOpponentMoveInvalidSituations(t *testing.T) {
	card := NewCard(Ace, Diamonds)
	opponentMove := NewMove(card)

	t.Run("not opponents move", func(t *testing.T) {
		game := createSampleGame()
		game.isOpponentMove = false

		assert.PanicsWithValue(
			t, "not opponent's turn",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("card has been played already", func(t *testing.T) {
		game := createSampleGame()

		// simulate if one hand has been played already
		game.seenCards[NewCard(Queen, Diamonds)] = struct{}{}
		game.seenCards[card] = struct{}{}
		game.hand.RemoveCard(NewCard(Queen, Diamonds))
		game.opponentScore = 14

		assert.PanicsWithValue(
			t, "card has already been played",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("card is in our hand", func(t *testing.T) {
		game := createSampleGame()
		opponentMove := NewMove(NewCard(Nine, Diamonds))

		assert.PanicsWithValue(
			t, "card is in ai's hand",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("card is the one on the table", func(t *testing.T) {
		game := createSampleGame()

		// simulate if ai has played first move
		card := NewCard(Ace, Spades)
		game.cardPlayed = &card
		game.hand.RemoveCard(card)
		opponentMove := NewMove(card)

		assert.PanicsWithValue(
			t, "card is the same as the one on the table",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("playing before drawing", func(t *testing.T) {
		game := createSampleGame()

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14

		assert.PanicsWithValue(
			t, "should not play before drawing cards",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("playing trump card", func(t *testing.T) {
		game := createSampleGame()
		opponentMove := NewMove(game.trumpCard)

		assert.PanicsWithValue(
			t, "played card is the trump card",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("switching trump card when responding", func(t *testing.T) {
		game := createSampleGame()

		// simulate if ai has played first move
		card := NewCard(Ace, Spades)
		game.cardPlayed = &card
		game.hand.RemoveCard(card)
		opponentMove := NewMoveWithTrumpCardSwitch(opponentMove.Card)

		assert.PanicsWithValue(
			t, "cannot switch trump card when you're not first to play",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("switching trump card on first move", func(t *testing.T) {
		game := createSampleGame()
		opponentMove := NewMoveWithTrumpCardSwitch(opponentMove.Card)

		assert.PanicsWithValue(
			t, "cannot switch trump card on first move",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("switching trump card with only two cards left", func(t *testing.T) {
		// TODO
	})

	t.Run("switching trump card after all cards are taken", func(t *testing.T) {
		// TODO
	})

	t.Run("switching trump card with rank nine", func(t *testing.T) {
		game := createSampleGame()
		game.trumpCard.Rank = Nine
		opponentMove := NewMoveWithTrumpCardSwitch(opponentMove.Card)

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14
		game.UpdateDrawnCard(NewCard(Jack, Hearts))

		assert.PanicsWithValue(
			t, "cannot switch trump card - trump card is a nine",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("announcing in the middle of a play", func(t *testing.T) {
		game := createSampleGame()

		// simulate if ai has played first move
		card := NewCard(Ace, Spades)
		game.cardPlayed = &card
		game.hand.RemoveCard(card)

		opponentMove := NewMoveWithAnnouncement(NewCard(Queen, Clubs))

		assert.PanicsWithValue(
			t, "cannot announce when you're not first to play",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("announcing when the other card has already been played", func(t *testing.T) {
		game := createSampleGame()

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14
		game.UpdateDrawnCard(NewCard(Jack, Hearts))

		opponentMove := NewMoveWithAnnouncement(NewCard(Queen, Spades))
		assert.PanicsWithValue(
			t, "cannot be an announcement because other card has already been played",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("announcing on first move", func(t *testing.T) {
		game := createSampleGame()

		opponentMove := NewMoveWithAnnouncement(NewCard(Queen, Hearts))
		assert.PanicsWithValue(
			t, "cannot announce on first move",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("announcing when the other card is in ai's hand", func(t *testing.T) {
		game := createSampleGame()

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14
		game.UpdateDrawnCard(NewCard(Jack, Hearts))

		opponentMove := NewMoveWithAnnouncement(NewCard(King, Diamonds))
		assert.PanicsWithValue(
			t, "cannot be an announcement because other card is in ai's hand",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})

	t.Run("announcing when the other card is the trump card", func(t *testing.T) {
		game := createSampleGame()
		game.trumpCard.Rank = Queen

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14
		game.UpdateDrawnCard(NewCard(Jack, Hearts))

		opponentMove := NewMoveWithAnnouncement(NewCard(King, Clubs))
		assert.PanicsWithValue(
			t, "cannot be an announcement because other card is the trump card",
			func() { game.UpdateOpponentMove(opponentMove) },
		)
	})
}

func TestUpdateOpponentMoveEdgeCaseSituations(t *testing.T) {
	// TODO:
	// playing trump card after switching it
	// playing trump card after switching it + announcing
	// playing trump card after drawing it
}

func TestUpdateDrawnCardInvalidSituations(t *testing.T) {
	t.Run("in the middle of a play", func(t *testing.T) {
		game := createSampleGame()

		// simulate if ai has played first move
		card := NewCard(Ace, Spades)
		game.cardPlayed = &card
		game.hand.RemoveCard(card)

		assert.PanicsWithValue(
			t, "cannot draw cards in the middle of a play",
			func() { game.UpdateDrawnCard(NewCard(Jack, Hearts)) },
		)
	})

	t.Run("before first move", func(t *testing.T) {
		game := createSampleGame()

		assert.PanicsWithValue(
			t, "should not draw cards before the first play",
			func() { game.UpdateDrawnCard(NewCard(Jack, Hearts)) },
		)
	})

	t.Run("drawing twice in a row", func(t *testing.T) {
		game := createSampleGame()

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14
		game.UpdateDrawnCard(NewCard(Jack, Hearts))

		assert.PanicsWithValue(
			t, "should not draw cards twice before playing",
			func() { game.UpdateDrawnCard(NewCard(Nine, Hearts)) },
		)
	})

	t.Run("drawing seen card", func(t *testing.T) {
		game := createSampleGame()

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14

		assert.PanicsWithValue(
			t, "drawn card has been played before",
			func() { game.UpdateDrawnCard(NewCard(Ten, Spades)) },
		)
	})

	t.Run("drawing card in ai's hand", func(t *testing.T) {
		game := createSampleGame()

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14

		assert.PanicsWithValue(
			t, "cannot draw card that is in the hand already",
			func() { game.UpdateDrawnCard(NewCard(Ace, Spades)) },
		)
	})

	t.Run("drawing card in opponent's hand", func(t *testing.T) {
		// TODO
	})

	t.Run("drawing trump card", func(t *testing.T) {
		game := createSampleGame()

		// simulating playing one hand
		game.seenCards[NewCard(King, Spades)] = struct{}{}
		game.seenCards[NewCard(Ten, Spades)] = struct{}{}
		game.hand.RemoveCard(NewCard(King, Spades))
		game.opponentScore = 14

		assert.PanicsWithValue(
			t, "cannot draw trump card yet",
			func() { game.UpdateDrawnCard(game.trumpCard) },
		)
	})
}
