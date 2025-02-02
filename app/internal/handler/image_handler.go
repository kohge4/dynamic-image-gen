package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/browser"
	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/config"
)

type ImageHandler interface {
	GetImageByWebURL(c *gin.Context)
	GetImageByTemplate(c *gin.Context)
}

type imageHandler struct {
	browserDriver browser.BrowserDriver
}

func NewImageHandler(browserDriver browser.BrowserDriver) ImageHandler {
	return &imageHandler{
		browserDriver: browserDriver,
	}
}

func (h *imageHandler) GetImageByWebURL(c *gin.Context) {
	req := newGetImageByWebURLRequest()
	if err := req.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageBytes, err := h.browserDriver.ScreenShot(req.URL, req.Selector, req.Width, req.Height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, http.DetectContentType(imageBytes), imageBytes)
}

func (h *imageHandler) GetImageByTemplate(c *gin.Context) {
	req := newGetImageByTemplateRequest()
	if err := req.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templateConfig := config.NewDefaultTemplateConfig(req.Title, req.Message)
	imageBytes, err := h.browserDriver.ScreenShot(templateConfig.InnerHTMLURL, templateConfig.Selector, templateConfig.Width, templateConfig.Height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, http.DetectContentType(imageBytes), imageBytes)
}
