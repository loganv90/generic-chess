# Dev Journal

## Initial Commits
- Originally, I created the chess logic in the frontend using Javascript. This included the logic for piece movements, the logic for the chess game, and unit tests.
- I re-wrote the chess logic using the command pattern to allow undoing and redoing chess moves. The command pattern requires: commands with undo and redo methods, a receiver that the commands will act upon, and an invoker that keeps track of which commands have been executed. In our case, the receiver is the chess board object, and the commands are the chess move objects. The invoker contains references to the moves, and the moves contain references to the board.
- Then, I re-wrote the chess logic in Typescript because I enjoy types.

## Creating a Go Backend
- I wrote the chess logic additionally in Go. The plan was to have logic on the frontend and backend to minimize requests and maximize responsiveness, but, eventually, the frontend logic was removed.
- I adding functional tests to complement the unit tests. The functional tests simulated chess games while verifying the board state.
- I connected the Typescript frontend to the Go backend using websockets. I added logic so the server could handle many concurrent chess games with many concurrent client connections. This required having client objects to wrap the websocket connections, hub objects to keep track of the clients in each game room, and a server to keep track of the hubs.
- I added websocket APIs to communicate between the client and server. Clients send http requests to connect to games and to establish websocket connections, but, after that, the chess game communication solely uses websockets.

## Adding Chess Variations
- When I wrote the chess logic, I purposely kept everything generic enough to support chess games with arbitrary board configurations and player counts. I configured a few chess variations and made them accessible in the UI.
- I made use of the command pattern again to enable player transitions. Player transition objects were added to the invoker and would: make changes to a player collection object to transition between turns, and make changes to a chess board object when players are eliminated. When players are checkmated, they are eliminated and their pieces are disabled. The game continues as long as there are at least two players left.

## Adding a Chess Engine
- I created a custom chess engine. The chess engine used the same client interface as the client objects from the hub, and was included in the hub object with the other clients.
- The chess engine can be divided into two major pieces: the searcher, and the evaluator. The searcher searches through board states, and the evaluator evaluates board states. These two pieces work togethter to find the best move in the current board state.
- I added benchmarks tests and used pprof to profile the performance of the chess engine.
- A lot of refactoring was made to the chess logic to improve the performance of the the chess engine.
- I refactored how moves were found in a position. Originally, we were finding the legal moves by validating the moves at each search level. When searching, it's faster to assume every move is legal and to validate the current position at each level instead.
- I refactored how objects were created during the search. Originally, we were creating many objects while searching. For example, I was creating new piece objects for each move. When we're searching through millions of moves, creating new objects causes the garbage collector to bottleneck the program. So, I replaced enums/strings with integers, I referenced piece objects from a table instead of creating new ones, I avoided using interfaces where possible, I avoided using maps where possible, and I created a custom array type.
- I could have also improved performance with bitboard representations. This would have been difficult to implement for an arbitrary board size though. This could be something interesting to work on in the future.

## Chess Engine Search
- The search portion of the chess engine attempts to find the move that will maximize the player's score, in the future, according to the evaluator. The search uses concepts like: minimax, transposition tables, iterative deeping, hashing, and multithreading.
- We search through board states using minimax. Given that the chess engine will need to play games with more than two players, we're unable to use alpha-beta pruning. Alpha-beta pruning relies on the idea that what's good for one players is equally bad for the other player. This is not always the case when there're more than two players.
- We use transposition tables when we search to prevent searching the same position multiple times.
- The transposition tables contain hashes of board states. The hashes are created using the Zobrist hashing technique. To hash a board state using Zobrist hashing, we assign a random number to every possible chess piece and square combination in advance, and, whenever we need a hash, we XOR all the numbers coresponding to the current board state to create a hash. The resulting hash will virtually always be unique, and it can be computed quickly.
- I multithreaded the first layer of the search. As expected, the search became faster by one level of depth. In the future, I'd like to write an algorithm to dynamically allocate work across threads so all threads stay busy throughout the search.
- We use iterative deepening when searching. The engine doesn't take much of a performance hit from iterative deepening, and iterative deepening helps the engine find the best move it can within any time limit.

## Chess Engine Evaluation
- The evaluation portion of the chess engine takes a board state and returns a score for each player. In our evaluator, a player's score is determined by three factors: piece values, piece mobility, and piece positions. A higher score is better.
- A player's score increases with each piece they have on the board. Some pieces yield more score than others.
- A player's score increases with each move they can make. Generally, positions with many available moves are better than positions with few available moves.
- A player's score increases depending on where their pieces are placed. Some pieces give extra score depending on where they're placed on the board.
- Also, a player's score is affected by whether they've been checkmated.

