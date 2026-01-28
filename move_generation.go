package main

// NOTE: CASTLE_MOVE_SQS is for testing non-occupancy of squares for castling.
// CASTLE_CHECK_SQS is used for checking if king is passing through attacked sqs.
// Include E1/E8 and exclude B1/B8 when validating if squares are attacked
func initCastleSqs() {
	CASTLE_MOVE_SQS[0] = SQUARE_BBS[SQ_F1] | SQUARE_BBS[SQ_G1]
	CASTLE_MOVE_SQS[1] = SQUARE_BBS[SQ_D1] | SQUARE_BBS[SQ_C1] | SQUARE_BBS[SQ_B1]
	CASTLE_MOVE_SQS[2] = SQUARE_BBS[SQ_F8] | SQUARE_BBS[SQ_G8]
	CASTLE_MOVE_SQS[3] = SQUARE_BBS[SQ_D8] | SQUARE_BBS[SQ_C8] | SQUARE_BBS[SQ_B8]

	CASTLE_CHECK_SQS[0] = SQUARE_BBS[SQ_E1] | SQUARE_BBS[SQ_F1] | SQUARE_BBS[SQ_G1]
	CASTLE_CHECK_SQS[1] = SQUARE_BBS[SQ_E1] | SQUARE_BBS[SQ_D1] | SQUARE_BBS[SQ_C1]
	CASTLE_CHECK_SQS[2] = SQUARE_BBS[SQ_E8] | SQUARE_BBS[SQ_F8] | SQUARE_BBS[SQ_G8]
	CASTLE_CHECK_SQS[3] = SQUARE_BBS[SQ_E8] | SQUARE_BBS[SQ_D8] | SQUARE_BBS[SQ_C8]
}

func kingCastles(moveList *[]Move, occ Bitboard, color Color, cr int, enemyAttacks Bitboard) {
	for _, side := range SIDES {
		if canCastle(side, color, occ, cr, enemyAttacks) {
			*moveList = append(*moveList, CASTLE_MOVES[int(color)*2+side])
		}
	}
}

func pawnPushes(sq Square, dir int, occ Bitboard) Bitboard {
	pushes := signedShift(SQUARE_BBS[sq], dir) &^ occ
	if squareRank(sq) == secondRank(dir) {
		pushes |= signedShift(pushes, dir) &^ occ
	}
	return pushes
}

func serializeForPseudosKing(pl *[]Move, pieceBB Bitboard, occ Bitboard, selfOcc Bitboard, color Color, cr int, enemyAttacks Bitboard) {
	src := Square(lsb(pieceBB))
	serializeNormalMoves(pl, src, kingAttacks(occ, src)&^selfOcc, occ)
	kingCastles(pl, occ, color, cr, enemyAttacks)
}

func serializeForPseudosOther(pl *[]Move, pieceBB Bitboard, occ Bitboard, selfOcc Bitboard, fn AttackFunc) {
	for cursor := pieceBB; cursor != 0; cursor &= cursor - 1 {
		src := Square(lsb(cursor))
		serializeNormalMoves(pl, src, fn(occ, src)&^selfOcc, occ)
	}
}

// NOTE: if promoting, 2 moves are added (queen and knight promotions)
func serializeForPseudosPawns(pl *[]Move, pawns Bitboard, occ Bitboard, selfOcc Bitboard, color Color, epSq Square) {
	fDir, lRank := forward(color), lastRank(color)
	epSqBB := Bitboard(0)
	// To avoid indexing problems
	if epSq != NULL_SQ {
		epSqBB = SQUARE_BBS[epSq]
	}
	for cursor := pawns; cursor != 0; cursor &= cursor - 1 {
		src := Square(lsb(cursor))
		attacks := (PAWN_ATTACK_BBS[src][color] & (occ | epSqBB)) &^ selfOcc
		pseudos := attacks | pawnPushes(src, fDir, occ)
		for dstCursor := pseudos; dstCursor != 0; dstCursor &= dstCursor - 1 {
			dst := Square(lsb(dstCursor))
			switch {
			case squareRank(dst) == lRank:
				cOQ := capOrQuiet(occ, dst)
				*pl = append(*pl, toMove(dst, src, KNIGHT_PROMOTION|cOQ))
				*pl = append(*pl, toMove(dst, src, QUEEN_PROMOTION|cOQ))
				*pl = append(*pl, toMove(dst, src, ROOK_PROMOTION|cOQ))
				*pl = append(*pl, toMove(dst, src, BISHOP_PROMOTION|cOQ))
			case dst == epSq:
				*pl = append(*pl, toMove(dst, src, EP_CAPTURE))
			case dst == twoUp(src, color):
				*pl = append(*pl, toMove(dst, src, DOUBLE_PAWN_PUSH))
			default:
				*pl = append(*pl, toMove(dst, src, capOrQuiet(occ, dst)))
			}
		}
	}
}

func serializeNormalMoves(ml *[]Move, src Square, moves Bitboard, occ Bitboard) {
	for dstCursor := moves; dstCursor != 0; dstCursor &= dstCursor - 1 {
		dst := Square(lsb(dstCursor))
		*ml = append(*ml, toMove(dst, src, capOrQuiet(occ, dst)))
	}
}
