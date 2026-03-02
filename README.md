# gzipstatic-gin

自动使用压缩文件响应，减少服务器性能消耗。

Automatically use compressed file responses, reducing server performance consumption.

Default encode:
> br .br  
> zstd .zst  
> gzip .gz  

---

## Get

```
go get -u github.com/bddjr/gzipstatic-gin@latest
```

---

## Example

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

### ExtFilterMap

```go
gzipstatic.ExtFilterMap = map[string]struct{}{
	".css":  {},
	".htm":  {},
	".html": {},
	".js":   {},
	".json": {},
	".mjs":  {},
	".svg":  {},
	".wasm": {},
	".xml":  {},
	".yaml": {},
	".yml":  {},
	".toml": {},
}
```

### EncodeNameExtMap

```go
gzipstatic.EncodeNameExtMap = map[string]string{
	"br":   ".br",
	"zstd": ".zst",
	"gzip": ".gz",
}
```

### EncodeNamePriority

```go
// from high to low
gzipstatic.EncodeNamePriority = []string{
	"br",
	"zstd",
	"gzip",
}
```

### EnableDebugHeader

```go
// X-Encoding-By: gzipstatic-gin
gzipstatic.EnableDebugHeader = true
```
