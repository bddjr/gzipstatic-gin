package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bddjr/gzipstatic-gin"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use(func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache")
	})

	gzipstatic.Static(engine, "/assets", "vite-project/dist/assets")
	gzipstatic.StaticFile(engine, "/", "vite-project/dist/index.html")
	gzipstatic.StaticFile(engine, "/vite.svg", "vite-project/dist/vite.svg")

	group := engine.Group("group")
	gzipstatic.Static(group, "/", "vite-project/dist")

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.Writer.WriteString("test404 " + time.Now().Format("2006-01-02 15:04:05 UTC-0700"))
	})

	fmt.Println("http://localhost:8080")
	fmt.Println("http://localhost:8080/group")
	err := http.ListenAndServe(":8080", engine)
	panic(err)
}
