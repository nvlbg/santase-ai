package santase

import (
	"math/rand"
	"sort"
)

// Suit represents one of the four suits a card can have.
type Suit int

// Clubs(♣), Diamonds(♦), Hearts(♥) and Spades(♠).
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

// Rank represents one of the ranks a card can have.
type Rank int

// Nine(9), Jack(J), Queen(Q), King(K), Ten(10) and Ace(A).
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

// Card is represented by a suit and a rank.
type Card struct {
	Suit Suit
	Rank Rank
}

// NewCard initializes a new Card by its suit and rank.
func NewCard(rank Rank, suit Suit) Card {
	return Card{
		Rank: rank,
		Suit: suit,
	}
}

func (c Card) String() string {
	return c.Rank.String() + c.Suit.String()
}

// AllCards is a utility slice containing all valid cards in the game.
var AllCards = []Card{
	NewCard(Nine, Clubs),
	NewCard(Jack, Clubs),
	NewCard(Queen, Clubs),
	NewCard(King, Clubs),
	NewCard(Ten, Clubs),
	NewCard(Ace, Clubs),
	NewCard(Nine, Diamonds),
	NewCard(Jack, Diamonds),
	NewCard(Queen, Diamonds),
	NewCard(King, Diamonds),
	NewCard(Ten, Diamonds),
	NewCard(Ace, Diamonds),
	NewCard(Nine, Hearts),
	NewCard(Jack, Hearts),
	NewCard(Queen, Hearts),
	NewCard(King, Hearts),
	NewCard(Ten, Hearts),
	NewCard(Ace, Hearts),
	NewCard(Nine, Spades),
	NewCard(Jack, Spades),
	NewCard(Queen, Spades),
	NewCard(King, Spades),
	NewCard(Ten, Spades),
	NewCard(Ace, Spades),
}

// Hand represents a collection of up to 6 different cards that a player may hold.
type Hand map[Card]struct{}

// NewHand initializes a new Hand.
//
// The cards in the hand can be passed as variadic parameters
// to the function or can be later added with AddCard.
//
// If too many cards are passed a panic will occur.
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

// AddCard adds the passed card to the hand.
//
// If the card is already in the hand this is a noop.
//
// If the hand has 6 cards already a panic will occur.
func (h *Hand) AddCard(c Card) {
	if len(*h) == 6 {
		panic("hand has 6 cards already")
	}

	(*h)[c] = struct{}{}
}

// HasCard checks if a card is in the hand.
func (h *Hand) HasCard(c Card) bool {
	_, ok := (*h)[c]
	return ok
}

// RemoveCard removes a card from the hand.
//
// If the passed card is not in the hand this is a noop.
func (h *Hand) RemoveCard(c Card) {
	delete(*h, c)
}

// Clone returns a (deep) copy of the hand.
func (h *Hand) Clone() Hand {
	hand := NewHand()
	for card := range *h {
		hand.AddCard(card)
	}
	return hand
}

// ToSlice converts the hand to a slice of cards.
func (h *Hand) ToSlice() []Card {
	result := make([]Card, 0, len(*h))
	for card := range *h {
		result = append(result, card)
	}
	return result
}

// GetValidResponses returns a new hand containing only the
// cards that would be a valid response if the game is closed.
func (h *Hand) GetValidResponses(played Card, trump Suit) Hand {
	var allowedResponses []Card

	for card := range *h {
		if card.Suit == played.Suit && card.Rank > played.Rank {
			allowedResponses = append(allowedResponses, card)
		}
	}

	if allowedResponses == nil {
		for card := range *h {
			if card.Suit == played.Suit {
				allowedResponses = append(allowedResponses, card)
			}
		}
	}

	if allowedResponses == nil && played.Suit != trump {
		for card := range *h {
			if card.Suit == trump {
				allowedResponses = append(allowedResponses, card)
			}
		}
	}

	if allowedResponses == nil {
		return h.Clone()
	}

	return NewHand(allowedResponses...)
}

// GetRandomCard returns a card chosen at random from the hand.
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

// Pile represents a collection of (different) cards.
//
// Unlike Hand there is no constraint on the number of cards in a Pile
type Pile map[Card]struct{}

// NewPile initializes a new empty Pile.
func NewPile() Pile {
	return make(map[Card]struct{})
}

// AddCard adds the passed card to the pile.
//
// If the card is already in the pile this is a noop.
func (p *Pile) AddCard(c Card) {
	(*p)[c] = struct{}{}
}

// HasCard checks if a card is in the pile.
func (p *Pile) HasCard(c Card) bool {
	_, ok := (*p)[c]
	return ok
}

// RemoveCard removes a card from the pile.
//
// If the passed card is not in the pile this is a noop.
func (p *Pile) RemoveCard(c Card) {
	delete(*p, c)
}

// Clone returns a (deep) copy of the pile.
func (p *Pile) Clone() Pile {
	pile := NewPile()
	for card := range *p {
		pile.AddCard(card)
	}
	return pile
}

// ToSlice converts the pile to a slice of cards.
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

// Move contains the information for a move that a player.
// in the game chooses to play
//
// Each move has a card that is played (places on the table).
// Additionally a move can be an announcement (if the player
// has king and queen of matching suit), can switch the trump
// card (if valid) and can close the game (if valid).
type Move struct {
	Card            Card
	IsAnnouncement  bool
	SwitchTrumpCard bool
	CloseGame       bool
}

// Agent represents a player in the game and is used to
// to determine which move is to be played by the AI.
//
// Game uses such agents to choose what move the AI plays.
// The GetMove() method will be called when it is the AI's
// turn to play. A pointer to the game state will be passed
// so the agent can obtain relevant information needed to
// determine its move.
//
// Agent is also an extension point that can be used to
// create different bots.
//
// See packages in "github.com/nvlbg/santase-ai/agents" for examples of
// how agents can be implemented.
type Agent interface {
	GetMove(*Game) Move
}
