package santase

import "fmt"

func ExampleGame() {
	// create initial hand for the ai
	hand := NewHand()
	hand.AddCard(NewCard(Nine, Diamonds))
	hand.AddCard(NewCard(King, Spades))
	hand.AddCard(NewCard(Queen, Diamonds))
	hand.AddCard(NewCard(Nine, Spades))
	hand.AddCard(NewCard(Ace, Spades))
	hand.AddCard(NewCard(Ten, Hearts))

	// create trump card for the game
	trumpCard := NewCard(Ten, Clubs)

	// is the opponent first to move
	isOpponentMove := true

	// create a game
	game := CreateGame(hand, trumpCard, isOpponentMove)

	// update the game with the move the opponent makes
	game.UpdateOpponentMove(NewMove(NewCard(Nine, Hearts)))

	// start the AI
	move := game.GetMove()

	fmt.Println(move.Card)

	// finish the first round by updating what card the AI draws
	game.UpdateDrawnCard(NewCard(Jack, Hearts))
	// Output:
	// 10♥
}
