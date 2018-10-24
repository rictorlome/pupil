TODO:

- [ ] - Create parse fen and generate fen functions to place pieces
- [ ] - Create API to visually test results
- [ ] - Create test suite

REFERENCE:
Layout:

|       | A   | B   | C   | D   | E   | F   | G   | H   |
| ----- | --- | --- | --- | --- | --- | --- | --- | --- |
| **8** | 56  | 57  | 58  | 59  | 60  | 61  | 62  | 63  |
| **7** | 48  | 49  | 50  | 51  | 52  | 53  | 54  | 55  |
| **6** | 40  | 41  | 42  | 43  | 44  | 45  | 46  | 47  |
| **5** | 32  | 33  | 34  | 35  | 36  | 37  | 38  | 39  |
| **4** | 24  | 25  | 26  | 27  | 28  | 29  | 30  | 31  |
| **3** | 16  | 17  | 18  | 19  | 20  | 21  | 22  | 23  |
| **2** | 08  | 09  | 10  | 11  | 12  | 13  | 14  | 15  |
| **1** | 00  | 01  | 02  | 03  | 04  | 05  | 06  | 07  |

H8 - Most significant bit
A1 - Least significant bit

H8 H7 ... G8 G7 ... A2 A1

Bitboard
64 bit uint64. Each bit represent presence or absence of a piece.

Pieces
Represented by 12 bitboards. 1 for each piece, times 2 for each color.

1.  Pawn
2.  Rook
3.  Knight
4.  Bishop
5.  Queen
6.  King
