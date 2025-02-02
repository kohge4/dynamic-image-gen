package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/config"
)

type GetImageByWebURLRequest struct {
	URL      string `form:"url" binding:"required"`
	Selector string `form:"selector" binding:"required"`
	Width    int    `form:"width"`
	Height   int    `form:"height"`
}

func (r *GetImageByWebURLRequest) Bind(c *gin.Context) error {
	if err := c.ShouldBind(r); err != nil {
		return err
	}

	if r.Width == 0 {
		r.Width = config.TwitterCardWidth
	}

	if r.Height == 0 {
		r.Height = config.TwitterCardHeight
	}

	return nil
}

func newGetImageByWebURLRequest() *GetImageByWebURLRequest {
	return &GetImageByWebURLRequest{}
}

type GetImageByTemplateRequest struct {
	Title   string `form:"title"`
	Message string `form:"message"`
}

func (r *GetImageByTemplateRequest) Bind(c *gin.Context) error {
	if err := c.ShouldBind(r); err != nil {
		return err
	}
	return nil
}

func newGetImageByTemplateRequest() *GetImageByTemplateRequest {
	return &GetImageByTemplateRequest{}
}
