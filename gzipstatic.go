// https://github.com/bddjr/gzipstatic-gin
package gzipstatic

import (
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

// X-Encoding-By: gzipstatic-gin
var EnableDebugHeader = true

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
}

var EncodeNameExtMap = map[string]string{
	"br":   ".br",
	"zstd": ".zst",
	"gzip": ".gz",
}

func getEngineElem(group gin.IRoutes) (reflect.Value, bool) {
	if group != nil {
		if engine, ok := group.(*gin.Engine); ok {
			return reflect.ValueOf(engine).Elem(), true
		}
		if rg, ok := group.(*gin.RouterGroup); ok {
			return reflect.ValueOf(rg).Elem().FieldByName("engine").Elem(), true
		}
	}
	return reflect.Value{}, false
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

	var file http.File
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var minSize int64 = -1
	var minSizeFileModTime time.Time
	minSizeFileEncodeName := ""

	if f, err := fs.Open(filePath); err == nil {
		s, err := f.Stat()
		if err != nil || s.IsDir() {
			f.Close()
		} else {
			file = f
			minSize = s.Size()
			minSizeFileModTime = s.ModTime()
		}
	}

	for encodeName := range strings.SplitSeq(headerAcceptEncoding, ",") {
		if i := strings.IndexByte(encodeName, ';'); i != -1 {
			encodeName = encodeName[:i]
		}
		encodeName = strings.TrimSpace(encodeName)
		if encodeName == "" {
			continue
		}
		encodeExt, ok := EncodeNameExtMap[encodeName]
		if !ok || encodeExt == "" {
			continue
		}
		if encodeExt[0] != '.' {
			encodeExt = "." + encodeExt
		}

		f, err := fs.Open(filePath + encodeExt)
		if err != nil {
			continue
		}

		s, err := f.Stat()
		if err != nil || s.IsDir() {
			f.Close()
		} else if minSize == -1 || s.Size() < minSize {
			if file != nil {
				file.Close()
			}
			file = f
			minSize = s.Size()
			minSizeFileModTime = s.ModTime()
			minSizeFileEncodeName = encodeName
		} else {
			f.Close()
		}
	}

	if minSize == -1 {
		// Not Found
		return false
	}

	wh := ctx.Writer.Header()
	if EnableDebugHeader {
		wh.Add("X-Encoding-By", "gzipstatic-gin")
	}
	if minSizeFileEncodeName != "" {
		wh.Set("Content-Encoding", minSizeFileEncodeName)
		wh.Add("Vary", "Accept-Encoding")
	}

	http.ServeContent(ctx.Writer, ctx.Request, filePath, minSizeFileModTime, file)
	return true
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
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
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
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
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
			engineElem, ok := getEngineElem(group)
			if !ok {
				return
			}
			noRoute := (*gin.HandlersChain)(unsafe.Pointer(engineElem.FieldByName("noRoute").UnsafeAddr()))
			ctxElem := reflect.ValueOf(ctx).Elem()
			*(*gin.HandlersChain)(unsafe.Pointer(ctxElem.FieldByName("handlers").UnsafeAddr())) = *noRoute
			*(*int8)(unsafe.Pointer(ctxElem.FieldByName("index").UnsafeAddr())) = -1
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
