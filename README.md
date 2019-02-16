# Pupil

Pupil is a chess engine written in `Go`. The goal of writing this program was to create an engine capable of beating me in a 10 minute game.

---

## Description

| Feature                | Implementation                                                                                                                                                                                                                                                                   |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Board representation   | [Bitboard](https://www.chessprogramming.org/Bitboards)                                                                                                                                                                                                                           |
| Square mapping         | [Little-endian Rank-file Mapping](https://www.chessprogramming.org/Square_Mapping_Considerations#Little-Endian_Rank-File_Mapping)                                                                                                                                                |
| Bitboard Serialization | [Forward-scanning](https://www.chessprogramming.org/Bitboard_Serialization#Scanning_Forward)                                                                                                                                                                                     |
| Move encoding          | [16 bit From-to based](https://www.chessprogramming.org/Encoding_Moves#From-To_Based)                                                                                                                                                                                            |
| Move Generation        | [Magic Bitboard approach](https://www.chessprogramming.org/Magic_Bitboards)                                                                                                                                                                                                      |
| Search                 | [AlphaBeta search](https://www.chessprogramming.org/Alpha-Beta) within the [Negamax](https://www.chessprogramming.org/Negamax) framework. Currently uses [transposition tables](https://www.chessprogramming.org/Transposition_Table) to cache static evaluations of leaf nodes. |
| Evaluation             | [Simplified Evaluation Function](https://www.chessprogramming.org/Simplified_Evaluation_Function): Material value and positional value based on precomputed arrays.                                                                                                              |

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
- [ ] Quiescence search
- [ ] Iterative deepening with time limit
- [ ] Flesh out frontend for better UX
- [ ] Experiment compiling to WebAssembly
- [ ] Deploy (WebAssembly or no)

## Nice to haves:

- [ ] Concurrent alpha-beta
- [ ] Lint
