package santase

import (
	"math/rand"
	"sort"
)

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

	(*h)[c] = struct{}{}
}

func (h *Hand) HasCard(c Card) bool {
	_, ok := (*h)[c]
	return ok
}

func (h *Hand) RemoveCard(c Card) {
	delete(*h, c)
}

func (h *Hand) GetRandomCard() Card {
	i := rand.Intn(len(*h))
	var card Card
	for card = range *h {
		if i == 0 {
			break
		}
		i--
	}
	return card
}

func (h Hand) String() string {
	var cards []Card
	for card := range h {
		cards = append(cards, card)
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Suit < cards[j].Suit || (cards[i].Suit == cards[j].Suit && cards[i].Rank < cards[j].Rank)
	})

	result := "{ "
	for _, card := range cards {
		result += card.String() + " "
	}
	result += "}"
	return result
}

type Pile map[Card]struct{}

func NewPile() Pile {
	return make(map[Card]struct{})
}

func (p *Pile) AddCard(c Card) {
	(*p)[c] = struct{}{}
}

func (p *Pile) HasCard(c Card) bool {
	_, ok := (*p)[c]
	return ok
}

func (p *Pile) RemoveCard(c Card) {
	delete(*p, c)
}

func (p Pile) String() string {
	result := "{ "
	for card := range p {
		result += card.String() + " "
	}
	result += "}"
	return result
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

func getHiddenCards(hand Hand, trumpCard Card) Pile {
	remaining := NewPile()
	for _, card := range allCards {
		isInHand := hand.HasCard(card)
		isTrumpCard := card == trumpCard
		if !isInHand && !isTrumpCard {
			remaining.AddCard(card)
		}
	}
	return remaining
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
	}
}

func (g *Game) getMove() Move {
	return singleObserverInformationSetMCTS(g)
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
		g.hand.AddCard(*g.trumpCard)
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

		if g.trumpCard == nil {
			panic("cannot switch trump card after it has been taken")
		}

		if g.trumpCard.Rank == Nine {
			panic("cannot switch trump card - trump card is a nine")
		}

		g.knownOpponentCards.AddCard(*g.trumpCard)
		g.trumpCard.Rank = Nine
		g.knownOpponentCards.RemoveCard(*g.trumpCard)
		g.unseenCards.RemoveCard(*g.trumpCard)
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
		stronger := strongerCard(g.cardPlayed, &opponentMove.Card, g.trump)
		if g.cardPlayed == stronger {
			g.score += points(g.cardPlayed) + points(&opponentMove.Card)
			g.isOpponentMove = false
		} else {
			g.opponentScore += points(g.cardPlayed) + points(&opponentMove.Card)
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
		g.unseenCards = NewPile()
	}
}
