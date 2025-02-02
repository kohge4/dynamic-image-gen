package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/handler"
)

func NewInnerRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("static/templates/*.html")

	internal := r.Group("/internal")
	{
		internal.Static("/css", "static/css")
		internal.GET("/card", handler.GetCardHTML)
	}
	return r
}
