santase-ai [![GoDoc](https://godoc.org/github.com/nvlbg/santase-ai?status.svg)](https://godoc.org/github.com/nvlbg/santase-ai)
==========

Santase-ai is an interface for different artificial
intelligence agents for the game santase (also known as
[sixty-six](https://en.wikipedia.org/wiki/Sixty-Six_(card_game\))).
It is useful in two cases:

1. You need out of the box artificial intelligence agent for the game
of santase. You can use the interface santase-ai provides and you can
switch to different implementations easily.
2. You want to create an artificial intelligence for the game of santase.
You can use santase-ai to skip writing common logic such as checking for
user input and focus on writing the AI.

Agents
------
This project includes two implementations for such agents.

### Random agent [![GoDoc](https://godoc.org/github.com/nvlbg/santase-ai/agents/random?status.svg)](https://godoc.org/github.com/nvlbg/santase-ai/agents/random)
Always plays random valid moves. This is more for demonstration purposes
than actually useful.

### Information Set Monte Carlo Tree Search agent [![GoDoc](https://godoc.org/github.com/nvlbg/santase-ai/agents/ismcts?status.svg)](https://godoc.org/github.com/nvlbg/santase-ai/agents/ismcts)
This is more advanced agent that uses monte carlo methods to search for
good moves. If you are interested of how this works I recommend looking into
[[1]](http://mcts.ai/)
[[2]](http://orangehelicopter.com/academic/papers/tciaig_ismcts.pdf)
[[3]](https://www-users.cs.york.ac.uk/~nsephton/papers/wcci2014-ismcts-parallelization.pdf).

Usage
-----
Here is how to use this library if you want to use an AI out of the box:

```go
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
```

To create your own agent all you need to do is implement the `Agent` interface:
```go
type Agent interface {
	GetMove(*Game) Move
}
```

Your agent will be called whenever it is time to play a move with a reference
to the game state, from which information about the game can be obtained. You
can see how the two agents are implemented for more information. The
[random agent](https://github.com/nvlbg/santase-ai/blob/master/agents/random/agent.go)
is pretty simple.

santase-gui
-----------
[santase-gui](https://github.com/nvlbg/santase-gui/) is a graphical interface
that uses santase-ai for the artificial intelligence part. You can use it to test
how your AI compares to different ones.

License
-------
This library is licensed under the MIT License.
