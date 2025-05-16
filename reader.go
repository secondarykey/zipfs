package zipfs

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

type reader struct {
	dir  fs.FS
	name string

	root    *rootFile
	zipFile fs.File

	data      []byte
	zipReader *zip.Reader
}

func NewReader(dir fs.FS, name string) (*reader, error) {

	info, err := fs.Stat(dir, name)
	if err != nil {
		return nil, xerrors.Errorf("FS Stat() error: %w", err)
	}

	var r reader
	r.root = newRoot(info)
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
	if name == "" {
		return r.root, nil
	}

	zf, err := r.zipReader.Open(name)
	if err != nil {
		return nil, xerrors.Errorf("zipReader.Open() error: %w", err)
	}

	return zf, nil
}

func (r *reader) readDir(name string) ([]fs.DirEntry, error) {

	files := r.zipReader.File
	var rtn []fs.DirEntry

	for _, f := range files {
		n := strings.ReplaceAll(f.Name, "\\", "/")
		if isChild(n, name) {
			rtn = append(rtn, newDir(f))
		}
	}
	return rtn, nil
}

func isChild(f1, f2 string) bool {

	idx := strings.Index(f1, f2)
	if idx != 0 {
		return false
	}

	f := f1
	if f2 != "" {
		f = strings.Replace(f1, f2, "", 1)
	} else {
		f = "/" + f1
	}

	if f == "/" {
		return false
	}

	idx = strings.LastIndex(f, "/")
	if idx > 0 {
		if len(f) != idx+1 {

			return false
		}
	}

	return true
}

func (r *reader) glob(ptn string) ([]string, error) {

	files := r.zipReader.File
	var rtn []string
	var err error

	for _, f := range files {

		name := strings.ReplaceAll(f.Name, "\\", "/")
		m, me := filepath.Match(ptn, name)
		errors.Join(err, me)
		if m && !f.FileInfo().IsDir() {
			rtn = append(rtn, name)
		}
	}
	return rtn, err
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
