package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/t-kky/rag/sprite"
)

var head int
var palette int
var direction int

func init() {
	flag.IntVar(&head, "head", 1, "head id")
	flag.IntVar(&palette, "palette", 1, "pal id")
	flag.IntVar(&direction, "direction", 0, "direction")
	flag.Parse()
}

func headSprite(headID int, paletteID int, direction int) *image.RGBA {
	spritePath := fmt.Sprintf("resources/data/sprite/인간족/머리통/남/%d_남.spr", headID)
	palettePath := fmt.Sprintf("resources/data/palette/머리/머리%d_남_%d.pal", headID, paletteID)

	spriteFile, _ := os.Open(spritePath)
	defer spriteFile.Close()
	spr := sprite.Read(spriteFile)
	fmt.Printf("%+v\n", spr.Version())

	paletteFile, _ := os.Open(palettePath)
	defer paletteFile.Close()
	palette := sprite.LoadPal(paletteFile)

	spr.Palette = palette.Colors

	return spr.IndexedBitmapFrame(direction)
}

func main() {
	img := headSprite(head, palette, direction)
	imgFile, _ := os.Create("output.png")
	defer imgFile.Close()
	png.Encode(imgFile, img)
}
