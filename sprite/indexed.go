package sprite

type IndexedBitmap struct {
	Width     uint16
	Height    uint16
	RLELength uint16
	RLEData   []uint8
}

func (b *IndexedBitmap) decode() []uint8 {
	// TODO: decoded := make([]uint8, b.Width * b.Height)
	var decoded []uint8

	for i := 0; i < int(b.RLELength); i++ {
		idx := b.RLEData[i]

		if idx == 0 {
			n := int(b.RLEData[i+1])
			for j := 0; j < n; j++ {
				decoded = append(decoded, 0)
			}

			i = i + 1
		} else {
			decoded = append(decoded, idx)
		}
	}

	return decoded
}
