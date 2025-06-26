package bntrainer

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	DarkSquareColor                              = color.RGBA{128, 100, 64, 255}
	LightSquareColor                             = color.RGBA{228, 200, 164, 255}
	BackgroundSquareColor                        = color.RGBA{100, 100, 100, 255}
	SelectedSquareColor                          = color.RGBA{167, 167, 167, 255}
	HoveredSquareColor                           = color.RGBA{200, 200, 200, 255}
	CheckmateMessageColor                        = color.RGBA{0, 255, 0, 255}
	DrawMessageColor                             = color.RGBA{255, 0, 0, 255}
	DefaultMessageColor                          = color.RGBA{0, 0, 0, 255}
	TitleFontSize         int                    = 96
	FontSize              int                    = 36
	TextFont              *text.GoTextFaceSource = nil
)

func init() {
	font, error := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if error != nil {
		panic("Failed to read font!")
	}
	TextFont = font
}

type SpriteSheet struct {
	TileWidth                                              int
	King, Bishop, Knight, EnemyKing, Board, SelectedSquare *ebiten.Image
	Marks                                                  [4](*ebiten.Image)
}

func DrawText(image *ebiten.Image, message string, drawPoint image.Point) {
	DrawTextWithSizeAndColor(image, message, drawPoint, FontSize, DefaultMessageColor)
}
func DrawTextWithSize(image *ebiten.Image, message string, drawPoint image.Point, size int) {
	DrawTextWithSizeAndColor(image, message, drawPoint, size, DefaultMessageColor)
}

func DrawTextWithColor(image *ebiten.Image, message string, drawPoint image.Point, color color.RGBA) {
	DrawTextWithSizeAndColor(image, message, drawPoint, FontSize, color)
}

func DrawTextWithSizeAndColor(image *ebiten.Image, message string, drawPoint image.Point, size int, color color.RGBA) {
	face := &text.GoTextFace{
		Source: TextFont,
		Size:   float64(size),
	}
	messageWidth, messageHeight := text.Measure(message, face, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(drawPoint.X-int(messageWidth)/2), float64(drawPoint.Y-int(messageHeight)/2))
	op.ColorScale.ScaleWithColor(color)
	text.Draw(image, message, face, op)
}

func GetSpriteSheet(boardSize int) *SpriteSheet {
	fullPath, fullPathError := filepath.Abs("bntrainer/resources/sprites.png")
	if fullPathError != nil {
		panic("Failed to get sprite file path!")
	}
	sprite_file, sprite_file_read_error := os.Open(fullPath)
	if sprite_file_read_error != nil {
		panic("Failed to read sprite file!")
	}
	sprite_image, _, sprite_image_decode_error := image.Decode(sprite_file)
	if sprite_image_decode_error != nil {
		panic("Failed to decode sprite file!")
	}
	sprites := ebiten.NewImageFromImage(sprite_image)
	tileWidth := sprites.Bounds().Size().Y / 2 // sprite sheet is 4x2
	sprite := func(x, y int) *ebiten.Image {
		x *= tileWidth
		y *= tileWidth
		minPoint := image.Point{x, y}
		x += tileWidth
		y += tileWidth
		maxPoint := image.Point{x, y}
		return sprites.SubImage(image.Rectangle{minPoint, maxPoint}).(*ebiten.Image)
	}
	return &SpriteSheet{
		tileWidth,
		sprite(0, 0),
		sprite(1, 0),
		sprite(2, 0),
		sprite(3, 0),
		getBoardImage(boardSize, tileWidth),
		getSquareImage(tileWidth, SelectedSquareColor),
		[4](*ebiten.Image){sprite(0, 1), sprite(1, 1), sprite(2, 1), sprite(3, 1)},
	}
}

func getBoardImage(size, tileWidth int) *ebiten.Image {
	length := tileWidth * size
	boardImage := ebiten.NewImage(length, length)
	boardImage.Fill(LightSquareColor)
	x := 0
	for x < length {
		y := 0
		for y < length {
			vector.DrawFilledRect(boardImage, float32(x)+float32(tileWidth), float32(y), float32(tileWidth), float32(tileWidth), DarkSquareColor, false)
			vector.DrawFilledRect(boardImage, float32(x), float32(y)+float32(tileWidth), float32(tileWidth), float32(tileWidth), DarkSquareColor, false)
			y += tileWidth * 2
		}
		x += tileWidth * 2
	}
	return boardImage
}

func getSquareImage(tileWidth int, color color.RGBA) *ebiten.Image {
	boxImage := ebiten.NewImage(tileWidth, tileWidth)
	boxImage.Fill(color)
	return boxImage
}
