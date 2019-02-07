package random

import (
	"math/rand"

	santase "github.com/nvlbg/santase-ai"
)

type agent struct{}

func (a *agent) GetMove(game *santase.Game) santase.Move {
	hand := game.GetHand()
	cardPlayed := game.GetCardPlayed()
	if cardPlayed != nil && (game.IsClosed() || game.GetTrumpCard() == nil) {
		var allowedResponses []santase.Card
		for card := range hand {
			if card.Suit == cardPlayed.Suit && card.Rank > cardPlayed.Rank {
				allowedResponses = append(allowedResponses, card)
			}
		}

		if allowedResponses == nil {
			for card := range hand {
				if card.Suit == cardPlayed.Suit {
					allowedResponses = append(allowedResponses, card)
				}
			}
		}

		if allowedResponses == nil && cardPlayed.Suit != game.GetTrump() {
			for card := range hand {
				if card.Suit == game.GetTrump() {
					allowedResponses = append(allowedResponses, card)
				}
			}
		}

		if allowedResponses == nil {
			allowedResponses = hand.ToSlice()
		}

		return santase.Move{
			Card: allowedResponses[rand.Intn(len(allowedResponses))],
		}
	}

	return santase.Move{
		Card: hand.GetRandomCard(),
	}
}

func NewAgent() santase.Agent {
	return &agent{}
}
