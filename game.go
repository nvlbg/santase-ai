package santase

type dummyAgent struct{}

func (a dummyAgent) GetMove(g *Game) Move {
	panic("no agent provided")
}

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

func (g *Game) GetCardPlayed() *Card {
	if g.cardPlayed == nil {
		return nil
	}
	card := *g.cardPlayed
	return &card
}

func (g *Game) GetHand() Hand {
	return g.hand.Clone()
}

func (g *Game) IsClosed() bool {
	return g.isClosed
}

func (g *Game) IsOpponentMove() bool {
	return g.isOpponentMove
}

func (g *Game) GetKnownOpponentCards() Hand {
	return g.knownOpponentCards.Clone()
}

func (g *Game) GetScore() int {
	return g.score
}

func (g *Game) GetOpponentScore() int {
	return g.opponentScore
}

func (g *Game) GetTrump() Suit {
	return g.trump
}

func (g *Game) GetTrumpCard() *Card {
	if g.trumpCard == nil {
		return nil
	}
	card := *g.trumpCard
	return &card
}

func (g *Game) GetSeenCards() Pile {
	return g.seenCards.Clone()
}

func (g *Game) GetUnseenCards() Pile {
	return g.unseenCards.Clone()
}

func (g *Game) SetAgent(agent Agent) {
	g.agent = agent
}

func (g *Game) StrongerCard(a *Card, b *Card) *Card {
	return StrongerCard(a, b, g.trump)
}

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
