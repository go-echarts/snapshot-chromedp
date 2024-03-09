package render

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/chromedp"
)

const (
	HTML               = "html"
	FileProtocol       = "file://"
	EchartsInstanceDom = "div[_echarts_instance_]"
	CanvasJs           = "echarts.getInstanceByDom(document.querySelector('div[_echarts_instance_]'))" +
		".getDataURL({type: '%s', pixelRatio: %d, excludeComponents: ['toolbox']})"
)

type SnapshotConfig struct {
	// Renderer canvas or svg, not used for now
	Renderer string
	// RenderContent the content bytes of charts after rendered
	RenderContent []byte
	// Path the path to save image
	Path string
	// FileName image name
	FileName string
	// Suffix image format, png, jpeg
	Suffix string
	// Quality the generated image quality, aka pixelRatio
	Quality int
	// KeepHtml whether keep the generated html also, default false
	KeepHtml bool
	// HtmlPath where to keep the generated html, default same to image path
	HtmlPath string
}

func MakeChartSnapshot(content []byte, image string) error {
	path, file := filepath.Split(image)
	suffix := filepath.Ext(file)[1:]
	fileName := file[0 : len(file)-len(suffix)-1]

	config := &SnapshotConfig{
		RenderContent: content,
		Path:          path,
		FileName:      fileName,
		Suffix:        suffix,
		Quality:       1,
		KeepHtml:      false,
	}
	return MakeSnapshot(config)
}

func MakeSnapshot(config *SnapshotConfig) error {
	path := config.Path
	fileName := config.FileName
	content := config.RenderContent
	quality := config.Quality
	suffix := config.Suffix
	keepHtml := config.KeepHtml
	htmlPath := config.HtmlPath

	if htmlPath == "" {
		htmlPath = path
	}

	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	if !filepath.IsAbs(htmlPath) {
		htmlPath, _ = filepath.Abs(htmlPath)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	htmlFullPath := filepath.Join(htmlPath, fileName+"."+HTML)

	if !keepHtml {
		defer func() {
			err := os.Remove(htmlFullPath)
			if err != nil {
				log.Printf("Failed to delete the file(%s), err: %s\n", htmlFullPath, err)
			}
		}()
	}

	err := os.WriteFile(htmlFullPath, content, 0o644)
	if err != nil {
		return err
	}

	if quality < 1 {
		quality = 1
	}

	var base64Data string
	executeJS := fmt.Sprintf(CanvasJs, suffix, quality)
	err = chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("%s%s", FileProtocol, htmlFullPath)),
		chromedp.WaitVisible(EchartsInstanceDom, chromedp.ByQuery),
		chromedp.Evaluate(executeJS, &base64Data),
	)
	if err != nil {
		return err
	}

	imgContent, err := base64.StdEncoding.DecodeString(strings.Split(base64Data, ",")[1])
	if err != nil {
		return err
	}

	imageFullPath := filepath.Join(path, fmt.Sprintf("%s.%s", fileName, suffix))
	if err := os.WriteFile(imageFullPath, imgContent, 0o644); err != nil {
		return err
	}

	log.Printf("Wrote %s.%s success", fileName, suffix)
	return nil
}
