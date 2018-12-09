# Pupil

A WIP `Go`-based chess engine destined to beat me in a 10 minute game.

---

## Description

| Feature                | Implementation                                                                                                                                                                |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Board representation   | [Bitboard](https://www.chessprogramming.org/Bitboards)                                                                                                                        |
| Square mapping         | [Little-endian Rank-file Mapping](https://www.chessprogramming.org/Square_Mapping_Considerations#Little-Endian_Rank-File_Mapping)                                             |
| Bitboard Serialization | [Forward-scanning](https://www.chessprogramming.org/Bitboard_Serialization#Scanning_Forward)                                                                                  |
| Move encoding          | [16 bit From-to based](https://www.chessprogramming.org/Encoding_Moves#From-To_Based)                                                                                         |
| Move Generation        | [Classical approach for sliders](https://www.chessprogramming.org/Classical_Approach), working on [Magic Bitboard approach](https://www.chessprogramming.org/Magic_Bitboards) |
| Search                 | TBD                                                                                                                                                                           |
| Evaluation             | TBD                                                                                                                                                                           |

---

## History

This is my third attempt at a chess engine.

### The predecessors:

1.  [Rhess](https://github.com/rictorlome/rhess) - a command line game written in `Ruby` using a [mailbox](https://www.chessprogramming.org/Mailbox) 8x8 board representation. Styled with 100% unicode and playable via keyboard.

2.  [Gogochess](https://github.com/rictorlome/gogochess) - a chess-move generator written in `Go` with a custom `HashMap` based move-list. Passing [perft tests](https://www.chessprogramming.org/Perft_Results) up to depth 5 and capable of solving `Mate-in-Twos` in less than a minute. Configured with `HTTP` API for browser integration and easy testing.

---

## Todo

- [x] Revisit move encoding
- [ ] Precompute relevant bitboards for pawn pushes
- [ ] Brainstorm and implement `Position.do_move`
- [ ] Pseudo-legal move gen -> Legal move gen
