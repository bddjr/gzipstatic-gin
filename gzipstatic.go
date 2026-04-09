// https://github.com/bddjr/gzipstatic-gin
package gzipstatic

import (
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"

	"github.com/gin-gonic/gin"
)

var ExtFilterMap = map[string]struct{}{
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

var EncodeNameExtMap = map[string]string{
	"br": ".br",
	// "zstd": ".zst",
	"gzip": ".gz",
}

// from high to low
var EncodeNamePriority = []string{
	"br",
	// "zstd",
	"gzip",
}

// X-Encoding-By: gzipstatic-gin
var EnableDebugHeader = true

var offset_RouterGroup_engine = func() uintptr {
	sf, ok := reflect.TypeFor[gin.RouterGroup]().FieldByName("engine")
	if !ok {
		panic("gzipstatic-gin: failed to get offset of gin.RouterGroup.engine")
	}
	const errmsg = "gzipstatic-gin: failed to check type of gin.RouterGroup.engine"
	if sf.Type.Kind() != reflect.Pointer {
		panic(errmsg)
	}
	b := reflect.TypeFor[gin.Engine]()
	sfte := sf.Type.Elem()
	if sfte.Kind() != b.Kind() {
		panic(errmsg)
	}
	if sfte.Size() != b.Size() {
		panic(errmsg)
	}
	if sfte.PkgPath() != b.PkgPath() {
		panic(errmsg)
	}
	if sfte.Name() != b.Name() {
		panic(errmsg)
	}
	return sf.Offset
}()

var offset_Engine_noRoute = func() uintptr {
	sf, ok := reflect.TypeFor[gin.Engine]().FieldByName("noRoute")
	if !ok {
		panic("gzipstatic-gin: failed to get offset of gin.Engine.noRoute")
	}
	const errmsg = "gzipstatic-gin: failed to check type of gin.Engine.noRoute"
	b := reflect.TypeFor[gin.HandlersChain]()
	if sf.Type.Kind() != b.Kind() {
		panic(errmsg)
	}
	if sf.Type.Size() != b.Size() {
		panic(errmsg)
	}
	if sf.Type.PkgPath() != b.PkgPath() {
		panic(errmsg)
	}
	if sf.Type.Name() != b.Name() {
		panic(errmsg)
	}
	return sf.Offset
}()

var offset_Context_handlers, offset_Context_index = func() (offset_handlers uintptr, offset_index uintptr) {
	t := reflect.TypeFor[gin.Context]()

	{
		sf, ok := t.FieldByName("handlers")
		if !ok {
			panic("gzipstatic-gin: failed to get offset of gin.Context.handlers")
		}
		const errmsg = "gzipstatic-gin: failed to check type of gin.Context.handlers"
		b := reflect.TypeFor[gin.HandlersChain]()
		if sf.Type.Kind() != b.Kind() {
			panic(errmsg)
		}
		if sf.Type.Size() != b.Size() {
			panic(errmsg)
		}
		if sf.Type.PkgPath() != b.PkgPath() {
			panic(errmsg)
		}
		if sf.Type.Name() != b.Name() {
			panic(errmsg)
		}
		offset_handlers = sf.Offset
	}

	{
		sf, ok := t.FieldByName("index")
		if !ok {
			panic("gzipstatic-gin: failed to get offset of gin.Context.index")
		}
		if sf.Type.Kind() != reflect.Int8 {
			panic("gzipstatic-gin: failed to check type of gin.Context.index")
		}
		offset_index = sf.Offset
	}

	return
}()

func getEngineUnsafePointer(group gin.IRoutes) unsafe.Pointer {
	if group == nil {
		return nil
	}
	if engine, ok := group.(*gin.Engine); ok {
		return unsafe.Pointer(engine)
	}
	if rg, ok := group.(*gin.RouterGroup); ok {
		return *(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(rg), offset_RouterGroup_engine))
	}
	return nil
}

// returns ok?
func serveFile(ctx *gin.Context, filePath string, fs http.FileSystem) bool {
	headerAcceptEncoding := ctx.GetHeader("Accept-Encoding")
	if headerAcceptEncoding == "" {
		return false
	}

	// if strings.HasSuffix(ctx.Request.URL.Path, "/index.html") {
	// 	ctx.Header("Location", "./")
	// 	ctx.Status(301)
	// 	return true
	// }

	if filePath == "" || strings.HasSuffix(filePath, "/") {
		filePath += "index.html"
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	if _, ok := ExtFilterMap[ext]; !ok {
		return false
	}

	acceptEncoding := make(map[string]string, len(EncodeNameExtMap))

	for encoding := range strings.SplitSeq(headerAcceptEncoding, ",") {
		if i := strings.IndexByte(encoding, ';'); i != -1 {
			encoding = encoding[:i]
		}
		encoding = strings.TrimSpace(encoding)
		if encoding == "" {
			continue
		}
		if ext, ok := EncodeNameExtMap[encoding]; ok {
			acceptEncoding[encoding] = ext
		}
	}

	for _, encoding := range EncodeNamePriority {
		ext, ok := acceptEncoding[encoding]
		if !ok {
			continue
		}

		file, err := fs.Open(filePath + ext)
		if err != nil {
			continue
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil || stat.IsDir() {
			continue
		}

		wh := ctx.Writer.Header()
		if EnableDebugHeader {
			wh.Add("X-Encoding-By", "gzipstatic-gin")
		}
		wh.Set("Content-Encoding", encoding)
		wh.Add("Vary", "Accept-Encoding")
		http.ServeContent(ctx.Writer, ctx.Request, filePath, stat.ModTime(), file)
		return true
	}
	return false
}

func File(ctx *gin.Context, FilePath string) {
	dir, name := filepath.Split(FilePath)
	if !serveFile(ctx, name, http.Dir(dir)) {
		ctx.File(FilePath)
	}
}

func FileFromFS(ctx *gin.Context, name string, fs http.FileSystem) {
	if !serveFile(ctx, name, fs) {
		ctx.FileFromFS(name, fs)
	}
}

func staticFileHandler(group gin.IRoutes, relativePath string, handler gin.HandlerFunc) gin.IRoutes {
	if strings.IndexByte(relativePath, ':') != -1 || strings.IndexByte(relativePath, '*') != -1 {
		panic("URL parameters can not be used when serving a static file")
	}
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
	return group
}

func StaticFile(group gin.IRoutes, relativePath, filepath string) gin.IRoutes {
	return staticFileHandler(group, relativePath, func(ctx *gin.Context) {
		File(ctx, filepath)
	})
}
func StaticFileFS(group gin.IRoutes, relativePath, filepath string, fs http.FileSystem) gin.IRoutes {
	return staticFileHandler(group, relativePath, func(ctx *gin.Context) {
		FileFromFS(ctx, filepath, fs)
	})
}

func Static(group gin.IRoutes, relativePath, root string) gin.IRoutes {
	return StaticFS(group, relativePath, gin.Dir(root, false))
}

func StaticFS(group gin.IRoutes, relativePath string, fs http.FileSystem) gin.IRoutes {
	if strings.IndexByte(relativePath, ':') != -1 || strings.IndexByte(relativePath, '*') != -1 {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := func(ctx *gin.Context) {
		name := ctx.Param("filepath")
		if serveFile(ctx, name, fs) {
			return
		}
		f, err := fs.Open(name)
		if err != nil {
			// 404 Not Found
			ctx.Status(http.StatusNotFound)
			// gin.Engine
			engineUP := getEngineUnsafePointer(group)
			if engineUP == nil {
				return
			}
			// gin.Engine.noRoute
			noRoute := (*gin.HandlersChain)(unsafe.Add(engineUP, offset_Engine_noRoute))
			// gin.Context
			ctxUP := unsafe.Pointer(ctx)
			// gin.Context.handlers
			*(*gin.HandlersChain)(unsafe.Add(ctxUP, offset_Context_handlers)) = *noRoute
			// gin.Context.index
			*(*int8)(unsafe.Add(ctxUP, offset_Context_index)) = -1
			return
		}
		f.Close()

		ctx.FileFromFS(name, fs)
	}
	urlPattern := path.Join(relativePath, "/*filepath")

	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
	return group
}
