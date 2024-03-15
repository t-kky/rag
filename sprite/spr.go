package sprite

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
)

type RGBABitmap struct {
	Width  uint16
	Height uint16
	Data   []color.RGBA
}

type Header struct {
	Signature          [2]byte
	Version            Version
	IndexedBitmapCount uint16
	RGBABitmapCount    uint16
}

type Version struct {
	Minor uint8
	Major uint8
}

type Sprite struct {
	Header         Header
	IndexedBitmaps []IndexedBitmap
	RGBABitmaps    []RGBABitmap
	Palette        [256]color.RGBA
}

func (s *Sprite) Version() string {
	return fmt.Sprintf("%d.%d", s.Header.Version.Major, s.Header.Version.Minor)
}

func (s *Sprite) IndexedBitmapFrame(frame int) *image.RGBA {
	bitmap := s.IndexedBitmaps[frame]

	width := int(bitmap.Width)
	height := int(bitmap.Height)

	topLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	decoded := bitmap.decode()
	for x := 0; x < len(decoded); x++ {
		c := s.Palette[decoded[x]]
		var a uint8
		if decoded[x] == 0 {
			a = 0
		} else {
			a = 255
		}
		img.Set(x%width, x/width, color.RGBA{c.R, c.G, c.B, a})
	}

	return img
}

func ReadSpr(reader io.Reader) Sprite {
	sprite := Sprite{}
	endian := binary.LittleEndian

	binary.Read(reader, endian, &sprite.Header)

	for i := 0; i < int(sprite.Header.IndexedBitmapCount); i++ {
		indexedBitmap := IndexedBitmap{}
		binary.Read(reader, endian, &indexedBitmap.Width)
		binary.Read(reader, endian, &indexedBitmap.Height)
		binary.Read(reader, endian, &indexedBitmap.RLELength)
		indexedBitmap.RLEData = make([]uint8, indexedBitmap.RLELength)
		binary.Read(reader, endian, &indexedBitmap.RLEData)
		sprite.IndexedBitmaps = append(sprite.IndexedBitmaps, indexedBitmap)
	}

	for i := 0; i < int(sprite.Header.RGBABitmapCount); i++ {
		rgbaBitmap := RGBABitmap{}
		binary.Read(reader, endian, &rgbaBitmap.Width)
		binary.Read(reader, endian, &rgbaBitmap.Height)
		rgbaBitmap.Data = make([]color.RGBA, rgbaBitmap.Width*rgbaBitmap.Height)
		binary.Read(reader, endian, &rgbaBitmap.Data)
		sprite.RGBABitmaps = append(sprite.RGBABitmaps, rgbaBitmap)
	}

	binary.Read(reader, endian, &sprite.Palette)

	return sprite
}
