package santase

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

var suitStrings = map[Suit]string{
	Clubs:    "♣",
	Diamonds: "♦",
	Hearts:   "♥",
	Spades:   "♠",
}

func (s Suit) String() string {
	if str, ok := suitStrings[s]; ok {
		return str
	}
	return "invalid"
}

type Rank int

const (
	Nine Rank = iota
	Jack
	Queen
	King
	Ten
	Ace
)

var rankStrings = map[Rank]string{
	Nine:  "9",
	Jack:  "J",
	Queen: "Q",
	King:  "K",
	Ten:   "10",
	Ace:   "A",
}

func (r Rank) String() string {
	if str, ok := rankStrings[r]; ok {
		return str
	}
	return "invalid"
}

type Card struct {
	Suit Suit
	Rank Rank
}

func NewCard(rank Rank, suit Suit) Card {
	return Card{
		Rank: rank,
		Suit: suit,
	}
}

func (c Card) String() string {
	return c.Rank.String() + c.Suit.String()
}

var allCards []Card

func init() {
	for _, rank := range []Rank{Nine, Jack, Queen, King, Ten, Ace} {
		for _, suit := range []Suit{Clubs, Diamonds, Hearts, Spades} {
			allCards = append(allCards, NewCard(rank, suit))
		}
	}
}

type Hand map[Card]struct{}

func NewHand() Hand {
	return make(map[Card]struct{})
}

func (h *Hand) AddCard(c Card) {
	if len(*h) == 6 {
		panic("hand has 6 cards already")
	}

	if _, ok := (*h)[c]; ok {
		panic(c.String() + " is already in the hand")
	}

	(*h)[c] = struct{}{}
}

func (h *Hand) HasCard(c Card) bool {
	_, ok := (*h)[c]
	return ok
}

func (h *Hand) RemoveCard(c Card) {
	delete(*h, c)
}

type Move struct {
	Card            Card
	IsAnnouncement  bool
	SwitchTrumpCard bool
}

func NewMove(card Card) Move {
	return Move{
		Card: card,
	}
}

func NewMoveWithAnnouncement(card Card) Move {
	if card.Rank != Queen && card.Rank != King {
		panic("announcement moves are only possible with queens and kings")
	}

	return Move{
		Card:           card,
		IsAnnouncement: true,
	}
}

func NewMoveWithTrumpCardSwitch(card Card) Move {
	return Move{
		Card:            card,
		SwitchTrumpCard: true,
	}
}

func NewMoveWithAnnouncementAndTrumpCardSwitch(card Card) Move {
	if card.Rank != Queen && card.Rank != King {
		panic("announcement moves are only possible with queens and kings")
	}

	return Move{
		Card:            card,
		IsAnnouncement:  true,
		SwitchTrumpCard: true,
	}
}

func strongerCard(a *Card, b *Card, trump Suit) *Card {
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

func points(c *Card) int {
	if pts, ok := pointsMap[c.Rank]; ok {
		return pts
	}

	panic("invalid card")
}

func getRemainingCards(hand Hand, seenCards map[Card]struct{}) Hand {
	remaining := NewHand()
	for _, card := range allCards {
		isInHand := hand.HasCard(card)
		_, isSeen := seenCards[card]
		if !isInHand && !isSeen {
			remaining.AddCard(card)
		}
	}
	return remaining
}

type Game struct {
	score              int
	opponentScore      int
	seenCards          map[Card]struct{}
	knownOpponentCards Hand
	trump              Suit
	hand               Hand
	trumpCard          Card
	isOpponentMove     bool
	cardPlayed         *Card
}

func CreateGame(hand Hand, trumpCard Card, isOpponentMove bool) Game {
	if len(hand) != 6 {
		panic("player's hand is not complete")
	}

	return Game{
		score:              0,
		opponentScore:      0,
		seenCards:          make(map[Card]struct{}),
		knownOpponentCards: NewHand(),
		trump:              trumpCard.Suit,
		hand:               hand,
		trumpCard:          trumpCard,
		isOpponentMove:     isOpponentMove,
		cardPlayed:         nil,
	}
}

func (g *Game) getMove() Move {
	return Move{
		Card:            NewCard(Ace, Spades),
		IsAnnouncement:  false,
		SwitchTrumpCard: false,
	}
}

func (g *Game) GetMove() Move {
	if g.isOpponentMove {
		panic("not AI's turn")
	}

	if g.cardPlayed == nil && len(g.seenCards) <= 12 && len(g.hand) != 6 {
		panic("should not play before drawing cards")
	}

	move := g.getMove()

	if move.SwitchTrumpCard {
		g.hand.RemoveCard(NewCard(Nine, g.trump))
		g.hand.AddCard(g.trumpCard)
		g.trumpCard.Rank = Nine
	}

	if move.IsAnnouncement {
		if move.Card.Suit == g.trump {
			g.score += 40
		} else {
			g.score += 20
		}
	}

	g.hand.RemoveCard(move.Card)

	if g.cardPlayed == nil {
		g.cardPlayed = &move.Card
		g.isOpponentMove = true
	} else {
		stronger := strongerCard(g.cardPlayed, &move.Card, g.trump)
		if g.cardPlayed == stronger {
			g.score += points(g.cardPlayed) + points(&move.Card)
			g.isOpponentMove = true
		} else {
			g.opponentScore += points(g.cardPlayed) + points(&move.Card)
			g.isOpponentMove = false
		}
		g.seenCards[*g.cardPlayed] = struct{}{}
		g.seenCards[move.Card] = struct{}{}
		g.cardPlayed = nil
	}

	return move
}

func (g *Game) UpdateOpponentMove(opponentMove Move) {
	if !g.isOpponentMove {
		panic("not opponent's turn")
	}

	if _, ok := g.seenCards[opponentMove.Card]; ok {
		panic("card has already been played")
	}

	if _, ok := g.hand[opponentMove.Card]; ok {
		panic("card is in ai's hand")
	}

	if g.cardPlayed != nil && *g.cardPlayed == opponentMove.Card {
		panic("card is the same as the one on the table")
	}

	if g.cardPlayed == nil && len(g.seenCards) <= 12 && len(g.hand) != 6 {
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

		if len(g.seenCards) >= 12 {
			panic("cannot switch trump card after it has been taken")
		}

		if g.trumpCard.Rank == Nine {
			panic("cannot switch trump card - trump card is a nine")
		}

		g.knownOpponentCards.AddCard(g.trumpCard)
		g.trumpCard.Rank = Nine

		if g.knownOpponentCards.HasCard(g.trumpCard) {
			g.knownOpponentCards.RemoveCard(g.trumpCard)
		}
	}

	if opponentMove.Card == g.trumpCard && len(g.seenCards) < 12 {
		panic("played card is the trump card")
	}

	if g.knownOpponentCards.HasCard(opponentMove.Card) {
		g.knownOpponentCards.RemoveCard(opponentMove.Card)
	}

	if opponentMove.IsAnnouncement {
		if g.cardPlayed != nil {
			panic("cannot announce when you're not first to play")
		}

		if len(g.seenCards) == 0 {
			panic("cannot announce on first move")
		}

		var other Card
		if opponentMove.Card.Rank == Queen {
			other = NewCard(King, opponentMove.Card.Suit)
		} else {
			other = NewCard(Queen, opponentMove.Card.Suit)
		}

		if _, ok := g.seenCards[other]; ok {
			panic("cannot be an announcement because other card has already been played")
		}

		if _, ok := g.hand[other]; ok {
			panic("cannot be an announcement because other card is in ai's hand")
		}

		if other == g.trumpCard && len(g.seenCards) < 12 {
			panic("cannot be an announcement because other card is the trump card")
		}

		if opponentMove.Card.Suit == g.trump {
			g.opponentScore += 40
		} else {
			g.opponentScore += 20
		}

		if !g.knownOpponentCards.HasCard(other) {
			g.knownOpponentCards.AddCard(other)
		}
	}

	if g.cardPlayed == nil {
		g.cardPlayed = &opponentMove.Card
		g.isOpponentMove = false
	} else {
		stronger := strongerCard(g.cardPlayed, &opponentMove.Card, g.trump)
		if g.cardPlayed == stronger {
			g.score += points(g.cardPlayed) + points(&opponentMove.Card)
			g.isOpponentMove = false
		} else {
			g.opponentScore += points(g.cardPlayed) + points(&opponentMove.Card)
			g.isOpponentMove = true
		}
		g.seenCards[*g.cardPlayed] = struct{}{}
		g.seenCards[opponentMove.Card] = struct{}{}
		g.cardPlayed = nil
	}
}

func (g *Game) UpdateDrawnCard(card Card) {
	if g.cardPlayed != nil {
		panic("cannot draw cards in the middle of a play")
	}

	if len(g.hand) == 6 {
		if len(g.seenCards) == 0 {
			panic("should not draw cards before the first play")
		} else {
			panic("should not draw cards twice before playing")
		}
	}

	if _, ok := g.seenCards[card]; ok {
		panic("drawn card has been played before")
	}

	if g.knownOpponentCards.HasCard(card) {
		panic("cannot draw card that is in opponent's hand")
	}

	if g.hand.HasCard(card) {
		panic("cannot draw card that is in the hand already")
	}

	if g.trumpCard == card && len(g.seenCards) < 10 {
		panic("cannot draw trump card yet")
	}

	g.hand.AddCard(card)

	if len(g.seenCards) == 12 {
		g.knownOpponentCards = getRemainingCards(g.hand, g.seenCards)
	}
}
