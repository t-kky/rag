package grf

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/korean"
)

const HEADER_SIZE = 46

type GRF struct {
	Path string
	Header Header
	Files  []File
}

type Header struct {
	Signature         [16]byte
	EncryptionKey     [14]byte
	FileTableOffset   uint32
	ScramblingSeed    uint32
	ScramledFileCount uint32
	Version           uint32
}

type File struct {
	Name string
	FileData
}

type FileData struct {
	CompressedSize   uint32
	ByteAlignedSize  uint32
	DecompressedSize uint32
	Type             uint8
	Offset           uint32
}

func Read(path string) *GRF {
	grf := GRF{Path: path}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	order := binary.LittleEndian

	binary.Read(file, order, &grf.Header)

	_, err = file.Seek(int64(grf.Header.FileTableOffset)+HEADER_SIZE, 0)
	if err != nil {
		panic(err)
	}

	var compressedSize uint32
	// var uncompressedSize uint32

	binary.Read(file, order, &compressedSize)
	file.Seek(4, 1)
	// binary.Read(file, order, &uncompressedSize)

	buffer := make([]byte, int(compressedSize))
	binary.Read(file, order, &buffer)

	unpacked, err := zlib.NewReader(bytes.NewReader(buffer))
	if err != nil {
		panic(err)
	}

	filecount := int(grf.Header.ScramledFileCount - grf.Header.ScramblingSeed - 7)

	for i := 0; i < filecount; i++ {
		var name []byte
		for {
			char := make([]byte, 1)
			unpacked.Read(char)
			if char[0] == 0 {
				break
			}
			name = append(name, char[0])
		}
		decoder := korean.EUCKR.NewDecoder()
		name, _ = decoder.Bytes(name)
		normalized := normalizedPath(string(name)) 
		fileData := FileData{}
		binary.Read(unpacked, order, &fileData)
		file := File{normalized, fileData}
		grf.Files = append(grf.Files, file)
	}

	return &grf
}

func normalizedPath(s string) string {
	s = strings.ToLower(string(s))
	s = strings.ReplaceAll(s, "\\", "/")
	s = strings.ReplaceAll(s, "//", "/")
	return s
}

func (g *GRF) ReadFile(filename string) io.ReadCloser {


	var grfFile File

	found := false
	for _, f := range g.Files {
		if normalizedPath(f.Name) == normalizedPath(filename) {
			grfFile = f
			found = true
		}
	}
	if found == false {
		panic("not found!")
	}

	file, err := os.Open(g.Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content := make([]byte, grfFile.CompressedSize)
	file.Seek(int64(grfFile.Offset) + HEADER_SIZE, 0)
	binary.Read(file, binary.LittleEndian, content)

	unpacked, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		panic(err)
	}

	return unpacked
}
