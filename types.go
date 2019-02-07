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

var AllCards []Card

func init() {
	for _, rank := range []Rank{Nine, Jack, Queen, King, Ten, Ace} {
		for _, suit := range []Suit{Clubs, Diamonds, Hearts, Spades} {
			AllCards = append(AllCards, NewCard(rank, suit))
		}
	}
}

type Hand map[Card]struct{}

func NewHand(cards ...Card) Hand {
	if len(cards) > 6 {
		panic("too many cards given")
	}

	result := make(map[Card]struct{})
	for _, card := range cards {
		result[card] = struct{}{}
	}
	return result
}

func (h *Hand) Clone() Hand {
	hand := NewHand()
	for card := range *h {
		hand.AddCard(card)
	}
	return hand
}

func (h *Hand) ToSlice() []Card {
	result := make([]Card, 0, len(*h))
	for card := range *h {
		result = append(result, card)
	}
	return result
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

func (p *Pile) Clone() Pile {
	pile := NewPile()
	for card := range *p {
		pile.AddCard(card)
	}
	return pile
}

func (p *Pile) ToSlice() []Card {
	result := make([]Card, 0, len(*p))
	for card := range *p {
		result = append(result, card)
	}
	return result
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
	CloseGame       bool
}

type Agent interface {
	GetMove(*Game) Move
}
