package main

import (
	"fmt"
	"net/http"

	"github.com/bddjr/gzipstatic-gin"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache")
	})

	gzipstatic.Static(router, "/", "vite-project/dist")

	router.NoRoute(func(ctx *gin.Context) {
		ctx.Writer.WriteString("404")
	})

	fmt.Println("http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	panic(err)
}
