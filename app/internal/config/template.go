package config

import "fmt"

const (
	// スクリーンショット対象のHTML用のテンプレートファイル
	DefaultCardTemplateFileName = "default_card.html"
	DefaultCardTemplateFilePath = "static/templates/" + DefaultCardTemplateFileName
)

const (
	// TwitterのOGP画像用のサイズ
	TwitterCardWidth  = 1500
	TwitterCardHeight = 800
)

const (
	DefaultCardScreenshotTargetSelector = "div.screenshot-target"
	DefaultCardInnerHTMLURL             = "http://localhost:8081/internal/card?title=%s&message=%s"
)

// スクリーンショット対象のHTMLを内部で用意する際の設定を管理する構造体
type TemplateConfig struct {
	Width        int
	Height       int
	Selector     string
	InnerHTMLURL string
}

func NewDefaultTemplateConfig(title, message string) *TemplateConfig {
	url := fmt.Sprintf(DefaultCardInnerHTMLURL, title, message)
	return &TemplateConfig{
		Width:        TwitterCardWidth,
		Height:       TwitterCardHeight,
		Selector:     DefaultCardScreenshotTargetSelector,
		InnerHTMLURL: url,
	}
}

func NewTemplateConfig(width, height int, selector, innerHTMLURL string) *TemplateConfig {
	return &TemplateConfig{
		Width:        width,
		Height:       height,
		Selector:     selector,
		InnerHTMLURL: innerHTMLURL,
	}
}
