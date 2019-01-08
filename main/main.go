package main

import (
	"fmt"
	"santase"
)

func main() {
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

	// update the game with the move the opponent makes
	game.UpdateOpponentMove(santase.NewMove(santase.NewCard(santase.Nine, santase.Hearts)))

	// start the AI
	move := game.GetMove()

	fmt.Println(move)

	// finish the first round by updating what card the AI draws
	game.UpdateDrawnCard(santase.NewCard(santase.Jack, santase.Hearts))
	// Output:
	//
}