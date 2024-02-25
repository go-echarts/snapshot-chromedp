package render

import (
	_ "embed"
	"github.com/go-echarts/snapshot-chromedp/asset"
	"os"
	"testing"
)

// TODO: integration in ci
func TestFileCreation(t *testing.T) {

	fileName := "mock"
	fileImage := fileName + ".png"
	fileHtml := fileName + ".html"

	err := MakeChartSnapshot(asset.RenderContent(), fileImage)

	if err != nil {
		t.Fatalf("Failed to create file: %s", err)
	}

	//defer os.Remove(fileImage)

	_, err = os.Stat(fileImage)
	if os.IsNotExist(err) {
		t.Fatalf("Image File was not exist")
	}

	_, err = os.Stat(fileHtml)
	if os.IsExist(err) {
		t.Fatalf("Html File was not exist")
	}
}
