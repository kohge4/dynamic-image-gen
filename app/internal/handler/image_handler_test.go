package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"

	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/config"
	mock_browser "github.com/kohge4/dynamic-img-gen-cdk/app/internal/test/mock/browser"
)

func TestGetImageByWebURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBrowserDriver := mock_browser.NewMockBrowserDriver(ctrl)

	handler := imageHandler{
		browserDriver: mockBrowserDriver,
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/v1/image/web-card", handler.GetImageByWebURL)

	t.Run("正常系", func(t *testing.T) {
		mockBrowserDriver.EXPECT().ScreenShot("https://screenshot-target-web-site.jp", "div.target", 1500, 800).Return([]byte("mocked image bytes"), nil)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/image/web-card?url=https://screenshot-target-web-site.jp&selector=div.target", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("異常系(不足パラメータあり)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/image/web-card?selector=div.target", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetImageByTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBrowserDriver := mock_browser.NewMockBrowserDriver(ctrl)

	handler := imageHandler{
		browserDriver: mockBrowserDriver,
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/v1/image/card", handler.GetImageByTemplate)

	t.Run("正常系", func(t *testing.T) {
		templateConfig := config.NewDefaultTemplateConfig("タイトル", "メッセージです")
		mockBrowserDriver.EXPECT().ScreenShot(
			templateConfig.InnerHTMLURL, templateConfig.Selector, templateConfig.Width, templateConfig.Height,
		).Return([]byte("mocked image bytes"), nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/image/card?title=タイトル&message=メッセージです", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("異常系(不足パラメータあり)", func(t *testing.T) {
		templateConfig := config.NewDefaultTemplateConfig("", "")
		mockBrowserDriver.EXPECT().ScreenShot(
			templateConfig.InnerHTMLURL, templateConfig.Selector, templateConfig.Width, templateConfig.Height,
		).Return([]byte("mocked image bytes"), nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/image/card", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
