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
gzipstatic.EncodeList = []*gzipstatic.EncodeListItem{
    {
        name: "br",
        ext:  ".br",
    }, {
        name: "gzip",
        ext:  ".gz",
    },
}
```

### EnableDebugHeader

```go
// Encoding-By: gzipstatic-gin
gzipstatic.EnableDebugHeader = true
```

---

## Source Code

[gzipstatic.go](gzipstatic.go)

---

## Reference

https://developer.mozilla.org/docs/Web/HTTP/Headers/Content-Encoding  
https://developer.mozilla.org/docs/Web/HTTP/Headers/Accept-Encoding  
https://developer.mozilla.org/docs/Web/HTTP/Headers/Vary  
https://github.com/gin-gonic/gin  
https://github.com/BCSPanel/BCSPanel/blob/main/src/httprouter/init.go  
https://github.com/lpar/gzipped  
https://github.com/nanmu42/gzip  
https://github.com/vbenjs/vite-plugin-compression

---

## License

[BSD-3-clause license](LICENSE.txt)
