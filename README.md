# Pupil

[Live link](https://rictorlome.github.io/pupil/)
## Overview

Pupil is a chess engine written in `Go`. The goal of writing this program was to create an engine capable of beating me in a 10 minute game.

[Perft tests](https://www.chessprogramming.org/Perft) pass for the [initial board state](https://www.chessprogramming.org/Perft_Results#Initial_Position) and for [kiwipete](https://www.chessprogramming.org/Perft_Results#Position_2) up to depths 8 and 6 respectively.

Some of the major optimizations include: [bitboard](https://www.chessprogramming.org/Bitboards) piece representation, [magic bitboards](https://www.chessprogramming.org/Magic_Bitboards) for move generation, [quiescence search](https://www.chessprogramming.org/Quiescence_Search), [killer move heuristic](https://www.chessprogramming.org/Killer_Heuristic), [MVV-LVA](https://www.chessprogramming.org/MVV-LVA) move ordering, object pooling, and a [transposition table](https://www.chessprogramming.org/Transposition_Table) keyed via [zobrist hashes](https://www.chessprogramming.org/Zobrist_Hashing).

Although nowhere near as complete, the code is heavily inspired by the open-source [Stockfish](https://github.com/official-stockfish/Stockfish) project, whose source code I dipped into heavily.

## Running the code

```bash
GOOS=js GOARCH=wasm go build -o docs/pupil.wasm     # fresh build of wasm
go run .                                            # start local server
```

Then visit `http://localhost:8080`.

---

## Description

| Feature                | Implementation                                                                                                                                                                                                                                                                                                                            |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Board representation   | [Bitboard](https://www.chessprogramming.org/Bitboards)                                                                                                                                                                                                                                                                                    |
| Square mapping         | [Little-endian Rank-file Mapping](https://www.chessprogramming.org/Square_Mapping_Considerations#Little-Endian_Rank-File_Mapping)                                                                                                                                                                                                         |
| Bitboard Serialization | [Forward-scanning](https://www.chessprogramming.org/Bitboard_Serialization#Scanning_Forward)                                                                                                                                                                                                                                              |
| Move encoding          | [16 bit From-to based](https://www.chessprogramming.org/Encoding_Moves#From-To_Based)                                                                                                                                                                                                                                                     |
| Move Generation        | [Magic Bitboard approach](https://www.chessprogramming.org/Magic_Bitboards)                                                                                                                                                                                                                                                               |
| Search                 | [Alpha-beta](https://www.chessprogramming.org/Alpha-Beta) with [quiescence search](https://www.chessprogramming.org/Quiescence_Search), [transposition tables](https://www.chessprogramming.org/Transposition_Table), [killer moves](https://www.chessprogramming.org/Killer_Heuristic), and [MVV-LVA](https://www.chessprogramming.org/MVV-LVA) move ordering. |
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
- [x] Deploy with WebAssembly
- [x] Quiescence search
- [x] Killer moves
- [x] MVV-LVA move ordering
- [x] Move list pooling
- [x] Lint

## Nice to haves:

- [ ] Iterative deepening with time limit
- [ ] Concurrent alpha-beta
