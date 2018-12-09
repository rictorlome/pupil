package main

/*
Table of Contents:

1. Colors (int)
2. Directions (int)
3. File Bitboards (uint64)
4. Rank Bitboards (uint64)
5. Other Bitboards (uint64)
6. Squares (uint)
7. Ranks (int)
8. Files (int)
9. Pieces (int)
10. PieceType (int)
11. FEN (string)
12. CastlingRights (int)
13. Move
14. AttackFuncs

*/
type AttackFunc func(Bitboard, Square) Bitboard
type Bitboard uint64
type Color int
type Move uint16
type MoveType uint16
type Piece uint
type PieceType uint
type Position struct {
	move_count int
	placement  []Bitboard
	state      StateInfo
	to_move    Color
}
type Square uint
type StateInfo struct {
	castling_rights int
	ep_sq           Square
	prev            *StateInfo
	rule_50         int

	blockers_for_king Bitboard
}

// 1. COLORS
var WHITE Color = 0
var BLACK Color = 1
var COLORS = []Color{WHITE, BLACK}

// 2. SIDES
var KING_SIDE int = 0
var QUEEN_SIDE int = 1
var SIDES = []int{KING_SIDE, QUEEN_SIDE}

// 2. DIRECTIONS
var NORTH int = 8
var EAST int = 1
var SOUTH int = -NORTH
var WEST int = -EAST
var NORTH_WEST int = NORTH + WEST
var NORTH_EAST int = NORTH + EAST
var SOUTH_WEST int = SOUTH + WEST
var SOUTH_EAST int = SOUTH + EAST
var DIRECTIONS = []int{
	NORTH, NORTH_EAST, EAST, SOUTH_EAST,
	SOUTH, SOUTH_WEST, WEST, NORTH_WEST,
}

// By piece
var BISHOP_DIRECTIONS = []int{
	NORTH_EAST, SOUTH_EAST, SOUTH_WEST, NORTH_WEST,
}

// Each line (tuple) corresponds to a direction
var KNIGHT_DIRECTIONS = []int{
	NORTH, NORTH_EAST,
	NORTH, NORTH_WEST,
	EAST, NORTH_EAST,
	EAST, SOUTH_EAST,
	SOUTH, SOUTH_EAST,
	SOUTH, SOUTH_WEST,
	WEST, NORTH_WEST,
	WEST, SOUTH_WEST,
}
var ROOK_DIRECTIONS = []int{
	NORTH, EAST, SOUTH, WEST,
}

// 3. FILE BITBOARDS
var FILE_ABB Bitboard = 0x101010101010101
var FILE_BBB Bitboard = signed_shift(FILE_ABB, EAST*1)
var FILE_CBB Bitboard = signed_shift(FILE_ABB, EAST*2)
var FILE_DBB Bitboard = signed_shift(FILE_ABB, EAST*3)
var FILE_EBB Bitboard = signed_shift(FILE_ABB, EAST*4)
var FILE_FBB Bitboard = signed_shift(FILE_ABB, EAST*5)
var FILE_GBB Bitboard = signed_shift(FILE_ABB, EAST*6)
var FILE_HBB Bitboard = signed_shift(FILE_ABB, EAST*7)
var FILE_BBS = []Bitboard{
	FILE_ABB, FILE_BBB, FILE_CBB, FILE_DBB,
	FILE_EBB, FILE_FBB, FILE_GBB, FILE_HBB,
}
var FILES string = "abcdefgh"

// 4. RANK BITBOARDS
var RANK_1BB Bitboard = 0xff
var RANK_2BB Bitboard = signed_shift(RANK_1BB, NORTH*1)
var RANK_3BB Bitboard = signed_shift(RANK_1BB, NORTH*2)
var RANK_4BB Bitboard = signed_shift(RANK_1BB, NORTH*3)
var RANK_5BB Bitboard = signed_shift(RANK_1BB, NORTH*4)
var RANK_6BB Bitboard = signed_shift(RANK_1BB, NORTH*5)
var RANK_7BB Bitboard = signed_shift(RANK_1BB, NORTH*6)
var RANK_8BB Bitboard = signed_shift(RANK_1BB, NORTH*7)
var RANK_BBS = []Bitboard{
	RANK_1BB, RANK_2BB, RANK_3BB, RANK_4BB,
	RANK_5BB, RANK_6BB, RANK_7BB, RANK_8BB,
}

// 5. OTHER BITBOARDS
var ALL_SQS Bitboard = 0xffffffffffffffff

// Initialized in main.init
var SQUARE_BBS [64]Bitboard
var KNIGHT_ATTACK_BBS [64]Bitboard
var KING_ATTACK_BBS [64]Bitboard
var PAWN_ATTACK_BBS [64][2]Bitboard

var ROOK_ATTACK_MASKS [64]Bitboard
var BISHOP_ATTACK_MASKS [64]Bitboard

var ROOK_OCCUPANCY_MASKS [64]Bitboard
var BISHOP_OCCUPANCY_MASKS [64]Bitboard

var RELEVANT_ROOK_OCCUPANCY [64]Bitboard
var RELEVANT_BISHOP_OCCUPANCY [64]Bitboard

var LINE_BBS [64][64]Bitboard

var RELEVANT_ROOK_OCCUPANCY_BITS = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}
var RELEVANT_BISHOP_OCCUPANCY_BITS = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

var CASTLE_MOVE_SQS [4]Bitboard
var CASTLE_CHECK_SQS [4]Bitboard

// 6. SQUARES
const (
	SQ_A1 Square = iota
	SQ_B1
	SQ_C1
	SQ_D1
	SQ_E1
	SQ_F1
	SQ_G1
	SQ_H1
	SQ_A2
	SQ_B2
	SQ_C2
	SQ_D2
	SQ_E2
	SQ_F2
	SQ_G2
	SQ_H2
	SQ_A3
	SQ_B3
	SQ_C3
	SQ_D3
	SQ_E3
	SQ_F3
	SQ_G3
	SQ_H3
	SQ_A4
	SQ_B4
	SQ_C4
	SQ_D4
	SQ_E4
	SQ_F4
	SQ_G4
	SQ_H4
	SQ_A5
	SQ_B5
	SQ_C5
	SQ_D5
	SQ_E5
	SQ_F5
	SQ_G5
	SQ_H5
	SQ_A6
	SQ_B6
	SQ_C6
	SQ_D6
	SQ_E6
	SQ_F6
	SQ_G6
	SQ_H6
	SQ_A7
	SQ_B7
	SQ_C7
	SQ_D7
	SQ_E7
	SQ_F7
	SQ_G7
	SQ_H7
	SQ_A8
	SQ_B8
	SQ_C8
	SQ_D8
	SQ_E8
	SQ_F8
	SQ_G8
	SQ_H8
	NULL_SQ // NULL SQUARE = 64
)

// For looping, excludes null square
var SQUARES = []Square{
	SQ_A1, SQ_B1, SQ_C1, SQ_D1, SQ_E1, SQ_F1, SQ_G1, SQ_H1,
	SQ_A2, SQ_B2, SQ_C2, SQ_D2, SQ_E2, SQ_F2, SQ_G2, SQ_H2,
	SQ_A3, SQ_B3, SQ_C3, SQ_D3, SQ_E3, SQ_F3, SQ_G3, SQ_H3,
	SQ_A4, SQ_B4, SQ_C4, SQ_D4, SQ_E4, SQ_F4, SQ_G4, SQ_H4,
	SQ_A5, SQ_B5, SQ_C5, SQ_D5, SQ_E5, SQ_F5, SQ_G5, SQ_H5,
	SQ_A6, SQ_B6, SQ_C6, SQ_D6, SQ_E6, SQ_F6, SQ_G6, SQ_H6,
	SQ_A7, SQ_B7, SQ_C7, SQ_D7, SQ_E7, SQ_F7, SQ_G7, SQ_H7,
	SQ_A8, SQ_B8, SQ_C8, SQ_D8, SQ_E8, SQ_F8, SQ_G8, SQ_H8,
}

// 7. RANKS
var RANK_1 int = 0
var RANK_2 int = 1
var RANK_3 int = 2
var RANK_4 int = 3
var RANK_5 int = 4
var RANK_6 int = 5
var RANK_7 int = 6
var RANK_8 int = 7

// 8. FILES
var FILE_A int = 0
var FILE_B int = 1
var FILE_C int = 2
var FILE_D int = 3
var FILE_E int = 4
var FILE_F int = 5
var FILE_G int = 6
var FILE_H int = 7

// 9. PIECES
const (
	WK Piece = iota
	WQ
	WB
	WN
	WR
	WP
	BK
	BQ
	BB
	BN
	BR
	BP
	NULL_PIECE
)

var WHITE_PIECES = []Piece{
	WK, WQ, WB, WN, WR, WP,
}
var BLACK_PIECES = []Piece{
	BK, BQ, BB, BN, BR, BP,
}
var PIECES = []Piece{
	WK, WQ, WB, WN, WR, WP,
	BK, BQ, BB, BN, BR, BP,
}
var PIECE_STRING string = "KQBNRPkqbnrp"

// 10. PieceType
const (
	KING PieceType = iota
	QUEEN
	BISHOP
	KNIGHT
	ROOK
	PAWN
	NULL_PIECE_TYPE
)

var PIECE_TYPES = []PieceType{
	KING, QUEEN, BISHOP, KNIGHT, ROOK, PAWN, NULL_PIECE_TYPE,
}

// 11. FEN
var INITIAL_FEN string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
var INITIAL_FEN_JUST_PIECES string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"

// 12. CASTLING RIGHTS - KQkq
var BQ_CASTLE int = 0x1
var BK_CASTLE int = BQ_CASTLE << 1
var WQ_CASTLE int = BQ_CASTLE << 2
var WK_CASTLE int = BQ_CASTLE << 3
var CASTLING_RIGHTS = [4]int{
	WK_CASTLE, WQ_CASTLE, BK_CASTLE, BQ_CASTLE,
}

var BLACK_CASTLES int = BQ_CASTLE | BK_CASTLE
var WHITE_CASLTES int = WQ_CASTLE | BK_CASTLE
var NO_CASTLE int = 0

var CHAR_TO_CASTLE = map[string]int{
	"q": BQ_CASTLE, "k": BK_CASTLE,
	"Q": WQ_CASTLE, "K": WK_CASTLE,
	"-": NO_CASTLE,
}

// 13. MOVE (uin16)
// bit 0 - 5 (dest sq)
// bit 6 - 11 (source sq)
// bit 12 - 15 (move type)
var DST_MASK Move = 0xFFC0
var SRC_MASK Move = 0xF03F
var MOVE_TYPE_MASK Move = 0xFFF

const (
	QUIET MoveType = iota << 12
	DOUBLE_PAWN_PUSH
	KING_CASTLE
	QUEEN_CASTLE
	CAPTURE
	EP_CAPTURE
)
const (
	KNIGHT_PROMOTION MoveType = (iota + 8) << 12
	BISHOP_PROMOTION
	ROOK_PROMOTION
	QUEEN_PROMOTION
	KNIGHT_PROMOTION_CAPTURE
	BISHOP_PROMOTION_CAPTURE
	ROOK_PROMOTION_CAPTURE
	QUEEN_PROMOTION_CAPTURE
)

var MOVE_TYPES = []MoveType{
	QUIET, DOUBLE_PAWN_PUSH, KING_CASTLE, QUEEN_CASTLE,
	CAPTURE, EP_CAPTURE /* 6-7 unused */, QUIET, QUIET,
	KNIGHT_PROMOTION, BISHOP_PROMOTION, ROOK_PROMOTION, QUEEN_PROMOTION,
	KNIGHT_PROMOTION_CAPTURE, BISHOP_PROMOTION_CAPTURE, ROOK_PROMOTION_CAPTURE, QUEEN_PROMOTION_CAPTURE,
}

var PROMOTION_STRINGS = []string{
	"", "", "", "", "", "", "", "",
	"n", "b", "r", "q", "n", "b", "r", "q",
}

var BLACK_KINGSIDE Move = to_move(SQ_E8, SQ_G8, KING_CASTLE)
var BLACK_QUEENSIDE Move = to_move(SQ_E8, SQ_C8, QUEEN_CASTLE)
var WHITE_KINGSIDE Move = to_move(SQ_E1, SQ_G1, KING_CASTLE)
var WHITE_QUEENSIDE Move = to_move(SQ_E1, SQ_C1, QUEEN_CASTLE)

var CASTLE_MOVES = []Move{
	WHITE_KINGSIDE, WHITE_QUEENSIDE, BLACK_KINGSIDE, BLACK_QUEENSIDE,
}

// 14. AttackFuncs
// For indexing by PieceType.
// PAWN returns null_attacks because pawn attacks are computed separately, in bulk.
var AttackFuncs = []AttackFunc{
	king_attacks, queen_attacks, bishop_attacks, knight_attacks, rook_attacks, null_attacks,
}
