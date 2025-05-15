package zipfs

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"path/filepath"
)

type dirEntry struct {
	f *zip.File
}

func newDir(f *zip.File) *dirEntry {
	var d dirEntry
	d.f = f
	return &d
}

func (d dirEntry) Name() string {
	n := filepath.Base(d.f.Name)
	return n
}

func (d dirEntry) IsDir() bool {
	return d.f.Mode().IsDir()
}

func (d dirEntry) Type() fs.FileMode {
	return d.f.Mode().Type()
}

func (d dirEntry) Info() (fs.FileInfo, error) {
	return d.f.FileInfo(), nil
}

type rootFile struct {
	info *rootInfo
}

type rootInfo struct {
	fs.FileInfo
}

func newRoot(z fs.FileInfo) *rootFile {
	var root rootFile
	var info rootInfo
	info.FileInfo = z

	root.info = &info
	return &root
}

func (r *rootFile) Read(data []byte) (int, error) {
	return 0, fmt.Errorf("ZIP File root")
}

func (r *rootFile) Stat() (fs.FileInfo, error) {
	return r.info, nil
}

func (r *rootFile) Close() error {
	return fmt.Errorf("ZIP File root")
}

func (r *rootInfo) Name() string {
	return ""
}

func (r *rootInfo) IsDir() bool {
	return true
}
