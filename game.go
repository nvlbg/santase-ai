package santase

type dummyAgent struct{}

func (a dummyAgent) GetMove(g *Game) Move {
	panic("no agent provided")
}

// Game is an instance of a game of santase.  In its state
// in contains information for only what a player in the game
// can see and deduce.
type Game struct {
	trump              Suit
	score              int
	opponentScore      int
	hand               Hand
	knownOpponentCards Hand
	seenCards          Pile
	unseenCards        Pile
	trumpCard          *Card
	cardPlayed         *Card
	isOpponentMove     bool
	isClosed           bool
	agent              Agent
}

// CreateGame creates a new instance of a Game given the
// initial hand for the AI, the trump card on the table and
// whether the opponent (from the point of view of the AI)
// plays first.
//
// Panics if the hand does not have 6 cards.
func CreateGame(hand Hand, trumpCard Card, isOpponentMove bool) Game {
	if len(hand) != 6 {
		panic("player's hand is not complete")
	}

	return Game{
		trump:              trumpCard.Suit,
		score:              0,
		opponentScore:      0,
		hand:               hand,
		knownOpponentCards: NewHand(),
		seenCards:          NewPile(),
		unseenCards:        getHiddenCards(hand, trumpCard),
		trumpCard:          &trumpCard,
		cardPlayed:         nil,
		isOpponentMove:     isOpponentMove,
		isClosed:           false,
		agent:              dummyAgent{},
	}
}

// GetCardPlayed returns a pointer to the card placed on the
// table by one of the players. If there is no card played
// the result will be nil.
func (g *Game) GetCardPlayed() *Card {
	if g.cardPlayed == nil {
		return nil
	}
	card := *g.cardPlayed
	return &card
}

// GetHand returns the hand of the AI player.
func (g *Game) GetHand() Hand {
	return g.hand.Clone()
}

// IsClosed returns if the game has been explicitly closed
// by one of the players.
func (g *Game) IsClosed() bool {
	return g.isClosed
}

// IsOpponentMove returns if it is turn for the opponent
// (from the point of view of the AI) to play next.
func (g *Game) IsOpponentMove() bool {
	return g.isOpponentMove
}

// GetKnownOpponentCards returns the cards in the opponent's
// hand that have been deduced.
func (g *Game) GetKnownOpponentCards() Hand {
	return g.knownOpponentCards.Clone()
}

// GetScore returns the points that the AI player has collected.
func (g *Game) GetScore() int {
	return g.score
}

// GetOpponentScore returns the points that the opponent has
// collected.
func (g *Game) GetOpponentScore() int {
	return g.opponentScore
}

// GetTrump returns the trump suit of the game.
func (g *Game) GetTrump() Suit {
	return g.trump
}

// GetTrumpCard returns a pointer to the trump card placed on the table.
// If all cards have been drawn the result will be nil.
func (g *Game) GetTrumpCard() *Card {
	if g.trumpCard == nil {
		return nil
	}
	card := *g.trumpCard
	return &card
}

// GetSeenCards returns the set of cards that have been taken by one of
// the players and cannot be played anymore in the game. If there is a
// card placed on the table it will not be included here (see GetCardPlayed).
func (g *Game) GetSeenCards() Pile {
	return g.seenCards.Clone()
}

// GetUnseenCards returns the set of cards that are not yet deretmined
// if they are in the opponent's hand or in the pile.
func (g *Game) GetUnseenCards() Pile {
	return g.unseenCards.Clone()
}

// SetAgent sets the Agent that will be used to choose the moves that
// will be played by the AI. There is a random agent and a monte carlo
// agent included in the library that can be used, or you can write your own.
func (g *Game) SetAgent(agent Agent) {
	g.agent = agent
}

// StrongerCard returns which of the two cards will be stronger in the
// context of the game.
func (g *Game) StrongerCard(first *Card, second *Card) *Card {
	return StrongerCard(first, second, g.trump)
}

// GetMove returns the move that the AI agent chose to play. It should be
// called only when it is the AI's turn to play, otherwise a panic will occur.
// If there is a bug in the agent and it chooses an invalid move a panic will
// occur as well.
//
// Note: the order of calls to GetMove, UpdateOpponentMove and
// UpdateDrawnCard matters.
func (g *Game) GetMove() Move {
	if g.isOpponentMove {
		panic("not AI's turn")
	}

	if !g.isClosed && g.cardPlayed == nil && len(g.seenCards) <= 12 && len(g.hand) != 6 {
		panic("should not play before drawing cards")
	}

	move := g.agent.GetMove(g)

	if move.SwitchTrumpCard {
		if g.cardPlayed != nil {
			panic("cannot switch trump card when you're not first to play")
		}

		if len(g.seenCards) == 0 {
			panic("cannot switch trump card on first move")
		}

		if len(g.seenCards) == 10 {
			panic("cannot switch trump card with only two cards left in the stack")
		}

		if g.trumpCard == nil {
			panic("cannot switch trump card after it has been taken")
		}

		if g.isClosed {
			panic("cannot switch trump card after the game has been closed")
		}

		if !g.hand.HasCard(NewCard(Nine, g.trump)) {
			panic("cannot switch trump card withouth nine of trump in hand")
		}

		g.hand.RemoveCard(NewCard(Nine, g.trump))
		g.hand.AddCard(*g.trumpCard)
		g.trumpCard.Rank = Nine
	}

	if move.CloseGame {
		if g.cardPlayed != nil {
			panic("cannot close game when second to move")
		}

		if len(g.seenCards) == 0 {
			panic("cannot close game on first move")
		}

		if len(g.seenCards) == 10 {
			panic("cannot close game with only two cards left in the stack")
		}

		if len(g.seenCards) >= 12 {
			panic("cannot close game after all cards have been drawn")
		}

		if g.isClosed {
			panic("cannot close game because it is already closed")
		}

		g.isClosed = true
	}

	if move.IsAnnouncement {
		if g.cardPlayed != nil {
			panic("cannot announce when you're not first to play")
		}

		if len(g.seenCards) == 0 {
			panic("cannot announce on first move")
		}

		if move.Card.Rank != Queen && move.Card.Rank != King {
			panic("invalid announcement card")
		}

		var other Card
		if move.Card.Rank == Queen {
			other = NewCard(King, move.Card.Suit)
		} else {
			other = NewCard(Queen, move.Card.Suit)
		}

		if !g.hand.HasCard(other) {
			panic("invalid announcement - not both cards of announcement are in hand")
		}

		if move.Card.Suit == g.trump {
			g.score += 40
		} else {
			g.score += 20
		}
	}

	if !g.hand.HasCard(move.Card) {
		panic("played card in not in hand")
	}

	if g.cardPlayed != nil && (g.isClosed || g.trumpCard == nil) {
		possibleResponses := g.hand.GetValidResponses(*g.cardPlayed, g.trump)
		if !possibleResponses.HasCard(move.Card) {
			panic("invalid response card: " + move.Card.String())
		}
	}

	g.hand.RemoveCard(move.Card)

	if g.cardPlayed == nil {
		g.cardPlayed = &move.Card
		g.isOpponentMove = true
	} else {
		stronger := StrongerCard(g.cardPlayed, &move.Card, g.trump)
		if g.cardPlayed == stronger {
			g.opponentScore += Points(g.cardPlayed) + Points(&move.Card)
			g.isOpponentMove = true
		} else {
			g.score += Points(g.cardPlayed) + Points(&move.Card)
			g.isOpponentMove = false
		}
		g.seenCards.AddCard(*g.cardPlayed)
		g.seenCards.AddCard(move.Card)
		g.cardPlayed = nil
	}

	return move
}

// UpdateOpponentMove updates the game state with the move that the opponent
// has played. If the move is invalid a panic will occur.
//
// Note: the order of calls to GetMove, UpdateOpponentMove and
// UpdateDrawnCard matters.
func (g *Game) UpdateOpponentMove(opponentMove Move) {
	if !g.isOpponentMove {
		panic("not opponent's turn")
	}

	if g.seenCards.HasCard(opponentMove.Card) {
		panic("card has already been played")
	}

	if g.hand.HasCard(opponentMove.Card) {
		panic("card is in ai's hand")
	}

	if g.cardPlayed != nil && *g.cardPlayed == opponentMove.Card {
		panic("card is the same as the one on the table")
	}

	if !g.isClosed && g.cardPlayed == nil && len(g.seenCards) <= 12 && len(g.hand) != 6 {
		panic("should not play before drawing cards")
	}

	if opponentMove.SwitchTrumpCard {
		if g.cardPlayed != nil {
			panic("cannot switch trump card when you're not first to play")
		}

		if len(g.seenCards) == 0 {
			panic("cannot switch trump card on first move")
		}

		if len(g.seenCards) == 10 {
			panic("cannot switch trump card with only two cards left in the stack")
		}

		if g.trumpCard == nil {
			panic("cannot switch trump card after it has been taken")
		}

		if g.isClosed {
			panic("cannot switch trump card after the game has been closed")
		}

		if g.trumpCard.Rank == Nine {
			panic("cannot switch trump card - trump card is a nine")
		}

		g.knownOpponentCards.AddCard(*g.trumpCard)
		g.trumpCard.Rank = Nine
		g.knownOpponentCards.RemoveCard(*g.trumpCard)
		g.unseenCards.RemoveCard(*g.trumpCard)
	}

	if opponentMove.CloseGame {
		if g.cardPlayed != nil {
			panic("cannot close game when second to move")
		}

		if len(g.seenCards) == 0 {
			panic("cannot close game on first move")
		}

		if len(g.seenCards) == 10 {
			panic("cannot close game with only two cards left in the stack")
		}

		if len(g.seenCards) >= 12 {
			panic("cannot close game after all cards have been drawn")
		}

		if g.isClosed {
			panic("cannot close game because it is already closed")
		}

		g.isClosed = true
	}

	if g.trumpCard != nil && opponentMove.Card == *g.trumpCard {
		panic("played card is the trump card")
	}

	g.knownOpponentCards.RemoveCard(opponentMove.Card)

	if opponentMove.IsAnnouncement {
		if g.cardPlayed != nil {
			panic("cannot announce when you're not first to play")
		}

		if len(g.seenCards) == 0 {
			panic("cannot announce on first move")
		}

		if opponentMove.Card.Rank != Queen && opponentMove.Card.Rank != King {
			panic("invalid announcement card: " + opponentMove.Card.String())
		}

		var other Card
		if opponentMove.Card.Rank == Queen {
			other = NewCard(King, opponentMove.Card.Suit)
		} else {
			other = NewCard(Queen, opponentMove.Card.Suit)
		}

		if g.seenCards.HasCard(other) {
			panic("cannot be an announcement because other card has already been played")
		}

		if g.hand.HasCard(other) {
			panic("cannot be an announcement because other card is in ai's hand")
		}

		if g.trumpCard != nil && other == *g.trumpCard {
			panic("cannot be an announcement because other card is the trump card")
		}

		if opponentMove.Card.Suit == g.trump {
			g.opponentScore += 40
		} else {
			g.opponentScore += 20
		}

		g.knownOpponentCards.AddCard(other)
		g.unseenCards.RemoveCard(other)
	}

	g.unseenCards.RemoveCard(opponentMove.Card)

	if g.cardPlayed == nil {
		g.cardPlayed = &opponentMove.Card
		g.isOpponentMove = false
	} else {
		stronger := StrongerCard(g.cardPlayed, &opponentMove.Card, g.trump)
		if g.cardPlayed == stronger {
			g.score += Points(g.cardPlayed) + Points(&opponentMove.Card)
			g.isOpponentMove = false
		} else {
			g.opponentScore += Points(g.cardPlayed) + Points(&opponentMove.Card)
			g.isOpponentMove = true
		}
		g.seenCards.AddCard(*g.cardPlayed)
		g.seenCards.AddCard(opponentMove.Card)
		g.cardPlayed = nil
	}
}

// UpdateDrawnCard updates the game state with the card that the AI player
// draws from the stack of cards after a move. If the drawn card is invalid
// or it is not the time to draw cards a panic will occur.
//
// Note: the order of calls to GetMove, UpdateOpponentMove and
// UpdateDrawnCard matters.
func (g *Game) UpdateDrawnCard(card Card) {
	if g.cardPlayed != nil {
		panic("cannot draw cards in the middle of a play")
	}

	if g.isClosed {
		panic("should not draw cards when the game is closed")
	}

	if len(g.hand) == 6 {
		if len(g.seenCards) == 0 {
			panic("should not draw cards before the first play")
		} else {
			panic("should not draw cards twice before playing")
		}
	}

	if g.seenCards.HasCard(card) {
		panic("drawn card has been played before")
	}

	if g.knownOpponentCards.HasCard(card) {
		panic("cannot draw card that is in opponent's hand")
	}

	if g.hand.HasCard(card) {
		panic("cannot draw card that is in the hand already")
	}

	if g.trumpCard == nil {
		panic("all cards are drawn already")
	}

	if *g.trumpCard == card && len(g.seenCards) < 10 {
		panic("cannot draw trump card yet")
	}

	g.hand.AddCard(card)
	g.unseenCards.RemoveCard(card)

	if len(g.seenCards) == 12 {
		for card := range g.unseenCards {
			g.knownOpponentCards.AddCard(card)
		}
		if card != *g.trumpCard {
			g.knownOpponentCards.AddCard(*g.trumpCard)
		}
		g.unseenCards = NewPile()
		g.trumpCard = nil
	}
}
