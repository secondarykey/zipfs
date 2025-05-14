package zipfs_test

import (
	"testing"

	"github.com/secondarykey/zipfs"

	"embed"
	"io/fs"
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
	}

	if len(dirs) != 2 {
		t.Errorf("2 not %v", len(dirs))
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
