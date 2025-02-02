package handler

import "github.com/gin-gonic/gin"

type GetCardHTMLRequest struct {
	Title   string `form:"title"`
	Message string `form:"message"`
}

func (r *GetCardHTMLRequest) Bind(c *gin.Context) error {
	return c.ShouldBind(r)
}

func newGetCardHTMLRequest() *GetCardHTMLRequest {
	return &GetCardHTMLRequest{}
}
