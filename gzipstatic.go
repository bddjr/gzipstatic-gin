// https://github.com/bddjr/gzipstatic-gin
package gzipstatic

import (
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var ExtFillter = regexp.MustCompile(`\.(html|htm|js|json|css)$`)

type EncodeListItem struct {
	name string
	ext  string
}

// Priority from high to low
var EncodeList = []*EncodeListItem{
	{
		name: "br",
		ext:  ".br",
	}, {
		name: "gzip",
		ext:  ".gz",
	},
}

var NoRoute gin.HandlerFunc = nil

func tryCompress(ctx *gin.Context, name string, fs http.FileSystem) (next bool) {
	if strings.HasSuffix(ctx.Request.URL.Path, "/index.html") {
		ctx.Header("Location", "../")
		ctx.Status(301)
		return false
	}

	if name == "" || strings.HasSuffix(name, "/") {
		name += "index.html"
	}

	ext := ExtFillter.FindString(name)
	if ext == "" {
		return true
	}

	ae := ctx.GetHeader("Accept-Encoding")

	for _, encode := range EncodeList {
		if !strings.Contains(ae, encode.name) {
			continue
		}

		f, err := fs.Open(name + encode.ext)
		if err != nil {
			continue
		}
		defer f.Close()

		s, err := f.Stat()
		if err != nil || s.IsDir() {
			f.Close()
			continue
		}

		h := ctx.Writer.Header()
		h.Del("Content-Length")
		h.Add("Vary", "Accept-Encoding")
		h.Set("X-Content-Encoding-By", "gzipstatic-gin")
		h.Set("Content-Encoding", encode.name)
		h.Set("Content-Type", mime.TypeByExtension(ext))

		http.ServeContent(ctx.Writer, ctx.Request, name, s.ModTime(), f)
		return false
	}
	return true
}

func File(ctx *gin.Context, FilePath string) {
	dir, name := filepath.Split(FilePath)
	if tryCompress(ctx, name, http.Dir(dir)) {
		ctx.File(FilePath)
	}
}

func FileFromFS(ctx *gin.Context, name string, fs http.FileSystem) {
	if tryCompress(ctx, name, fs) {
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
		if !tryCompress(ctx, name, fs) {
			return
		}
		if NoRoute != nil {
			f, err := fs.Open(name)
			if err != nil {
				ctx.Status(404)
				NoRoute(ctx)
				return
			}
			f.Close()
		}
		ctx.FileFromFS(name, fs)
	}
	urlPattern := path.Join(relativePath, "/*filepath")

	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
	return group
}
