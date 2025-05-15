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
	mu     sync.Mutex
}

func NewZipFile(name string) (*ZipFS, error) {

	dir := filepath.Dir(name)
	base := filepath.Base(name)

	d := os.DirFS(dir)

	return New(d, base)
}

func New(dir fs.FS, name string) (*ZipFS, error) {

	var err error
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
	err := f.Init()
	if err != nil {
		return nil, xerrors.Errorf("Init() error: %w", err)
	}

	file, err := f.reader.Open(name)
	if err != nil {
		return nil, xerrors.Errorf("reader Open() error: %w", err)
	}

	return file, nil
}

// fs.ReadFileFS
func (f *ZipFS) ReadFile(name string) ([]byte, error) {

	fp, err := f.Open(name)
	if err != nil {
		return nil, xerrors.Errorf("fs.Open() error: %w", err)
	}
	defer fp.Close()

	data, err := io.ReadAll(fp)
	if err != nil {
		return nil, xerrors.Errorf("io.ReadAll() error: %w", err)
	}
	return data, nil
}

// fs.ReadDirFS
func (f *ZipFS) ReadDir(name string) ([]fs.DirEntry, error) {

	info, err := f.Stat(name)
	if err != nil {
		return nil, xerrors.Errorf("fs.Stat() error: %w", err)
	}

	if !info.IsDir() {
		return nil, xerrors.Errorf("%s not Directory", name)
	}

	return f.reader.readDir(name)
}

// fs.StatFS
func (f *ZipFS) Stat(name string) (fs.FileInfo, error) {

	fp, err := f.Open(name)
	if err != nil {
		return nil, xerrors.Errorf("fs.Open() error: %w", err)
	}
	defer fp.Close()

	return fp.Stat()
}

// fs.GlobFS
func (f *ZipFS) Glob(pattern string) ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	err := f.Init()
	if err != nil {
		return nil, xerrors.Errorf("Init() error: %w", err)
	}

	return f.reader.glob(pattern)
}

func (f *ZipFS) Init() error {

	if !f.reader.IsClose() {
		return nil
	}

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

// fs.SubFS
//func (f *ZipFS) Sub(name string) (fs.FS,error)

// fs.WalkDirFunc
// func(path string,d DirEntry,err error) error
