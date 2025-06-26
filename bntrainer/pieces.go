package bntrainer

type PieceType int

const (
	King PieceType = iota
	Knight
	Bishop
)

var movesByPiece = map[PieceType][]Movement{
	King:   {{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}},
	Bishop: {{-1, -1}, {-1, 1}, {1, -1}, {1, 1}},
	Knight: {{-2, -1}, {-1, -2}, {-2, 1}, {-1, 2}, {2, -1}, {1, -2}, {2, 1}, {1, 2}},
}

func isLongDistancePiece(pieceType PieceType) bool {
	return pieceType == Bishop
}

type ChessPiece struct {
	isEnemy    bool
	pieceType  PieceType
	colorIndex int
}

func (p *ChessPiece) GetMoves(board *Board, current Square) map[Square]bool {
	moves := movesByPiece[p.pieceType]
	result := map[Square]bool{}
	if isLongDistancePiece(p.pieceType) {
		for _, move := range moves {
			next := current
			for {
				next = next.add(move)
				if existing_piece, piece_exists := board.piecesBySquare[next]; piece_exists {
					if existing_piece.isEnemy != p.isEnemy {
						result[next] = true
					}
					break
				}
				if !board.Contains(next) {
					break
				}
				result[next] = true
			}
		}
	} else {
		for _, move := range moves {
			next := current.add(move)
			if existing_piece, piece_exists := board.piecesBySquare[next]; piece_exists {
				if existing_piece.isEnemy != p.isEnemy {
					result[next] = true
				}
			} else if board.Contains(next) {
				result[next] = true
			}
		}
	}
	return result
}
