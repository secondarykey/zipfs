# zipfs is zip filesystem

"zipfs" can handle the contents of a ZIP file as fs.FS.

The reason I created it was because I was building a static site on Google App Engine, but archiving it by year resulted in a large number of files and deployment took a long time, and since past archives received few hits, I thought it would be suitable for a site where slow processing was acceptable.

## zipfs.New()

```
dir := os.DirFS("./")
z,err := zipfs.New(dir,"examples.zip")
```

### zipfs.NewZipFile()

You can also open it by file name.

```
z,err := zipfs.NewZipFile("./examples.zip")
```

## ZipFS::Open()

Just like fs.FS, you can obtain a file with Open().

```
zf,err := z.Open("index.html")
logo,err := z.Open("assets/zipfs.webp")
```

ZipFS does not unpack the zip file in memory until it is opened.

## ZipFS::Release()

Once a file has been Opened, the Zip file is expanded in memory.
This can be released by calling Release().

```
z,err := zipfs.New(dir,"examples.zip")
z.Release()
zf,err := z.Open("index.html") // This works
```

You don't need to worry about Open() it again.

## ZipFS::Init()

If you don't want to initialize it the first time you open it, you can manually expand the memory by calling Init().

```
z,err := zipfs.NewZipFile("../examples/examples.zip")
err = z.Init()

zf,err := z.Open("index.html") // Fast because expansion is done at initialization
```

## Handler sample

I have implemented a simple way to use http.Handler.

```
> go run examples/main.go
```

Visit http://localhost:8080/examples/ and you will see the page.

```
    http.Handle("/examples/",
        http.StripPrefix("/examples/", http.FileServerFS(z)))
```

In this example, ZipFS is released every 5 seconds.

```
    go func() {
         for range time.Tick(5 * time.Second) {
             z.Release()
         }
    }()
```

The memory is displayed every second, so the memory increases when accessed, but you can see that the memory is released within 5 seconds.



## Tips for adjusting the name 

Although this is not a zipfs feature, I will describe how to adjust the names that are often used.

### Adjusting fs.FS

When accessing "./test/index.html", the name should be just "index.html"

```
f := fs.Sub(zf,"./test")
```

### Adjusting http.Handler

The URL contains "examples", but the FS does not contain examples.

```
http.Handle("/examples/",http.StripPrefix("/examples/",handler))
```

When "/examples/index.html" is accessed, index.html will be accessed on the FS side.

## Usage 


