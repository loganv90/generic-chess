package bot

/*
Chess Search:
Alpha beta search
Iterative deepening
Transposition table (repeated and colour flipped positions)
Keep track of killer moves (moves that caused cut off) either by [piece][to] or [from][to]
Quiescence search (do all captures on last capture/promotion square)
Move ordering (Hashed moves, winning captures/promotions, equal captures, killer moves, non-captures, losing captures)
Create a hash function to uniquely identify the chess position
q2k2q1/2nqn2b/1n1P1n1b/2rnr2Q/1NQ1QN1Q/3Q3B/2RQR2B/Q2K2Q1 w - - this position causes engines to explode

Chess Evaluation:
Material count
Piece mobility with move bonus for different pieces (including moves protecting ally pieces, and not including moves to controlled squares) low queen bonus
Piece locations piece-square tables (dp for knight 2 moves potential squares) tables for opening and endgame, then interpolate
Penalty isolated pawns worth less than chained pawns
Penalty for attacked squares close to king
Penalty for lots of mobility if king were a queen
Bonus for attacking close to own king
Bonus for pinning pieces to more valuable pieces
Bonus for queen-rook, queen-bishop, bishop-bishop, rook-rook combos
Consider lazy evaluations if some eval rules already exceed alpha or beta by some significant margin

How:
We need to connect the bot to the game such that the UI still receives the same actions
We can have a channel for the UI to send actions to the bot
We can have a channel for the bot to send actions to the UI

Ok so the bot will have a run function that will run in a goroutine
There will be two channels for the bot to communicate with the Hub
When the bot is created, it will get a reference to the game

When the bot is created:
- create a channel called send in the bot that it will use to receive actions from the Hub
- the bot can just use the existing send channel in the hub to send actions to the hub
- we're going to have to create a bot-client class to hold the communicate stuff for the hub

When the bot is running:
- the bot will be told which colors it is playing
- the bot will be pinged by the hub every time a move is made
- the bot will make a calculation and respond back to the hub

Changes to the UI:
- the UI will need a control to specify the colors the bot is playing
*/

/*
Responsible for:
- Evaluating a board and returning a score
*/
type Evaluator interface {
    Evaluate() (int, error) // TODO add board as parameter
}

func newSimpleEvaluator() (*SimpleEvaluator, error) {
    return &SimpleEvaluator{}, nil
}

type SimpleEvaluator struct {
    score int
}

func (s *SimpleEvaluator) Evaluate() int {
    return s.score
}

