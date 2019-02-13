package santase_test

import (
	"fmt"
	"time"

	santase "github.com/nvlbg/santase-ai"
	"github.com/nvlbg/santase-ai/agents/ismcts"
)

func ExampleGame() {
	// create initial hand for the ai
	hand := santase.NewHand()
	hand.AddCard(santase.NewCard(santase.Nine, santase.Diamonds))
	hand.AddCard(santase.NewCard(santase.King, santase.Spades))
	hand.AddCard(santase.NewCard(santase.Queen, santase.Diamonds))
	hand.AddCard(santase.NewCard(santase.Nine, santase.Spades))
	hand.AddCard(santase.NewCard(santase.Ace, santase.Spades))
	hand.AddCard(santase.NewCard(santase.Ten, santase.Hearts))

	// create trump card for the game
	trumpCard := santase.NewCard(santase.Ten, santase.Clubs)

	// is the opponent first to move
	isOpponentMove := true

	// create a game
	game := santase.CreateGame(hand, trumpCard, isOpponentMove)

	// specify which agent to use for choosing moves
	game.SetAgent(ismcts.NewAgent(5.4, time.Second))

	// update the game with the move the opponent makes
	game.UpdateOpponentMove(santase.Move{Card: santase.NewCard(santase.Nine, santase.Hearts)})

	// start the AI
	move := game.GetMove()

	fmt.Println(move.Card)

	// finish the first round by updating what card the AI draws
	game.UpdateDrawnCard(santase.NewCard(santase.Jack, santase.Hearts))
	// Output:
	// 10♥
}

func ExampleNewHand() {
	// constructs new empty hand
	_ = santase.NewHand()

	// constructs new hand with two cards
	hand := santase.NewHand(
		santase.NewCard(santase.Queen, santase.Hearts),
		santase.NewCard(santase.King, santase.Hearts),
	)

	// the cards in a hand can be iterated like so
	for card := range hand {
		fmt.Println(card)
	}

	// Unordered output:
	// Q♥
	// K♥
}

func ExampleHand_GetValidResponses() {
	hand := santase.NewHand(
		santase.NewCard(santase.Nine, santase.Diamonds),
		santase.NewCard(santase.King, santase.Spades),
		santase.NewCard(santase.Queen, santase.Diamonds),
		santase.NewCard(santase.Nine, santase.Spades),
		santase.NewCard(santase.Ace, santase.Spades),
		santase.NewCard(santase.Ten, santase.Hearts),
	)

	valid := hand.GetValidResponses(
		santase.NewCard(santase.Queen, santase.Spades),
		santase.Hearts,
	)

	fmt.Println(valid)

	// Output:
	// { K♠ A♠ }
}

func ExampleNewPile() {
	// constructs new empty pile
	pile := santase.NewPile()

	pile.AddCard(santase.NewCard(santase.Queen, santase.Hearts))
	pile.AddCard(santase.NewCard(santase.King, santase.Hearts))

	// the cards in a pile can be iterated like so
	for card := range pile {
		fmt.Println(card)
	}

	// Unordered output:
	// Q♥
	// K♥
}
