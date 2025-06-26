package bntrainer

import "errors"

type Square struct {
	x, y int
}

type Movement Square

func (s Square) add(m Movement) Square {
	return Square{s.x + m.x, s.y + m.y}
}

type Board struct {
	size           int
	piecesBySquare map[Square]ChessPiece
}

func (b *Board) Contains(s Square) bool {
	return s.x >= 0 && s.x < b.size && s.y >= 0 && s.y < b.size
}

func (b *Board) Move(fromSquare, toSquare Square, isPlayerMove bool) error {
	if fromSquare == toSquare {
		return errors.New("invalid move: from and to squares are the same")
	}
	fromPiece, fromPieceExists := b.piecesBySquare[fromSquare]
	if !fromPieceExists {
		return errors.New("invalid move: no piece on from square")
	}
	if fromPiece.isEnemy == isPlayerMove {
		return errors.New("invalid move: wrong color piece tried to move")
	}
	if fromPiece.pieceType == King {
		for square, piece := range b.piecesBySquare {
			if fromPiece.isEnemy != piece.isEnemy {
				delete(b.piecesBySquare, fromSquare)
				stepsIntoCheck, _ := piece.GetMoves(b, square)[toSquare]
				if stepsIntoCheck {
					b.piecesBySquare[fromSquare] = fromPiece
					return errors.New("invalid move: tried to step into check")
				}
				b.piecesBySquare[fromSquare] = fromPiece
			}
		}
	}
	canMove, _ := fromPiece.GetMoves(b, fromSquare)[toSquare]
	if !canMove {
		return errors.New("invalid move: piece cannot move this way")
	}
	toPiece, toPieceExists := b.piecesBySquare[toSquare]
	if toPieceExists {
		if fromPiece.isEnemy == toPiece.isEnemy {
			return errors.New("invalid move: trying to capture own piece")
		}
	}
	delete(b.piecesBySquare, fromSquare)
	b.piecesBySquare[toSquare] = fromPiece
	return nil
}
