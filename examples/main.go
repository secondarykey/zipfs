package main

import (
	"github.com/secondarykey/zipfs"

	"embed"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

//go:embed embed/*
var emb embed.FS

func main() {

	z, err := zipfs.New(emb, "embed/web.zip")
	if err != nil {
		panic(err)
	}

	go func() {
		for range time.Tick(1 * time.Second) {
			printMemory()
		}
	}()

	go func() {
		for range time.Tick(5 * time.Second) {
			z.Release()
		}
	}()

	serv := "localhost:8080"
	http.Handle("/examples/",
		http.StripPrefix("/examples/", http.FileServerFS(z)))
	fmt.Printf("Serv:%s\n", serv)
	http.ListenAndServe(serv, nil)
}

func printMemory() {
	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Printf("%s ==== \n", time.Now().Format(time.DateTime))
	fmt.Printf("  %-8s: %8dBytes\n", "Sys", memStats.Sys)
	fmt.Printf("  %-8s: %8dBytes\n", "Alloc", memStats.Alloc)
	fmt.Printf("  %-8s: %8dBytes\n", "Heap", memStats.HeapAlloc)
	//fmt.Printf("  %-10s: %v bytes\n", "Total", memStats.TotalAlloc)
}
