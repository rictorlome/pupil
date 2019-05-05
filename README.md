# Pupil

Pupil is a chess engine written in `Go`. The goal of writing this program was to create an engine capable of beating me in a 10 minute game.

[Perft tests](https://www.chessprogramming.org/Perft) pass for the [initial board state](https://www.chessprogramming.org/Perft_Results#Initial_Position) and for [kiwipete](https://www.chessprogramming.org/Perft_Results#Position_2) up to depths 6 and 5 respectively.

Some of the major optimizations include: [bitboard](https://www.chessprogramming.org/Bitboards) piece representation, [magic bitboards](https://www.chessprogramming.org/Magic_Bitboards) for move generation, multi-threaded perft tests (for development speed), object pooling, and a transposition table keyed via [zobrist hashes](https://www.chessprogramming.org/Zobrist_Hashing).

Although nowhere near as complete, the code is heavily inspired by the open-source [Stockfish](https://github.com/official-stockfish/Stockfish) project, whose source code I dipped into heavily.

#### Note to anyone reading the code:

Apologies for the weird code style. I began this project before installing a linter in my editor/becoming familiar with the golang style guide. I haven't been able to find a good codemod/script to automatically clean it up. If you know of any, please let me know.

---

## Description

| Feature                | Implementation                                                                                                                                                                                                                                                                                                                            |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Board representation   | [Bitboard](https://www.chessprogramming.org/Bitboards)                                                                                                                                                                                                                                                                                    |
| Square mapping         | [Little-endian Rank-file Mapping](https://www.chessprogramming.org/Square_Mapping_Considerations#Little-Endian_Rank-File_Mapping)                                                                                                                                                                                                         |
| Bitboard Serialization | [Forward-scanning](https://www.chessprogramming.org/Bitboard_Serialization#Scanning_Forward)                                                                                                                                                                                                                                              |
| Move encoding          | [16 bit From-to based](https://www.chessprogramming.org/Encoding_Moves#From-To_Based)                                                                                                                                                                                                                                                     |
| Move Generation        | [Magic Bitboard approach](https://www.chessprogramming.org/Magic_Bitboards)                                                                                                                                                                                                                                                               |
| Search                 | [AlphaBeta search](https://www.chessprogramming.org/Alpha-Beta) within the [Negamax](https://www.chessprogramming.org/Negamax) framework. Currently uses [transposition tables](https://www.chessprogramming.org/Transposition_Table) to cache information about the different [node types](https://www.chessprogramming.org/Node_Types). |
| Evaluation             | [Simplified Evaluation Function](https://www.chessprogramming.org/Simplified_Evaluation_Function): Material value and positional value based on precomputed arrays.                                                                                                                                                                       |

---

## History

This is my third attempt at a chess engine.

### The predecessors:

1.  [Rhess](https://github.com/rictorlome/rhess) - a command line game written in `Ruby` using a [mailbox](https://www.chessprogramming.org/Mailbox) 8x8 board representation. Styled with 100% unicode and playable via keyboard.

2.  [Gogochess](https://github.com/rictorlome/gogochess) - a chess-move generator written in `Go` with a custom `HashMap` based move-list. Passing [perft tests](https://www.chessprogramming.org/Perft_Results) up to depth 5 and capable of solving `Mate-in-Twos` in less than a minute. Configured with `HTTP` API for browser integration and easy testing.

---

## Todo

- [x] Test alpha-beta against negamax.
- [x] Update transposition table for non-perft search
- [x] Cache best move
- [x] Update alpha-beta to search best move first
- [x] Limit cache entry size
- [x] Implement LRU/better cache clearing (simple array)
- [x] Flesh out frontend for better UX
- [x] Experiment compiling to WebAssembly
- [x] Deploy (WebAssembly or no)

## Nice to haves:

- [ ] Quiescence search
- [ ] Iterative deepening with time limit
- [ ] Concurrent alpha-beta
- [ ] Lint
