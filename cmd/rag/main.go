package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/t-kky/rag/grf"
	"github.com/t-kky/rag/sprite"
)

var (
	head       int
	palette    int
	direction  int
	grfPath    string
	outputPath string
)

func init() {
	flag.IntVar(&head, "head", 1, "head id")
	flag.IntVar(&palette, "palette", 1, "pal id")
	flag.IntVar(&direction, "direction", 0, "direction")
	flag.StringVar(&grfPath, "grf", "", "path to grf")
	flag.StringVar(&outputPath, "output", "output.png", "output path")
	flag.Parse()
}

func headSprite(headID int, paletteID int, direction int) *image.RGBA {
	spritePath := fmt.Sprintf("data/sprite/인간족/머리통/남/%d_남.spr", headID)
	palettePath := fmt.Sprintf("data/palette/머리/머리%d_남_%d.pal", headID, paletteID)

	fs := grf.NewFS(grfPath)
	spriteFile, _ := fs.Open(spritePath)
	defer spriteFile.Close()

	spr := sprite.ReadSpr(spriteFile)
	// fmt.Printf("%+v\n", spr.Version())

	palFile, _ := fs.Open(palettePath)
	defer palFile.Close()
	palette := sprite.LoadPal(palFile)

	spr.Palette = palette.Colors

	return spr.IndexedBitmapFrame(direction)
}

func main() {
	if grfPath == "" {
		fmt.Println("path to grf required, see cli -h")
		os.Exit(1)
	}

	img := headSprite(head, palette, direction)
	imgFile, err := os.Create(outputPath)
	if err != nil {
		panic("unable to open output")
	}
	defer imgFile.Close()
	png.Encode(imgFile, img)
}
