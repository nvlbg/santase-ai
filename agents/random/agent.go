// Package random provides a random agent that always plays random valid move.
package random

import (
	santase "github.com/nvlbg/santase-ai"
)

type agent struct{}

func (a *agent) GetMove(game *santase.Game) santase.Move {
	hand := game.GetHand()
	cardPlayed := game.GetCardPlayed()
	if cardPlayed != nil && (game.IsClosed() || game.GetTrumpCard() == nil) {
		hand = hand.GetValidResponses(*cardPlayed, game.GetTrump())
	}

	return santase.Move{
		Card: hand.GetRandomCard(),
	}
}

// NewAgent creates a new random agent.
func NewAgent() santase.Agent {
	return &agent{}
}
