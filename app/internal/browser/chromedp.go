package browser

import (
	"context"

	"github.com/chromedp/chromedp"
)

type BrowserDriver interface {
	ScreenShot(url, selector string, width, height int) ([]byte, error)
}

type browserDriver struct{}

func NewBrowserDriver() BrowserDriver {
	return &browserDriver{}
}

func (b *browserDriver) ScreenShot(url, selector string, width, height int) ([]byte, error) {
	// 参考: chromiumのオプションの説明一覧 https://peter.sh/experiments/chromium-command-line-switches/
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(width, height),
		chromedp.Flag("lang", "ja"),
		chromedp.NoSandbox, // no-zygoteを有効にする場合は必要
		// 参考: Lambda用 https://github.com/chromedp/chromedp/issues/1074#issuecomment-1188370109 (chromedp.DefaultExecAllocatorOptionsにないもののみ追加)
		chromedp.Flag("single-process", true),
		chromedp.Flag("no-zygote", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Screenshot(selector, &buf, chromedp.NodeVisible),
	}); err != nil {
		return nil, err
	}

	return buf, nil
}

func elementScreenshot(url, selector string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Screenshot(selector, res, chromedp.NodeVisible),
	}
}
