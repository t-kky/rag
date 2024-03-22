package grf

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io/fs"
	"os"
	"strings"

	"golang.org/x/text/encoding/korean"
)

const HEADER_SIZE = 46

type grfFile struct {
	name   string
	header header
	files  map[string]File
}

type header struct {
	Signature         [16]byte
	EncryptionKey     [14]byte
	FileTableOffset   uint32
	ScramblingSeed    uint32
	ScramledFileCount uint32
	Version           uint32
}

type FS struct {
	grf *grfFile
}

func NewFS(name string) *FS {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	order := binary.LittleEndian

	grf := grfFile{name: name, files: make(map[string]File)}

	binary.Read(file, order, &grf.header)

	_, err = file.Seek(int64(grf.header.FileTableOffset)+HEADER_SIZE, 0)
	if err != nil {
		panic(err)
	}

	var compressedSize uint32
	binary.Read(file, order, &compressedSize)
	file.Seek(4, 1)

	buffer := make([]byte, int(compressedSize))
	binary.Read(file, order, &buffer)

	unpacked, err := zlib.NewReader(bytes.NewReader(buffer))
	if err != nil {
		panic(err)
	}

	filecount := int(grf.header.ScramledFileCount - grf.header.ScramblingSeed - 7)

	for i := 0; i < filecount; i++ {
		var fname []byte
		for {
			char := make([]byte, 1)
			unpacked.Read(char)
			if char[0] == 0 {
				break
			}
			fname = append(fname, char[0])
		}
		decoder := korean.EUCKR.NewDecoder()
		fname, _ = decoder.Bytes(fname)
		normalized := normalize(string(fname))
		fileData := data{}
		binary.Read(unpacked, order, &fileData)
		file := File{normalized, grf, nil, fileData}
		grf.files[file.name] = file
	}

	return &FS{&grf}
}

func (fs *FS) Open(name string) (fs.File, error) {
	file, ok := fs.grf.files[normalize(name)]
	if !ok {
		return nil, errors.New("does not exist")
	}

	return &file, nil
}

func normalize(s string) string {
	s = strings.ToLower(string(s))
	s = strings.ReplaceAll(s, "\\", "/")
	s = strings.ReplaceAll(s, "//", "/")
	return s
}
