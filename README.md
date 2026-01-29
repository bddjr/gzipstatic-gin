# gzipstatic-gin

自动使用最小的压缩文件响应，减少服务器性能消耗。

Automatically use the smallest compressed file response to reduce server performance consumption.

---

## Get

```
go get -u github.com/bddjr/gzipstatic-gin
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

### EnableDebugHeader

```go
// X-Encoding-By: gzipstatic-gin
gzipstatic.EnableDebugHeader = true
```

---

## Test

Dependencies: Git, Go, Node.js

```
git clone https://github.com/bddjr/gzipstatic-gin
cd gzipstatic-gin
cd testdata
cd vite-project
npm i -g pnpm
pnpm i
pnpm build
cd ..
go build
./testdata
```
