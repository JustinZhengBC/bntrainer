package bntrainer

import (
	"math/rand"
)

type GameType int

const (
	GameType_BN GameType = iota
	GameType_BB
	GameType_NNN
)

func initializeRandomSquares(boardSize int) map[Square]bool {
	randomSquares := map[Square]bool{}
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			randomSquares[Square{i, j}] = true
		}
	}
	return randomSquares
}

func getRandomSquare(availableSquares map[Square]bool) Square {
	randInt := rand.Intn(len(availableSquares))
	for square := range availableSquares {
		if randInt == 0 {
			return square
		}
		randInt--
	}
	panic("Ran out of available squares!")
}

func placeRandomly(board *Board, availableSquares map[Square]bool, piece ChessPiece) {
	square := getRandomSquare(availableSquares)
	board.piecesBySquare[square] = piece
	delete(availableSquares, square)
	for attackedSquare := range piece.GetMoves(board, square) {
		delete(availableSquares, attackedSquare)
	}
}

func GenerateBoard(boardSize int, gameType GameType) *Board {
	board := &Board{boardSize, map[Square]ChessPiece{}}
	availableSquares := initializeRandomSquares(boardSize)
	placeRandomly(board, availableSquares, ChessPiece{false, King, 0})
	switch gameType {
	case GameType_BN:
		placeRandomly(board, availableSquares, ChessPiece{false, Bishop, 1})
		placeRandomly(board, availableSquares, ChessPiece{false, Knight, 2})
		placeRandomly(board, availableSquares, ChessPiece{true, King, -1})
	case GameType_BB:
		availableLightSquares := map[Square]bool{}
		availableDarkSquares := map[Square]bool{}
		for square := range availableSquares {
			if (square.x+square.y)%2 == 0 {
				availableLightSquares[square] = true
			} else {
				availableDarkSquares[square] = true
			}
		}
		placeRandomly(board, availableLightSquares, ChessPiece{false, Bishop, 1})
		placeRandomly(board, availableDarkSquares, ChessPiece{false, Bishop, 3})
		if rand.Intn(2) == 0 {
			placeRandomly(board, availableLightSquares, ChessPiece{true, King, -1})
		} else {
			placeRandomly(board, availableLightSquares, ChessPiece{true, King, -1})
		}
	case GameType_NNN:
		placeRandomly(board, availableSquares, ChessPiece{false, Knight, 1})
		placeRandomly(board, availableSquares, ChessPiece{false, Knight, 2})
		placeRandomly(board, availableSquares, ChessPiece{false, Knight, 3})
		placeRandomly(board, availableSquares, ChessPiece{true, King, -1})
	}
	return board
}
