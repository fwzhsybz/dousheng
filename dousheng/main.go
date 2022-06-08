package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() //初始化路由引擎

	initRouter(r) //调用多路由模块

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
