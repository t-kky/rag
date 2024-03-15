package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/t-kky/rag/sprite"
	"github.com/t-kky/rag/grf"
)

var (
	head      int
	palette   int
	direction int
	grfPath   string
)

func init() {
	flag.IntVar(&head, "head", 1, "head id")
	flag.IntVar(&palette, "palette", 1, "pal id")
	flag.IntVar(&direction, "direction", 0, "direction")
	flag.StringVar(&grfPath, "grf", "", "path to grf")
	flag.Parse()
}

func headSprite(headID int, paletteID int, direction int) *image.RGBA {
	spritePath := fmt.Sprintf("data/sprite/인간족/머리통/남/%d_남.spr", headID)
	palettePath := fmt.Sprintf("data/palette/머리/머리%d_남_%d.pal", headID, paletteID)

	data := grf.Read(grfPath)
	spriteReader := data.ReadFile(spritePath)
	defer spriteReader.Close()

	spr := sprite.ReadSpr(spriteReader)
	// fmt.Printf("%+v\n", spr.Version())

	palReader := data.ReadFile(palettePath)
	defer palReader.Close()
	palette := sprite.LoadPal(palReader)

	spr.Palette = palette.Colors

	return spr.IndexedBitmapFrame(direction)
}

func main() {
	if grfPath == "" {
		fmt.Println("path to grf required, see cli -h")
		os.Exit(1)
	}

	img := headSprite(head, palette, direction)
	imgFile, _ := os.Create("output.png")
	defer imgFile.Close()
	png.Encode(imgFile, img)
}
