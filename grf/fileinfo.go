package grf

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	mtime time.Time
	name  string
	size  int64
	mode  fs.FileMode
}

func (fi FileInfo) Name() string {
	return fi.name
}

func (fi FileInfo) Size() int64 {
	return fi.size
}

func (fi FileInfo) Mode() fs.FileMode {
	return 0
}

func (fi FileInfo) ModTime() time.Time {
	return fi.mtime
}

func (fi FileInfo) IsDir() bool {
	return false
}

func (fi FileInfo) Sys() interface{} {
	return nil
}
