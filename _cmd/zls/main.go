package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"

	"github.com/secondarykey/zipfs"
	"golang.org/x/xerrors"
)

var filter string

func init() {
	flag.StringVar(&filter, "filter", "*", "File Filter")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Required filename(zip)")
		return
	}
	fn := args[0]

	err := run(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "zls error: %+v", err)
		os.Exit(1)
	}

	os.Exit(0)
	return
}

func run(n string) error {
	zfs, err := zipfs.NewZipFile(n)
	if err != nil {
		return xerrors.Errorf("zipfs.NewFile() error: %w", err)
	}

	files, err := zfs.Glob(filter)
	if err != nil {
		return xerrors.Errorf("zfs.Glob() error: %w", err)
	}

	printFiles(os.Stdout, files)
	return nil
}

func printFiles(w io.Writer, files []string) {

	dirs := createDirs(files)
	for _, d := range dirs {
		fmt.Fprintf(w, "%s\n", d.Name)
		for _, f := range d.Files {
			fmt.Fprintf(w, "  %s\n", f.Name)
		}
	}
}

const RootKey = "<Root>/"

func createDirs(files []string) []*dir {

	dirs := make(map[string]*dir)

	for _, fn := range files {

		key := path.Dir(fn)
		name := path.Base(fn)

		d, ok := dirs[key]
		if !ok {
			var wk dir
			wk.Name = key
			d = &wk
			dirs[key] = d
		}

		var f file
		f.Name = name
		d.Files = append(d.Files, &f)
	}

	var rtn []*dir
	for _, d := range dirs {
		rtn = append(rtn, d)
		//TODO ファイル名ソート
	}

	//ディレクトリ名ソート
	sort.Slice(rtn, func(i, j int) bool {
		d1 := rtn[i]
		d2 := rtn[j]
		return less(d1.Name, d2.Name)
	})

	if len(rtn) > 0 {
		rtn[0].Name = RootKey
	}
	return rtn
}

func less(d1, d2 string) bool {
	if d1 == "./" {
		return true
	} else if d2 == "./" {
		return false
	}

	return d1 < d2

	//s1 := strings.Split(d1, "/")
	//s2 := strings.Split(d2, "/")

	//l1 := len(s1)
	//l2 := len(s2)

	//minLen = l1
}

type dir struct {
	Name  string
	Files []*file
}

type file struct {
	Name string
}
