package zipfs_test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/secondarykey/zipfs"
)

//go:embed examples/embed/*
var emb embed.FS
var dir fs.FS

func init() {
	var err error
	dir, err = fs.Sub(emb, "examples/embed")
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {

	z, err := zipfs.New(dir, "web.zip")
	if err != nil {
		t.Fatalf("zipfs.New() error: %v", err)
	}

	_, err = z.Open("index.html")
	if err != nil {
		t.Errorf("zipfs Open() error: %v", err)
	}

	_, err = z.Open("assets/zipfs.webp")
	if err != nil {
		t.Errorf("zipfs Open() error: %v", err)
	}

	_, err = z.Open("notfound.jpg")
	if err == nil {
		t.Errorf("not found is error")
	}

	err = z.Release()
	if err != nil {
		t.Errorf("zipfs Release() error: %v", err)
	}

	//It work
	_, err = z.Open("index.html")
	if err != nil {
		t.Errorf("zipfs re Open() error: %v", err)
	}
}

func TestNewZipFile(t *testing.T) {
	z, err := zipfs.NewZipFile("examples/embed/web.zip")
	if err != nil {
		t.Fatalf("zipfs.NewZipFile() error: %v", err)
	}

	_, err = z.Open("index.html")
	if err != nil {
		t.Errorf("zipfs Open() error: %v", err)
	}

}

func TestReadDir(t *testing.T) {
	z, err := zipfs.New(dir, "web.zip")
	if err != nil {
		t.Fatalf("zipfs.New() error: %v", err)
	}

	dirs, err := z.ReadDir("assets")
	if err != nil {
		t.Fatalf("zipfs.ReadDir() error: %v", err)
	} else if len(dirs) != 2 {
		t.Errorf("assets children 2 not %d", len(dirs))
	}

	dirs, err = z.ReadDir("")
	if err != nil {
		t.Fatalf("zipfs.ReadDir() error: %v", err)
	} else if len(dirs) != 2 {
		t.Errorf("root 2 not %v", len(dirs))
	}
}

func TestGlob(t *testing.T) {
	z, err := zipfs.New(dir, "web.zip")
	if err != nil {
		t.Fatalf("zipfs.New() error: %v", err)
	}

	files, err := z.Glob("*.webp")
	if err != nil {
		t.Errorf("Glob(*.webp) error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Glob(*.webp) length error: %v", len(files))
	}

	if files[0] != "assets/zipfs.webp" {
		t.Errorf("Glob(*.webp) name error: %v", files[0])
	}

	files, err = z.Glob("*.*")
	if err != nil {
		t.Errorf("Glob(*.*) error: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("Glob(*.*) length error: %v", len(files))
	}
}

//func TestReadFile(t *testing.T) {
//func TestStat(t *testing.T) {

func TestInit(t *testing.T) {

	z, err := zipfs.New(dir, "web.zip")
	if err != nil {
		t.Fatalf("zipfs.New() error: %v", err)
	}

	err = z.Init()
	if err != nil {
		t.Errorf("zipfs Init() error: %v", err)
	}

	_, err = z.Open("assets/styles.css")
	if err != nil {
		t.Errorf("zipfs Open() error: %v", err)
	}

	err = z.Release()
	if err != nil {
		t.Errorf("zipfs Release() error: %v", err)
	}
	err = z.Init()
	if err != nil {
		t.Errorf("zipfs Init() error: %v", err)
	}

	_, err = z.Open("assets/styles.css")
	if err != nil {
		t.Errorf("zipfs Open() error: %v", err)
	}
}

func TestWalkDir(t *testing.T) {

	z, err := zipfs.New(dir, "web.zip")
	if err != nil {
		t.Fatalf("zipfs.New() error: %v", err)
	}

	cnt := 0
	fs.WalkDir(z, "", func(path string, d fs.DirEntry, err error) error {
		cnt++
		return nil
	})

	if cnt != 5 {
		t.Errorf("WalkDir cnt error: want 5 got %d ", cnt)
	}

}
