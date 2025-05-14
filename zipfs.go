package zipfs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/xerrors"
)

type ZipFS struct {
	reader *reader

	mu sync.Mutex
}

func NewZipFile(name string) (*ZipFS, error) {
	dir := filepath.Dir(name)
	base := filepath.Base(name)

	d := os.DirFS(dir)

	return New(d, base)
}

func New(dir fs.FS, name string) (*ZipFS, error) {

	_, err := fs.Stat(dir, name)
	if err != nil {
		return nil, xerrors.Errorf("FS Stat() error: %w", err)
	}

	var z ZipFS
	z.reader, err = NewReader(dir, name)
	if err != nil {
		return nil, xerrors.Errorf("zipfs.NewReader() error: %w", err)
	}

	return &z, nil
}

// fs.FS
func (f *ZipFS) Open(name string) (fs.File, error) {

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.reader.IsClose() {
		err := f.Init()
		if err != nil {
			return nil, xerrors.Errorf("Init() error: %w", err)
		}
	}

	file, err := f.reader.Open(name)
	if err != nil {
		return nil, xerrors.Errorf("reader Open() error: %w", err)
	}
	return file, nil
}

// fs.ReadFileFS
func (f *ZipFS) ReadFile(name string) ([]byte, error) {

	f, err := f.Open(name)
	if err != nil {
		return nil, xerrors.Errorf("fs.Open() error: %w", err)
	}
	defer f.Close()

	reader, err := f.Open()
	if err != nil {
		return nil, xerrors.Errorf("file.Open() error: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, xerrors.Errorf("io.ReadAll() error: %w", err)
	}

	return nil
}

func (f *ZipFS) Init() error {

	err := f.reader.init()
	if err != nil {
		return xerrors.Errorf("reader init() error: %w", err)
	}
	return nil
}

func (f *ZipFS) Release() error {

	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.reader.IsClose() {
		return f.reader.Close()
	}
	return nil
}

// fs.GlobFS
//func (f *ZipFS) Glob(pattern string) ([]string,error)

// fs.ReadDirFS
//func (f *ZipFS) ReadDir(name string) ([]DirEntry,error)

// fs.SubFS
//func (f *ZipFS) Sub(name string) (fs.FS,error)

// fs.WalkDirFunc
// func(path string,d DirEntry,err error) error
