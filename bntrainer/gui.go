package bntrainer

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var empty = Square{-1, -1}

type GameState = int

const (
	GameState_Ongoing = iota
	GameState_Checkmate
	GameState_Stalemate
	GameState_InsufficientMaterial
)

var (
	board          *Board      = nil
	cursorPosition image.Point = image.Point{0, 0}
	gameType       GameType    = GameType_BN
	isPlayerMove   bool        = true
	selectedSquare Square      = empty
	showHelp       bool        = true
	isDragging     bool        = false
	gameState      GameState   = GameState_Ongoing
)

var (
	sprites                    *SpriteSheet  = GetSpriteSheet(BoardSize)
	startGameButton_BN         *ebiten.Image = buildStartGameButton(GameType_BN, SelectedSquareColor)
	hoveredStartGameButton_BN  *ebiten.Image = buildStartGameButton(GameType_BN, HoveredSquareColor)
	startGameButtonPoint_BN    image.Point   = image.Point{200, 300}
	startGameButton_BB         *ebiten.Image = buildStartGameButton(GameType_BB, SelectedSquareColor)
	hoveredStartGameButton_BB  *ebiten.Image = buildStartGameButton(GameType_BB, HoveredSquareColor)
	startGameButtonPoint_BB    image.Point   = image.Point{ScreenWidth / 2, 300}
	startGameButton_NNN        *ebiten.Image = buildStartGameButton(GameType_NNN, SelectedSquareColor)
	hoveredStartGameButton_NNN *ebiten.Image = buildStartGameButton(GameType_NNN, HoveredSquareColor)
	startGameButtonPoint_NNN   image.Point   = image.Point{ScreenWidth - 200, 300}
	boardPoint                 image.Point   = image.Point{ScreenWidth - sprites.Board.Bounds().Dx()/2 - (ScreenHeight-sprites.Board.Bounds().Dy())/2, ScreenHeight / 2}
	boardTopLeft               image.Point   = image.Point{boardPoint.X - sprites.Board.Bounds().Dx()/2, boardPoint.Y - sprites.Board.Bounds().Dy()/2}
)

func ToggleHelp() {
	showHelp = !showHelp
}

func UpdateCursorPosition(position image.Point) {
	cursorPosition = position
}

func processTitlePageClick() {
	if inImage(startGameButton_BN, startGameButtonPoint_BN, cursorPosition) {
		gameType = GameType_BN
		startNewGame()
	} else if inImage(startGameButton_BB, startGameButtonPoint_BB, cursorPosition) {
		gameType = GameType_BB
		startNewGame()
	} else if inImage(startGameButton_NNN, startGameButtonPoint_NNN, cursorPosition) {
		gameType = GameType_NNN
		startNewGame()
	}
}

func updateGameState() {
	if board == nil {
		gameState = GameState_Ongoing
		return
	}
	if isPlayerMove {
		numBishops := 0
		numKnights := 0
		for _, piece := range board.piecesBySquare {
			if piece.pieceType == Knight {
				numKnights++
			} else if piece.pieceType == Bishop {
				numBishops++
			}
		}
		if numBishops*2+numKnights < 3 {
			gameState = GameState_InsufficientMaterial
		}
	} else {
		enemyKingInCheck := false
		enemyKingSquare := empty
		var availableSquares = map[Square]bool{}
		for square, piece := range board.piecesBySquare {
			if piece.isEnemy && piece.pieceType == King {
				availableSquares = piece.GetMoves(board, square)
				enemyKingSquare = square
			}
		}
		for square, piece := range board.piecesBySquare {
			if !piece.isEnemy {
				for attackedSquare, _ := range piece.GetMoves(board, square) {
					delete(availableSquares, attackedSquare)
					if attackedSquare == enemyKingSquare {
						enemyKingInCheck = true
					}
				}
			}
		}
		if len(availableSquares) == 0 {
			if enemyKingInCheck {
				gameState = GameState_Checkmate
			} else {
				gameState = GameState_Stalemate
			}
		}
	}
}

func processGameClick() {
	if gameState != GameState_Ongoing {
		return
	}
	square := getBoardSquare(board, sprites.TileWidth, boardTopLeft, cursorPosition)
	if square == empty {
		return
	}
	if selectedSquare == empty {
		selectedSquare = square
	} else {
		err := board.Move(selectedSquare, square, isPlayerMove)
		if err == nil {
			selectedSquare = empty
			isPlayerMove = !isPlayerMove
		} else {
			selectedSquare = square
		}
		updateGameState()
	}
}

func ProcessClick() {
	isDragging = true
	if board == nil {
		processTitlePageClick()
	} else {
		processGameClick()
	}
}

func ProcessRelease() {
	isDragging = false
	if board == nil || selectedSquare == empty {
		return
	}
	_, pieceExists := board.piecesBySquare[selectedSquare]
	if !pieceExists {
		return
	}
	processGameClick()
}

func DrawGui(screen *ebiten.Image) {
	screen.Fill(BackgroundSquareColor)
	if board == nil {
		DrawTextWithSize(screen, "BN Trainer", image.Point{ScreenWidth / 2, 100}, TitleFontSize)
		if inImage(startGameButton_BN, startGameButtonPoint_BN, cursorPosition) {
			drawImage(screen, hoveredStartGameButton_BN, startGameButtonPoint_BN)
		} else {
			drawImage(screen, startGameButton_BN, startGameButtonPoint_BN)
		}
		if inImage(startGameButton_BB, startGameButtonPoint_BB, cursorPosition) {
			drawImage(screen, hoveredStartGameButton_BB, startGameButtonPoint_BB)
		} else {
			drawImage(screen, startGameButton_BB, startGameButtonPoint_BB)
		}
		if inImage(startGameButton_NNN, startGameButtonPoint_NNN, cursorPosition) {
			drawImage(screen, hoveredStartGameButton_NNN, startGameButtonPoint_NNN)
		} else {
			drawImage(screen, startGameButton_NNN, startGameButtonPoint_NNN)
		}
	} else {
		boardLength := sprites.Board.Bounds().Dy()
		boardMargin := (ScreenHeight - boardLength) / 2
		drawImage(screen, sprites.Board, image.Point{ScreenWidth - boardLength/2 - boardMargin, ScreenHeight / 2})
		if selectedSquare != empty {
			drawInSquare(screen, sprites.SelectedSquare, selectedSquare)
		}
		for square, piece := range board.piecesBySquare {
			if !isDragging || square != selectedSquare || ((piece.colorIndex == -1) == isPlayerMove) {
				drawInSquare(screen, getSpriteForPiece(piece), square)
			}
		}
		for colorIndex := 0; colorIndex <= 3; colorIndex++ {
			for square, piece := range board.piecesBySquare {
				if showHelp && piece.colorIndex == colorIndex {
					markSprite := getSpriteForColorIndex(piece.colorIndex)
					drawInSquare(screen, markSprite, square)
					for toSquare := range piece.GetMoves(board, square) {
						drawInSquare(screen, markSprite, toSquare)
					}
				}
			}
		}
		if selectedSquare != empty && isDragging {
			piece, pieceExists := board.piecesBySquare[selectedSquare]
			hoveredSquare := getBoardSquare(board, sprites.TileWidth, boardTopLeft, cursorPosition)
			if pieceExists && ((piece.colorIndex != -1) == isPlayerMove) {
				drawImage(screen, getSpriteForPiece(piece), cursorPosition)
				if piece.colorIndex != -1 && piece.GetMoves(board, selectedSquare)[hoveredSquare] {
					markSprite := getSpriteForColorIndex(piece.colorIndex)
					drawInSquareWithFade(screen, markSprite, hoveredSquare, true)
					for toSquare := range piece.GetMoves(board, hoveredSquare) {
						drawInSquareWithFade(screen, markSprite, toSquare, true)
					}
				}
			}
		}
		textX := (ScreenWidth - boardLength - boardMargin) / 2
		DrawText(screen, "R for new game", image.Point{textX, 400})
		DrawText(screen, "X for menu", image.Point{textX, 450})
		DrawText(screen, "Z to toggle color", image.Point{textX, 500})
		gameStateText := ""
		gameStateColor := DefaultMessageColor
		switch gameState {
		case GameState_Checkmate:
			gameStateText = "Checkmate!"
			gameStateColor = CheckmateMessageColor
		case GameState_Stalemate:
			gameStateText = "Stalemate"
			gameStateColor = DrawMessageColor
		case GameState_InsufficientMaterial:
			gameStateText = "Insufficient Material"
			gameStateColor = DrawMessageColor
		}
		if gameStateText == "" {
			if isPlayerMove {
				gameStateText = "White to move"
			} else {
				gameStateText = "Black to move"
			}
		}
		DrawTextWithColor(screen, gameStateText, image.Point{textX, 100}, gameStateColor)
	}
}

func ReturnToMenu() {
	board = nil
	selectedSquare = empty
}

func buildStartGameButton(gameType GameType, backgroundColor color.RGBA) *ebiten.Image {
	buttonWidth := sprites.TileWidth * 2
	button := ebiten.NewImage(buttonWidth, buttonWidth)
	button.Fill(backgroundColor)
	switch gameType {
	case GameType_BN:
		drawImage(button, sprites.Knight, image.Point{sprites.TileWidth / 2, sprites.TileWidth})
		drawImage(button, sprites.Bishop, image.Point{sprites.TileWidth * 3 / 2, sprites.TileWidth})
	case GameType_BB:
		drawImage(button, sprites.Bishop, image.Point{sprites.TileWidth / 2, sprites.TileWidth})
		drawImage(button, sprites.Bishop, image.Point{sprites.TileWidth * 3 / 2, sprites.TileWidth})
	case GameType_NNN:
		drawImage(button, sprites.Knight, image.Point{sprites.TileWidth, sprites.TileWidth / 2})
		drawImage(button, sprites.Knight, image.Point{sprites.TileWidth / 2, sprites.TileWidth * 3 / 2})
		drawImage(button, sprites.Knight, image.Point{sprites.TileWidth * 3 / 2, sprites.TileWidth * 3 / 2})
	}
	return button
}

func drawImage(screen *ebiten.Image, image *ebiten.Image, point image.Point) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(point.X-image.Bounds().Dx()/2), float64(point.Y-image.Bounds().Dy()/2))
	screen.DrawImage(image, op)
}

func inImage(image *ebiten.Image, drawPoint, clickPoint image.Point) bool {
	return abs(drawPoint.X-clickPoint.X) <= image.Bounds().Dx()/2 && abs(drawPoint.Y-clickPoint.Y) <= image.Bounds().Dy()/2
}

func drawInSquareWithFade(screen *ebiten.Image, image *ebiten.Image, square Square, shrink bool) {
	op := &ebiten.DrawImageOptions{}
	if shrink {
		op.ColorScale.ScaleAlpha(0.33)
	}
	op.GeoM.Translate(float64(boardTopLeft.X), float64(boardTopLeft.Y))
	op.GeoM.Translate(float64(square.x*sprites.TileWidth), float64(square.y*sprites.TileWidth))
	screen.DrawImage(image, op)
}

func drawInSquare(screen *ebiten.Image, image *ebiten.Image, square Square) {
	drawInSquareWithFade(screen, image, square, false)
}

func getBoardSquare(board *Board, tileWidth int, drawPoint, clickPoint image.Point) Square {
	square := Square{(clickPoint.X - drawPoint.X) / tileWidth, (clickPoint.Y - drawPoint.Y) / tileWidth}
	if !board.Contains(square) {
		return empty
	}
	return square
}

func HandleStartNewGameShortcut() {
	if board != nil {
		startNewGame()
	}
}

func startNewGame() {
	board = GenerateBoard(BoardSize, gameType)
	isPlayerMove = true
	selectedSquare = empty
	gameState = GameState_Ongoing
}

func getSpriteForPiece(piece ChessPiece) *ebiten.Image {
	switch piece.pieceType {
	case King:
		if piece.isEnemy {
			return sprites.EnemyKing
		} else {
			return sprites.King
		}
	case Knight:
		return sprites.Knight
	case Bishop:
		return sprites.Bishop
	}
	panic("Failed to find sprite for piece!")
}

func getSpriteForColorIndex(colorIndex int) *ebiten.Image {
	if colorIndex < 0 || colorIndex >= len(sprites.Marks) {
		panic("Failed to find sprite for color index!")
	}
	return sprites.Marks[colorIndex]
}
func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}
