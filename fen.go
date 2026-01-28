package main

import (
	"strconv"
	"strings"
)

func bitboardsToGrid(bitboards []Bitboard) [8][8]string {
	var grid [8][8]string
	for _, s := range SQUARES {
		for _, piece := range PIECES {
			if occupiedAtSq(bitboards[piece], s) {
				grid[7-squareRank(s)][squareFile(s)] = piece.String()
			}
		}
	}
	return grid
}

func generateCastleString(s StateInfo) string {
	var castleString string
	// Reversed because bit indices are right to left
	for idx, char := range "qkQK" {
		if hasRight(s.castlingRights, uint(idx)) {
			castleString = string(char) + castleString
		}
	}
	if castleString == "" {
		return "-"
	}
	return castleString
}

func generateColorString(pos Position) string {
	if pos.sideToMove() == WHITE {
		return "w"
	}
	return "b"
}

func generateRule50String(s StateInfo) string {
	return strconv.Itoa(s.rule50)
}

func generateFen(pos Position) string {
	var fenArr []string
	fenArr = append(fenArr, gridToFen(bitboardsToGrid(pos.placement)))
	fenArr = append(fenArr, generateColorString(pos))
	fenArr = append(fenArr, generateCastleString(*pos.state))
	fenArr = append(fenArr, pos.state.epSq.String())
	fenArr = append(fenArr, generateRule50String(*pos.state))
	fenArr = append(fenArr, strconv.Itoa(pos.moveCount()))
	return strings.Join(fenArr, " ")
}

func gridToFen(grid [8][8]string) string {
	var fenArr []string
	for _, row := range grid {
		fenString := ""
		offset := 0
		for _, sq := range row {
			if sq == "" {
				offset += 1
			} else {
				if offset == 0 {
					fenString += sq
				} else {
					fenString += strconv.Itoa(offset) + sq
					offset = 0
				}
			}
		}
		if offset != 0 {
			fenString += strconv.Itoa(offset)
		}
		fenArr = append(fenArr, fenString)
	}
	return strings.Join(fenArr, "/")
}

func parseColor(s string) Color {
	if s == "w" {
		return WHITE
	}
	return BLACK
}

func parseFen(fen string) Position {
	fields := strings.Split(fen, " ")
	moveCount, _ := strconv.Atoi(fields[5])
	// Core fen info
	p := Position{}
	p.state = &StateInfo{}
	p.setFenInfo(fields[0], fields[1], fields[2], fields[3], fields[4], moveCount)
	// Additional state info
	p.state.key = p.toZobrist()
	p.state.oppositeColorAttacks = p.getColorAttacks(opposite(p.sideToMove()))
	p.state.blockersForKing = p.sliderBlockers(opposite(p.sideToMove()), p.kingSquare(p.sideToMove()))
	p.state.prev = nil

	return p
}

func parseMoveType(promstring string, occ Bitboard, src Square, dst Square, moverType PieceType, captured Piece) MoveType {
	fileDiff := squareFile(dst) - squareFile(src)
	rankDiff := squareRank(dst) - squareRank(src)
	// Promotion
	var mt MoveType
	if promstring != "" {
		idx := strings.Index(strings.Join(PROMOTION_STRINGS, ""), promstring)
		mt = MoveType(idx<<12) | MoveType(PROMOTION_MASK)
	}
	// Capture
	mt |= capOrQuiet(occ, dst)
	// En passant
	if moverType == PAWN && squareFile(dst) != squareFile(src) && captured == NULL_PIECE {
		mt |= EP_CAPTURE
	}
	// Double Pawn Push
	if moverType == PAWN && (rankDiff == 2 || rankDiff == -2) {
		mt |= DOUBLE_PAWN_PUSH
	}
	// Castles
	if moverType == KING && fileDiff == 2 {
		mt |= KING_CASTLE
	}
	if moverType == KING && fileDiff == -2 {
		mt |= QUEEN_CASTLE
	}
	return mt
}

func parsePositions(positions string) []Bitboard {
	resultBBs := make([]Bitboard, 12)
	ranks := strings.Split(positions, "/")
	for rank, rankString := range ranks {
		offset := 0
		for file, sq := range rankString {
			if strings.Contains(PIECE_STRING, string(sq)) {
				index := makePiece(string(sq))
				sqNum := makeSquare(RANK_8-rank, file+offset)
				resultBBs[index] |= SQUARE_BBS[sqNum]
			} else {
				offset += int(sq-'0') - 1
			}
		}
	}
	return resultBBs
}

func parseSquare(sq string) Square {
	if sq == "-" {
		return NULL_SQ
	}
	rank := int(sq[1]-'0') - 1
	return makeSquare(rank, strings.Index(FILES, sq[0:1]))
}
