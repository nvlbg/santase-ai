package santase

import (
	"math"
	"math/rand"
	"time"
)

type node struct {
	parent       *node
	children     map[Card]*node
	availability int
	visits       int
	score        int
}

func (n *node) isTerminal() bool {
	return len(n.children) == 0
}

func (n *node) isExpanded(game *game) bool {
	hand := game.getHand()

	for card := range hand {
		isCardLegal := game.isCardLegal(card)
		child := n.children[card]
		if isCardLegal && (child == nil || child.visits == 0) {
			return false
		}
	}

	if game.cardPlayed == nil && len(game.stack) > 1 && len(game.stack) < 11 {
		nineTrump := NewCard(Nine, game.trump)
		child := n.children[*game.trumpCard]
		if hand.HasCard(nineTrump) && (child == nil || child.visits == 0) {
			return false
		}
	}

	return true
}

func (n *node) expandRandomChild(g *game) *node {
	hand := g.getHand()

	var unexpandedCards []Card
	for card := range hand {
		isLegal := g.isCardLegal(card)
		child := n.children[card]
		if isLegal && (child == nil || child.visits == 0) {
			unexpandedCards = append(unexpandedCards, card)
		}
	}

	if g.cardPlayed == nil && len(g.stack) > 1 && len(g.stack) < 11 {
		nineTrump := NewCard(Nine, g.trump)
		child := n.children[*g.trumpCard]
		if hand.HasCard(nineTrump) && (child == nil || child.visits == 0) {
			unexpandedCards = append(unexpandedCards, *g.trumpCard)
		}
	}

	for _, card := range unexpandedCards {
		n.children[card] = &node{
			parent:       n,
			children:     make(map[Card]*node),
			availability: 1,
			visits:       0,
			score:        0,
		}
	}
	card := unexpandedCards[rand.Intn(len(unexpandedCards))]
	g.simulate(card)
	n.children[card].visits++
	return n.children[card]
}

type game struct {
	score          int
	opponentScore  int
	hand           Hand
	opponentHand   Hand
	trump          Suit
	stack          []Card
	trumpCard      *Card
	cardPlayed     *Card
	isOpponentMove bool
}

func (g *game) getHand() Hand {
	if g.isOpponentMove {
		return g.opponentHand
	}
	return g.hand
}

func (g *game) isCardLegal(card Card) bool {
	// you're first to play or the game is not closed
	if g.cardPlayed == nil || g.trumpCard != nil {
		return true
	}

	// playing stronger card of the requested suit
	if card.Suit == g.cardPlayed.Suit && card.Rank > g.cardPlayed.Rank {
		return true
	}

	hand := g.getHand()
	if g.cardPlayed.Suit == card.Suit {
		for c := range hand {
			if c.Suit == g.cardPlayed.Suit && c.Rank > g.cardPlayed.Rank {
				// you're holding stronger card of the same suit that you must play
				return false
			}
		}
		// you don't have stronger card of the same suit
		return true
	}

	for c := range hand {
		if c.Suit == g.cardPlayed.Suit {
			// you're holding card of the requested suit that you must play
			return false
		}
	}

	if g.cardPlayed.Suit != g.trump && card.Suit == g.trump {
		// you are forced to play trump card in this case
		return true
	}

	if g.cardPlayed.Suit != g.trump {
		for c := range hand {
			if c.Suit == g.trump {
				// you're holding a trump card that you should play
				return false
			}
		}
	}

	// your move is valid
	return true
}

func (g *game) simulate(card Card) {
	hand := g.getHand()

	if g.cardPlayed == nil {
		// check if switching is possible
		if g.trumpCard != nil && g.trumpCard.Rank != Nine && len(g.stack) > 1 && len(g.stack) < 11 {
			nineTrump := NewCard(Nine, g.trump)
			if card != nineTrump && hand.HasCard(nineTrump) {
				hand.RemoveCard(nineTrump)
				hand.AddCard(*g.trumpCard)
				g.trumpCard = &nineTrump
			}
		}

		// check if announcing is possible
		if card.Rank == Queen || card.Rank == King && len(g.stack) < 11 {
			var other Card
			if card.Rank == Queen {
				other = NewCard(King, card.Suit)
			} else {
				other = NewCard(Queen, card.Suit)
			}

			if hand.HasCard(other) {
				var announcementPoints int
				if card.Suit == g.trump {
					announcementPoints = 40
				} else {
					announcementPoints = 20
				}

				if g.isOpponentMove {
					g.opponentScore += announcementPoints
				} else {
					g.score += announcementPoints
				}
			}
		}

		g.cardPlayed = &card
		hand.RemoveCard(card)
		g.isOpponentMove = !g.isOpponentMove
	} else {
		stronger := strongerCard(g.cardPlayed, &card, g.trump)
		var winnerScore *int
		if g.cardPlayed == stronger {
			if g.isOpponentMove {
				winnerScore = &g.score
			} else {
				winnerScore = &g.opponentScore
			}

			g.isOpponentMove = !g.isOpponentMove
		} else {
			if g.isOpponentMove {
				winnerScore = &g.opponentScore
			} else {
				winnerScore = &g.score
			}
		}

		*winnerScore += Points(g.cardPlayed) + Points(&card)
		g.cardPlayed = nil
		hand.RemoveCard(card)

		if len(g.stack) > 1 {
			if g.isOpponentMove {
				g.opponentHand.AddCard(g.stack[len(g.stack)-1])
				g.hand.AddCard(g.stack[len(g.stack)-2])
			} else {
				g.hand.AddCard(g.stack[len(g.stack)-1])
				g.opponentHand.AddCard(g.stack[len(g.stack)-2])
			}
			g.stack = g.stack[:len(g.stack)-2]
		} else if len(g.stack) == 1 {
			if g.isOpponentMove {
				g.opponentHand.AddCard(g.stack[0])
				g.hand.AddCard(*g.trumpCard)
			} else {
				g.hand.AddCard(g.stack[0])
				g.opponentHand.AddCard(*g.trumpCard)
			}
			g.stack = nil
			g.trumpCard = nil
		}
	}
}

func (g *game) runSimulation() int {
	var hand Hand
	var card Card
	for g.score < 66 && g.opponentScore < 66 && (len(g.hand) > 0 || len(g.opponentHand) > 0) {
		if g.isOpponentMove {
			hand = g.opponentHand
		} else {
			hand = g.hand
		}

		if g.cardPlayed == nil {
			card = hand.GetRandomCard()
			// check if switching is possible
			if card == NewCard(Nine, g.trump) && len(g.stack) > 1 && len(g.stack) < 11 {
				// TODO: this way playing without switching is not simulated
				card = *g.trumpCard
			}
		} else {
			if g.trumpCard != nil {
				card = hand.GetRandomCard()
			} else {
				var possibleResponses []Card
				for card := range hand {
					if card.Suit == g.cardPlayed.Suit && card.Rank > g.cardPlayed.Rank {
						possibleResponses = append(possibleResponses, card)
					}
				}
				if possibleResponses == nil {
					for card := range hand {
						if card.Suit == g.cardPlayed.Suit {
							possibleResponses = append(possibleResponses, card)
						}
					}
				}
				if possibleResponses == nil && g.cardPlayed.Suit != g.trump {
					for card := range hand {
						if card.Suit == g.trump {
							possibleResponses = append(possibleResponses, card)
						}
					}
				}
				if possibleResponses == nil {
					for card := range hand {
						possibleResponses = append(possibleResponses, card)
					}
				}
				card = possibleResponses[rand.Intn(len(possibleResponses))]
			}
		}
		g.simulate(card)
	}

	// TODO: 3 points here could potentially be only 2 if a player has a hand with 2 nines
	if g.score >= 66 && g.opponentScore >= 66 {
		if g.score > g.opponentScore {
			return 1
		} else if g.score < g.opponentScore {
			return -1
		} else {
			return 0
		}
	} else if g.score >= 66 {
		if g.opponentScore == 0 {
			return 3
		} else if g.opponentScore < 33 {
			return 2
		} else {
			return 1
		}
	} else if g.opponentScore >= 66 {
		if g.score == 0 {
			return -3
		} else if g.score < 33 {
			return -2
		} else {
			return -1
		}
	} else {
		return 0
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sample(g *Game) game {
	hiddenCards := make([]Card, 0, len(g.unseenCards))
	for card := range g.unseenCards {
		hiddenCards = append(hiddenCards, card)
	}
	rand.Shuffle(len(hiddenCards), func(i, j int) {
		hiddenCards[i], hiddenCards[j] = hiddenCards[j], hiddenCards[i]
	})

	hand := NewHand()
	for card := range g.hand {
		hand.AddCard(card)
	}

	splitAt := 6 - len(g.knownOpponentCards)
	if !g.isOpponentMove && g.cardPlayed != nil {
		splitAt--
	}

	opponentHand := NewHand()
	for card := range g.knownOpponentCards {
		opponentHand.AddCard(card)
	}
	for _, card := range hiddenCards[:min(len(hiddenCards), splitAt)] {
		opponentHand.AddCard(card)
	}

	var stack []Card
	if g.trumpCard != nil {
		stack = hiddenCards[min(len(hiddenCards), splitAt):]
	}

	return game{
		score:          g.score,
		opponentScore:  g.opponentScore,
		hand:           hand,
		opponentHand:   opponentHand,
		trump:          g.trump,
		stack:          stack,
		trumpCard:      g.trumpCard,
		cardPlayed:     g.cardPlayed,
		isOpponentMove: g.isOpponentMove,
	}
}

func selectNode(root *node, game *game, c float64) *node {
	v := root

	for v.isExpanded(game) && !v.isTerminal() {
		// descend down the tree using modified UCB1
		bestScore := math.Inf(-1)
		var bestChild *node
		var bestCard Card
		for card := range game.getHand() {
			if !game.isCardLegal(card) {
				continue
			}
			u := v.children[card]
			score := float64(u.score)/float64(u.visits) + c*math.Sqrt(2*math.Log(float64(u.availability))/float64(u.visits))
			if score > bestScore {
				bestScore = score
				bestChild = u
				bestCard = card
			}

			u.availability++
		}

		v = bestChild
		v.visits++
		game.simulate(bestCard)
	}

	return v
}

func singleObserverInformationSetMCTS(game *Game) Move {
	root := node{children: make(map[Card]*node)}
	quit := make(chan struct{})

	go func() {
		<-time.After(game.timePerMove)
		close(quit)
	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		default:
			// choose a determinization at random compatible with the game
			// this iteration will use only actions compatible with the
			// selected determinization
			g := sample(game)

			// select which node to expand
			v := selectNode(&root, &g, game.c)

			// expand the tree if the selected node is not fully expanded
			if !v.isExpanded(&g) {
				v = v.expandRandomChild(&g)
			}

			// simulate the game till the end using random moves
			points := g.runSimulation()

			// backpropagation
			for v.parent != nil {
				v.score += points
				v = v.parent
			}
		}
	}

	// return best move
	var bestCard Card
	var maxVisits = 0
	for card, v := range root.children {
		if v.visits > maxVisits {
			maxVisits = v.visits
			bestCard = card
		}
	}

	// check if switching is possible
	switchTrumpCard := false
	if game.cardPlayed == nil && len(game.seenCards) > 0 && len(game.seenCards) < 10 {
		nineTrump := NewCard(Nine, game.trump)
		if nineTrump != bestCard && game.hand.HasCard(nineTrump) {
			switchTrumpCard = true
		}
	}

	// check if announcing is possible
	isAnnouncement := false
	if game.cardPlayed == nil && (bestCard.Rank == Queen || bestCard.Rank == King) {
		var other Card
		if bestCard.Rank == Queen {
			other = NewCard(King, bestCard.Suit)
		} else {
			other = NewCard(Queen, bestCard.Suit)
		}

		if game.hand.HasCard(other) || (switchTrumpCard && *game.trumpCard == other) {
			isAnnouncement = true
		}
	}

	return Move{
		Card:            bestCard,
		SwitchTrumpCard: switchTrumpCard,
		IsAnnouncement:  isAnnouncement,
	}
}
