// Package ismcts provides an agent implemented using
// Information Set Monte Carlo Tree Search (ISMCTS)
//
// ISMCTS[1] is an extenstion to the Monte Carlo Tree Search (MCTS)
// algorithm for games of imperfect information, such as santase where
// there are hidden cards.
//
// MCTS algorithms are anytime algorithms, meaning they can be ran for as
// long as they are allowed and still give good results. Generally running
// them for longer achieves better results.
//
// This package implements a parallelization technique on top of ISMCTS[2]
// that will start as many goroutines as there are cores on the machine.
//
// [1] Peter I. Cowling, Edward Powley and Daniel Whitehouse, “Information Set Monte Carlo Tree Search” http://orangehelicopter.com/academic/papers/tciaig_ismcts.pdf
//
// [2] Nick Sephton, Peter I. Cowling, Edward Powley, and Daniel Whitehouse, “Parallelization of Information Set Monte Carlo Tree Search” https://www-users.cs.york.ac.uk/~nsephton/papers/wcci2014-ismcts-parallelization.pdf
package ismcts

import (
	"math"
	"math/rand"
	"runtime"
	"time"

	santase "github.com/nvlbg/santase-ai"
)

type action struct {
	card      santase.Card
	closeGame bool
}

type node struct {
	parent       *node
	children     map[action]*node
	availability int
	visits       int
	score        int
}

func (n *node) isTerminal() bool {
	return len(n.children) == 0
}

func (n *node) isExpanded(game *game) bool {
	hand := game.getHand()
	canCloseGame := game.canClose()

	for card := range hand {
		isCardLegal := game.isCardLegal(card)
		if !isCardLegal {
			continue
		}

		// check if the move of playing this card is explored
		child := n.children[action{card: card}]
		if child == nil || child.visits == 0 {
			return false
		}

		// check if closing the game and playing the card is explored (if possible)
		if canCloseGame {
			child = n.children[action{card: card, closeGame: true}]
			if child == nil || child.visits == 0 {
				return false
			}
		}
	}

	if game.cardPlayed == nil && !game.isClosed && len(game.stack) > 1 && len(game.stack) < 11 {
		// check if switching is explored (if possible)
		nineTrump := santase.NewCard(santase.Nine, game.trump)
		child := n.children[action{card: *game.trumpCard}]
		if hand.HasCard(nineTrump) && (child == nil || child.visits == 0) {
			return false
		}

		// check if switching and closing the game is explored (if possible)
		child = n.children[action{card: *game.trumpCard, closeGame: true}]
		if hand.HasCard(nineTrump) && (child == nil || child.visits == 0) {
			return false
		}
	}

	return true
}

func (n *node) expandRandomChild(g *game) *node {
	hand := g.getHand()
	canCloseGame := g.canClose()

	var unexpandedActions []action
	for card := range hand {
		isLegal := g.isCardLegal(card)
		if !isLegal {
			continue
		}

		a := action{card: card}
		child := n.children[a]
		if child == nil || child.visits == 0 {
			unexpandedActions = append(unexpandedActions, a)
		}

		if canCloseGame {
			a := action{card: card, closeGame: true}
			child := n.children[a]
			if child == nil || child.visits == 0 {
				unexpandedActions = append(unexpandedActions, a)
			}
		}
	}

	if g.cardPlayed == nil && !g.isClosed && len(g.stack) > 1 && len(g.stack) < 11 {
		nineTrump := santase.NewCard(santase.Nine, g.trump)
		a := action{card: *g.trumpCard}
		child := n.children[a]
		if hand.HasCard(nineTrump) && (child == nil || child.visits == 0) {
			unexpandedActions = append(unexpandedActions, a)
		}

		a = action{card: *g.trumpCard, closeGame: true}
		child = n.children[a]
		if hand.HasCard(nineTrump) && (child == nil || child.visits == 0) {
			unexpandedActions = append(unexpandedActions, a)
		}
	}

	for _, a := range unexpandedActions {
		if n.children[a] == nil {
			n.children[a] = &node{
				parent:       n,
				children:     make(map[action]*node),
				availability: 1,
				visits:       0,
				score:        0,
			}
		}
	}
	action := unexpandedActions[rand.Intn(len(unexpandedActions))]
	g.simulate(action)
	n.children[action].visits++
	return n.children[action]
}

type game struct {
	score          int
	opponentScore  int
	hand           santase.Hand
	opponentHand   santase.Hand
	trump          santase.Suit
	stack          []santase.Card
	trumpCard      *santase.Card
	cardPlayed     *santase.Card
	isOpponentMove bool
	isClosed       bool
}

func (g *game) canClose() bool {
	return g.cardPlayed == nil && !g.isClosed && len(g.stack) > 1 && len(g.stack) < 11
}
func (g *game) getHand() santase.Hand {
	if g.isOpponentMove {
		return g.opponentHand
	}
	return g.hand
}

func (g *game) isCardLegal(card santase.Card) bool {
	// you're first to play or the game is not closed
	if g.cardPlayed == nil || (g.trumpCard != nil && !g.isClosed) {
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

func (g *game) simulate(a action) {
	hand := g.getHand()

	if g.cardPlayed == nil {
		// check if switching is possible
		if g.trumpCard != nil && !g.isClosed && g.trumpCard.Rank != santase.Nine && len(g.stack) > 1 && len(g.stack) < 11 {
			nineTrump := santase.NewCard(santase.Nine, g.trump)
			if a.card != nineTrump && hand.HasCard(nineTrump) {
				hand.RemoveCard(nineTrump)
				hand.AddCard(*g.trumpCard)
				g.trumpCard = &nineTrump
			}
		}

		if a.closeGame {
			g.isClosed = true
		}

		// check if announcing is possible
		if a.card.Rank == santase.Queen || a.card.Rank == santase.King && len(g.stack) < 11 {
			var other santase.Card
			if a.card.Rank == santase.Queen {
				other = santase.NewCard(santase.King, a.card.Suit)
			} else {
				other = santase.NewCard(santase.Queen, a.card.Suit)
			}

			if hand.HasCard(other) {
				var announcementPoints int
				if a.card.Suit == g.trump {
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

		g.cardPlayed = &a.card
		hand.RemoveCard(a.card)
		g.isOpponentMove = !g.isOpponentMove
	} else {
		stronger := santase.StrongerCard(g.cardPlayed, &a.card, g.trump)
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

		*winnerScore += santase.Points(g.cardPlayed) + santase.Points(&a.card)
		g.cardPlayed = nil
		hand.RemoveCard(a.card)

		if !g.isClosed {
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
}

func (g *game) runSimulation() int {
	var hand santase.Hand
	var a action
	for g.score < 66 && g.opponentScore < 66 && (len(g.hand) > 0 || len(g.opponentHand) > 0) {
		hand = g.getHand()

		if g.cardPlayed == nil {
			card := hand.GetRandomCard()
			// check if switching is possible
			if card == santase.NewCard(santase.Nine, g.trump) && !g.isClosed && len(g.stack) > 1 && len(g.stack) < 11 {
				// TODO: this way playing without switching is not simulated
				card = *g.trumpCard
			}

			a = action{card: card}
			// with probability = 1/7 decide wether to close the game at this turn
			if g.canClose() && rand.Intn(7) == 0 {
				a.closeGame = true
			}
		} else {
			if g.trumpCard != nil && !g.isClosed {
				a = action{card: hand.GetRandomCard()}
			} else {
				possibleResponses := hand.GetValidResponses(*g.cardPlayed, g.trump)
				a = action{card: possibleResponses.GetRandomCard()}
			}
		}
		g.simulate(a)
	}

	// TODO: 3 points here could potentially be only 2 if a player has a hand with 2 nines
	if g.score >= 66 && g.opponentScore >= 66 {
		if g.score > g.opponentScore {
			return 1
		} else if g.score < g.opponentScore {
			return -1
		} else {
			if g.isOpponentMove {
				return -1
			}
			return 1
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
		if g.isOpponentMove {
			return -1
		}
		return 1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sample(g *santase.Game) game {
	unseenCards := g.GetUnseenCards()
	hiddenCards := unseenCards.ToSlice()
	rand.Shuffle(len(hiddenCards), func(i, j int) {
		hiddenCards[i], hiddenCards[j] = hiddenCards[j], hiddenCards[i]
	})

	hand := g.GetHand()
	knownOpponentCards := g.GetKnownOpponentCards()

	splitAt := len(hand) - len(knownOpponentCards)
	if !g.IsOpponentMove() && g.GetCardPlayed() != nil {
		splitAt--
	}

	opponentHand := knownOpponentCards
	for _, card := range hiddenCards[:min(len(hiddenCards), splitAt)] {
		opponentHand.AddCard(card)
	}

	trumpCard := g.GetTrumpCard()
	var stack []santase.Card
	if trumpCard != nil {
		stack = hiddenCards[min(len(hiddenCards), splitAt):]
	}

	return game{
		score:          g.GetScore(),
		opponentScore:  g.GetOpponentScore(),
		hand:           hand,
		opponentHand:   opponentHand,
		trump:          g.GetTrump(),
		stack:          stack,
		trumpCard:      trumpCard,
		cardPlayed:     g.GetCardPlayed(),
		isOpponentMove: g.IsOpponentMove(),
		isClosed:       g.IsClosed(),
	}
}

func selectNode(root *node, game *game, c float64) *node {
	v := root

	for v.isExpanded(game) && !v.isTerminal() {
		// descend down the tree using modified UCB1
		bestScore := math.Inf(-1)
		var bestChild *node
		var bestAction action
		canCloseGame := game.canClose()
		for card := range game.getHand() {
			if !game.isCardLegal(card) {
				continue
			}
			a := action{card: card}
			u := v.children[a]

			f := float64(u.score) / float64(u.visits)
			if game.isOpponentMove {
				f *= -1
			}
			g := c * math.Sqrt(2*math.Log(float64(u.availability))/float64(u.visits))
			score := f + g
			if score > bestScore {
				bestScore = score
				bestChild = u
				bestAction = a
			}
			u.availability++

			if canCloseGame {
				a := action{card: card, closeGame: true}
				u = v.children[a]
				score = float64(u.score)/float64(u.visits) + c*math.Sqrt(2*math.Log(float64(u.availability))/float64(u.visits))
				if score > bestScore {
					bestScore = score
					bestChild = u
					bestAction = a
				}
				u.availability++
			}
		}

		v = bestChild
		v.visits++
		game.simulate(bestAction)
	}

	return v
}

func toMove(game *santase.Game, bestAction action) santase.Move {
	hand := game.GetHand()
	seenCards := game.GetSeenCards()
	cardPlayed := game.GetCardPlayed()

	// check if switching is possible
	switchTrumpCard := false
	if cardPlayed == nil && len(seenCards) > 0 && len(seenCards) < 10 {
		nineTrump := santase.NewCard(santase.Nine, game.GetTrump())
		if nineTrump != bestAction.card && hand.HasCard(nineTrump) {
			switchTrumpCard = true
		}
	}

	// check if announcing is possible
	isAnnouncement := false
	if cardPlayed == nil && len(seenCards) != 0 &&
		(bestAction.card.Rank == santase.Queen || bestAction.card.Rank == santase.King) {
		var other santase.Card
		if bestAction.card.Rank == santase.Queen {
			other = santase.NewCard(santase.King, bestAction.card.Suit)
		} else {
			other = santase.NewCard(santase.Queen, bestAction.card.Suit)
		}

		if hand.HasCard(other) || (switchTrumpCard && *game.GetTrumpCard() == other) {
			isAnnouncement = true
		}
	}

	return santase.Move{
		Card:            bestAction.card,
		SwitchTrumpCard: switchTrumpCard,
		IsAnnouncement:  isAnnouncement,
		CloseGame:       bestAction.closeGame,
	}
}

// SOISMCTS follows the pseudo code described in the paper
// "Information Set Monte Carlo Tree Search"
func (a *agent) SOISMCTS(game *santase.Game, results chan *node, quit chan struct{}) {
	root := node{children: make(map[action]*node)}

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
			v := selectNode(&root, &g, a.c)

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

	results <- &root
}

// singleObserverInformationSetMCTS is a single threaded ISMCTS
// implementation. It is equivalent with
// singleObserverInformationSetMCTSRootParallelization if ran
// on a machine with one cpu.
// This version can be easier to debug.
func (a *agent) singleObserverInformationSetMCTS(game *santase.Game) santase.Move {
	results := make(chan *node)
	quit := make(chan struct{})

	go func() {
		<-time.After(a.timePerMove)
		close(quit)
	}()

	go a.SOISMCTS(game, results, quit)
	root := <-results

	// return best move
	var bestAction action
	var maxVisits = 0
	for a, v := range root.children {
		if v.visits > maxVisits {
			maxVisits = v.visits
			bestAction = a
		}
	}

	return toMove(game, bestAction)
}

// singleObserverInformationSetMCTSRootParallelization implements
// ISMCTS with root parallelization as defined in the paper
// "Parallelization of Information Set Monte Carlo Tree Search"
// with as many workers as there are cores on the machine
func (a *agent) singleObserverInformationSetMCTSRootParallelization(game *santase.Game) santase.Move {
	results := make(chan *node)
	quit := make(chan struct{})

	go func() {
		<-time.After(a.timePerMove)
		close(quit)
	}()

	numCpus := runtime.NumCPU()

	for i := 0; i < numCpus; i++ {
		go a.SOISMCTS(game, results, quit)
	}

	stats := make(map[action]int)
	for i := 0; i < numCpus; i++ {
		root := <-results
		for a, v := range root.children {
			stats[a] += v.visits
		}
	}

	var bestAction action
	var maxVisits = 0
	for a, visits := range stats {
		if visits > maxVisits {
			maxVisits = visits
			bestAction = a
		}
	}

	return toMove(game, bestAction)
}

type agent struct {
	c           float64
	timePerMove time.Duration
}

func (a *agent) GetMove(game *santase.Game) santase.Move {
	return a.singleObserverInformationSetMCTSRootParallelization(game)
}

// NewAgent creates a new ISMCTS agent.
//
// The first parameter c is a constant used in the algorithm
// that balances exploitation and exploration
// (https://en.wikipedia.org/wiki/Monte_Carlo_tree_search#Exploration_and_exploitation).
// The choice of this parameter can affect playing strength.
// A value around 5.4 works good.
//
// The second parameter timePerMove chooses the maximum time
// per move the agent is allowed.
func NewAgent(c float64, timePerMove time.Duration) santase.Agent {
	return &agent{
		c:           5.4,
		timePerMove: 2 * time.Second,
	}
}
