package zipfs

import (
	"archive/zip"
	"io"
	"io/fs"
	"strings"

	"golang.org/x/xerrors"
)

type reader struct {
	dir  fs.FS
	name string

	zipFile fs.File

	data      []byte
	zipReader *zip.Reader
}

func NewReader(dir fs.FS, name string) (*reader, error) {
	var r reader
	r.dir = dir
	r.name = name
	return &r, nil
}

func (r *reader) init() error {

	if !r.IsClose() {
		r.Close()
	}

	file, err := r.dir.Open(r.name)
	if err != nil {
		return xerrors.Errorf("fs.Open() error: %w", err)
	}
	r.zipFile = file

	data, err := io.ReadAll(file)
	if err != nil {
		return xerrors.Errorf("io.ReadAll() error: %w", err)
	}
	r.data = data

	r.zipReader, err = zip.NewReader(r, int64(len(data)))
	if err != nil {
		return xerrors.Errorf("zip.NewReader() error: %w", err)
	}
	return nil
}

func (r *reader) Open(name string) (fs.File, error) {
	return r.zipReader.Open(name)
}

type dirEntry struct {
	f *zip.File
}

func newDir(f *zip.File) *dirEntry {
	var d dirEntry
	d.f = f
	return &d
}

func (d dirEntry) Name() string {
	return d.f.Name
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

func (r *reader) readDir(name string) ([]fs.DirEntry, error) {

	files := r.zipReader.File
	var rtn []fs.DirEntry

	for _, f := range files {
		if strings.HasPrefix(f.Name, name) {
			if f.Name != name {
				rtn = append(rtn, newDir(f))
			}
		}
	}
	return rtn, nil
}

func (r *reader) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= int64(len(r.data)) {
		return 0, io.EOF
	}
	n = copy(p, r.data[off:])
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (r *reader) Close() error {

	if r == nil {
		return nil
	}

	r.data = nil
	if r.zipReader != nil {
		//r.zipReader.Close()
	}

	if r.zipFile != nil {
		return r.zipFile.Close()
	}
	return nil
}

func (r *reader) IsClose() bool {
	if r == nil {
		return true
	}
	return r.data == nil
}
