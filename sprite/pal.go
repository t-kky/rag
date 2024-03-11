package sprite

import (
	"encoding/binary"
	"image/color"
	"io"
)

type Palette struct {
	Colors [256]color.RGBA
}

func LoadPal(reader io.Reader) Palette {
	p := Palette{}
	binary.Read(reader, binary.LittleEndian, &p.Colors)
	return p
}
