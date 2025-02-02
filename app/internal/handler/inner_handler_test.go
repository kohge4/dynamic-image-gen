package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestGetCardHTM(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.MkdirAll("static/templates", os.ModePerm)
	os.WriteFile("static/templates/default_card.html", []byte("<html>Mock Template</html>"), os.ModePerm)
	defer os.RemoveAll("static")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("static/templates/*.html")
	router.GET("/internal/card", GetCardHTML)

	t.Run("正常系", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/internal/card?title=タイトルです&message=メッセージです", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("異常系(不足パラメータあり)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/internal/card", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
