package grf

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"os"
)

type File struct {
	name   string
	grf    grfFile
	reader io.ReadCloser
	data
}

type data struct {
	CompressedSize   uint32
	ByteAlignedSize  uint32
	DecompressedSize uint32
	Type             uint8
	Offset           uint32
}

func (f *File) Close() error {
	if f.reader == nil {
		return nil
	}

	return f.reader.Close()
}

func (f *File) Read(p []byte) (int, error) {
	if f.reader == nil {
		file, err := os.Open(f.grf.name)
		if err != nil {
			return 0, errors.New("error opening grf")
		}
		defer file.Close()

		content := make([]byte, f.CompressedSize)
		file.Seek(int64(f.Offset)+HEADER_SIZE, 0)
		binary.Read(file, binary.LittleEndian, content)

		reader, err := zlib.NewReader(bytes.NewReader(content))
		if err != nil {
			return 0, errors.New("unable to read file")
		}

		f.reader = reader
	}

	return f.reader.Read(p)
}

func (f *File) Stat() (fs.FileInfo, error) {
	return FileInfo{}, nil
}
