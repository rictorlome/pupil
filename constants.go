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
type Key uint64
type Move uint16
type MoveType uint16
type Piece uint
type PieceType uint
type Position struct {
	occ                 Bitboard
	placement           []Bitboard
	placement_by_square []Piece
	ply                 int
	state               *StateInfo
	stm                 Color
}
type Square uint
type StateInfo struct {
	// Core fen info
	castling_rights int
	ep_sq           Square
	rule_50         int
	// Additional info
	key                    Key
	opposite_color_attacks Bitboard
	blockers_for_king      Bitboard
	prev                   *StateInfo
	captured               Piece
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
var NEIGHBOR_BBS [64]Bitboard
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
var BETWEEN_BBS [64][64]Bitboard

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

var ROOK_SRC_DST [64][2]Square

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
	KING, QUEEN, BISHOP, KNIGHT, ROOK, PAWN,
	KING, QUEEN, BISHOP, KNIGHT, ROOK, PAWN,
	NULL_PIECE_TYPE,
}

var PIECE_TYPE_STRINGS = []string{
	"KING", "QUEEN", "BISHOP", "KNIGHT", "ROOK", "PAWN",
}

var PROMOTION_PIECE_TYPES = []PieceType{
	KNIGHT, BISHOP, ROOK, QUEEN,
}

// 11. FEN
var INITIAL_FEN string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
var INITIAL_FEN_JUST_PIECES string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
var KIWIPETE_FEN string = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

// 12. CASTLING RIGHTS - KQkq
var BQ_CASTLE int = 0x1
var BK_CASTLE int = BQ_CASTLE << 1
var WQ_CASTLE int = BQ_CASTLE << 2
var WK_CASTLE int = BQ_CASTLE << 3
var CASTLING_RIGHTS = [4]int{
	WK_CASTLE, WQ_CASTLE, BK_CASTLE, BQ_CASTLE,
}

var BLACK_CASTLES int = BQ_CASTLE | BK_CASTLE
var WHITE_CASTLES int = WQ_CASTLE | WK_CASTLE
var ALL_CASTLES int = BLACK_CASTLES | WHITE_CASTLES
var NO_CASTLE int = 0

var CHAR_TO_CASTLE = map[string]int{
	"q": BQ_CASTLE, "k": BK_CASTLE,
	"Q": WQ_CASTLE, "K": WK_CASTLE,
	"-": NO_CASTLE,
}

var CASTLING_MASK_BY_SQ [64]int

// 13. MOVE (uint16)
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

var PROMOTION_MASK Move = Move(KNIGHT_PROMOTION)

var MOVE_TYPES = []MoveType{
	QUIET, DOUBLE_PAWN_PUSH, KING_CASTLE, QUEEN_CASTLE,
	CAPTURE, EP_CAPTURE /* 6-7 unused */, QUIET, QUIET,
	KNIGHT_PROMOTION, BISHOP_PROMOTION, ROOK_PROMOTION, QUEEN_PROMOTION,
	KNIGHT_PROMOTION_CAPTURE, BISHOP_PROMOTION_CAPTURE, ROOK_PROMOTION_CAPTURE, QUEEN_PROMOTION_CAPTURE,
}

var MOVE_TYPE_STRINGS = []string{
	"QUIET", "DOUBLE_PAWN_PUSH", "KING_CASTLE", "QUEEN_CASTLE",
	"CAPTURE", "EP_CAPTURE" /* 6-7 unused */, "QUIET", "QUIET",
	"KNIGHT_PROMOTION", "BISHOP_PROMOTION", "ROOK_PROMOTION", "QUEEN_PROMOTION",
	"KNIGHT_PROMOTION_CAPTURE", "BISHOP_PROMOTION_CAPTURE", "ROOK_PROMOTION_CAPTURE", "QUEEN_PROMOTION_CAPTURE",
}

var PROMOTION_STRINGS = []string{
	"", "", "", "", "", "", "", "",
	"n", "b", "r", "q", "n", "b", "r", "q",
}

var BLACK_KINGSIDE Move = to_move(SQ_G8, SQ_E8, KING_CASTLE)
var BLACK_QUEENSIDE Move = to_move(SQ_C8, SQ_E8, QUEEN_CASTLE)
var WHITE_KINGSIDE Move = to_move(SQ_G1, SQ_E1, KING_CASTLE)
var WHITE_QUEENSIDE Move = to_move(SQ_C1, SQ_E1, QUEEN_CASTLE)

var CASTLE_MOVES = []Move{
	WHITE_KINGSIDE, WHITE_QUEENSIDE, BLACK_KINGSIDE, BLACK_QUEENSIDE,
}

// 14. AttackFuncs
// For indexing by PieceType.
// PAWN returns null_attacks because pawn attacks are computed separately, in bulk.
var AttackFuncs = []AttackFunc{
	king_attacks, queen_attacks, bishop_attacks, knight_attacks, rook_attacks, null_attacks,
}

// 15. Magics
var RookAttackTable = make([]Bitboard, 0x19000)
var BishopAttackTable = make([]Bitboard, 0x1480)

var RookMagics [64]Magic
var BishopMagics [64]Magic

var RookMagicNums = [64]uint64{
	108104671511543808,
	18014535952633857,
	2954379670400794880,
	2954378947825565856,
	4647718131169495040,
	1297041099320199169,
	2449993940007651409,
	2449958755704832076,
	73324262571671584,
	297307944687521792,
	595741856929678464,
	4612953205712621568,
	4614219310666612992,
	844442110134528,
	5189272705265959172,
	32791861655505152,
	36029071901073408,
	2058145305659973632,
	4611759686243385619,
	9422938895039856928,
	596868237938001920,
	4981685425136468480,
	1103806726148,
	4612255565460160796,
	9331463119089737728,
	108262316138582018,
	3114096104833280,
	9223381936756494337,
	149753485850837120,
	564607811063816,
	2216203192336,
	72343475651102865,
	2326144666862354564,
	9817882440763457536,
	15903404912025600,
	9228438621365735488,
	36037595267859456,
	2201179128832,
	4621293557540456528,
	35463578517636,
	532438514434052,
	2621240552865792,
	142941348036624,
	595038238377377808,
	8796160163968,
	3463338499909353744,
	71751992344584,
	4561341172285441,
	36028934462115904,
	315287160573919360,
	9367911773858955904,
	432363156547863168,
	54183967376801920,
	10137499355709568,
	17626573177856,
	288231492847444480,
	9441820810189807874,
	22536148670808134,
	6919816012901188113,
	4535620731137,
	180988444385739269,
	563534337475586,
	328780511021891844,
	4616330940799844898,
}

var BishopMagicNums = [64]uint64{
	9008315954660864,
	11531467962850607104,
	1154122173485645832,
	649648646845384704,
	6920912362364077377,
	4612250085613379618,
	4789477016863872,
	144133059503988864,
	1297072985223987328,
	7789081117325395586,
	37391989604576,
	72062284217780288,
	288267768423285008,
	144117391663104012,
	45038198015410704,
	144678690335297536,
	882705664470782470,
	1252073281490830592,
	4724278242670813194,
	1125935407433728,
	1131056025585664,
	5046292197922643970,
	4901047794737350721,
	2306423590159191296,
	37159096049247232,
	40656160491307136,
	1244190053636047364,
	2260600203780128,
	576742302458855424,
	4701901084959051776,
	72629961009172486,
	6055657188738990595,
	22535607658489888,
	9228025208718495748,
	38487207183380,
	35218732156944,
	9409152393700180224,
	10380814742402122756,
	1531789031409713294,
	2605851684668417,
	31815503895921665,
	2306023346841584129,
	577130974567993348,
	2347079781972968704,
	2306973616674179584,
	721148579787326208,
	1751144284030596,
	440921359384704,
	576533328933421568,
	1196560977821952,
	9232520540542468096,
	288239242072689697,
	4614219362508865537,
	5260804732511330304,
	292879180041945744,
	585470167778075712,
	576533321161916416,
	18148586004480,
	846075306168320,
	576465150754719745,
	1981583836111193093,
	2666167540381141124,
	1153029840879027200,
	9241424497968423040,
}

// https://chess.stackexchange.com/questions/4490/maximum-possible-movement-in-a-turn
const MAX_BRANCHING int = 218

// Values taken from Michniewski's simple evaluation function
// https://www.chessprogramming.org/Simplified_Evaluation_Function
const K_VAL int = 20000
const Q_VAL int = 900
const B_VAL int = 330
const N_VAL int = 320
const R_VAL int = 500
const P_VAL int = 100

var MATERIAL_VALUES = []int{
	K_VAL, Q_VAL, B_VAL, N_VAL, R_VAL, P_VAL,
}

var KING_SQUARE_VALUES_MIDGAME = []int{
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	20, 20, 0, 0, 0, 0, 20, 20,
	20, 30, 10, 0, 0, 10, 30, 20,
}

var KING_SQUARE_VALUES_ENDGAME = []int{
	-50, -40, -30, -20, -20, -30, -40, -50,
	-30, -20, -10, 0, 0, -10, -20, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -30, 0, 0, 0, 0, -30, -30,
	-50, -30, -30, -30, -30, -30, -30, -50,
}

var QUEEN_SQUARE_VALUES = []int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, -5,
	-10, 5, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
}

var BISHOP_SQUARE_VALUES = []int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}

var KNIGHT_SQUARE_VALUES = []int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

var ROOK_SQUARE_VALUES = []int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var PAWN_SQUARE_VALUES = []int{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var POSITION_VALUES = []*[]int{
	&KING_SQUARE_VALUES_MIDGAME, &QUEEN_SQUARE_VALUES,
	&BISHOP_SQUARE_VALUES, &KNIGHT_SQUARE_VALUES,
	&ROOK_SQUARE_VALUES, &PAWN_SQUARE_VALUES,
}
