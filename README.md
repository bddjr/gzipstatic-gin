# gzipstatic-gin

使用静态 gzip 或 br 压缩文件响应，减少服务器性能占用。

Use static gzip or br compression file response to reduce server performance consumption.

---

## Get

```
go get github.com/bddjr/gzipstatic-gin
```

---

## Example

### NoRoute

```go
noRoute := func(ctx *gin.Context) {
    f, _ := os.ReadFile("frontend/dist/404.html")
    ctx.Data(404, gin.MIMEHTML, f)
}
Router.NoRoute(noRoute)

gzipstatic.NoRoute = noRoute
```

### Static

```go
// router.Static("/", "frontend/dist")
gzipstatic.Static(router, "/", "frontend/dist")
```

### StaticFile

```go
// router.StaticFile("/", "frontend/dist/index.html")
gzipstatic.StaticFile(router, "/", "frontend/dist/index.html")
```

### File

```go
// ctx.File("frontend/dist/index.html")
gzipstatic.File(ctx, "frontend/dist/index.html")
```

### StaticFS

```go
// router.StaticFS("/", "/", fs)
gzipstatic.StaticFS(router, "/", "/", fs)
```

### StaticFileFS

```go
// router.StaticFileFS("/", "index.html", fs)
gzipstatic.StaticFileFS(router, "/", "index.html", fs)
```

### FileFromFS

```go
// ctx.FileFromFS("index.html", fs)
gzipstatic.FileFromFS(ctx, "index.html", fs)
```

### ExtFillter

```go
gzipstatic.ExtFillter = regexp.MustCompile(`\.(html|htm|js|json|css)$`)
```

### EncodeList

```go
// Priority from high to low
gzipstatic.EncodeList = []gzipstatic.EncodeListItem{
    {
        name: "br",
        ext:  ".br",
    }, {
        name: "gzip",
        ext:  ".gz",
    },
}
```

---

## Response Header Characteristics

X-Content-Encoding-By: gzipstatic-gin

---

## Source Code

[gzipstatic.go](gzipstatic.go)

---

## License

[BSD-3-clause license](LICENSE.txt)
