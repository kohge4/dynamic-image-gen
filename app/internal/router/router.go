package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/browser"
	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/handler"
)

func New() *gin.Engine {
	r := gin.Default()

	imageHandler := handler.NewImageHandler(browser.NewBrowserDriver())

	v1 := r.Group("/v1")
	{
		v1.GET("/image/web-card", imageHandler.GetImageByWebURL)
		v1.GET("/image/card", imageHandler.GetImageByTemplate)
	}

	return r
}
