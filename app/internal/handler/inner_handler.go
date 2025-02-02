package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/config"
)

func GetCardHTML(c *gin.Context) {
	req := newGetCardHTMLRequest()
	if err := req.Bind(c); err != nil {
		c.HTML(http.StatusBadRequest, config.DefaultCardTemplateFileName, gin.H{"title": "400 Validation Error(title and message is required)"})
		return
	}

	if _, err := os.Stat(config.DefaultCardTemplateFilePath); err != nil {
		c.HTML(http.StatusNotFound, config.DefaultCardTemplateFileName, gin.H{"title": "404 Not Found"})
		return
	}

	c.HTML(http.StatusOK, config.DefaultCardTemplateFileName, gin.H{
		"title":   req.Title,
		"message": req.Message,
	})
}
