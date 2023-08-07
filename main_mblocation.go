package main

import (
	// L "./lib"
	"github.com/gin-gonic/gin"
	L "github.com/mncbank/api/lib"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(L.JSONLogMiddleware())
	v := r.Group("/mbank/v1/location/")
	{
		v.POST("inquiry/", L.MbLocInq)
	}
	r.Run(":80") // listen and serve on 0.0.0.0:80
}
