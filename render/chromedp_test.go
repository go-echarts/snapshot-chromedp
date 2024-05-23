package render

import (
	"context"
	_ "embed"
	"os"
	"testing"
	"time"

	"github.com/go-echarts/snapshot-chromedp/asset"
)

func TestFileCreation(t *testing.T) {
	fileName := "mock"
	fileImage := fileName + ".png"
	fileHtml := fileName + ".html"

	err := MakeChartSnapshot(asset.RenderContent(), fileImage)
	if err != nil {
		t.Fatalf("Failed to create file: %s", err)
	}

	_, err = os.Stat(fileImage)
	if os.IsNotExist(err) {
		t.Fatalf("Image File was exist")
	}

	_, err = os.Stat(fileHtml)
	if os.IsExist(err) {
		t.Fatalf("Html File was not exist")
	}
}

func TestFileCreationWithFullConfiguration(t *testing.T) {
	fileName := "keepBothResources"
	fileImage := fileName + ".jpeg"
	fileHtml := fileName + ".html"

	config := NewSnapshotConfig(asset.RenderContent(), fileImage)
	config.KeepHtml = true

	err := MakeSnapshot(config)
	if err != nil {
		t.Fatalf("Failed to create file: %s", err)
	}

	_, err = os.Stat(fileImage)
	if os.IsNotExist(err) {
		t.Fatalf("Image File was not exist")
	}

	_, err = os.Stat(fileHtml)
	if os.IsNotExist(err) {
		t.Fatalf("Html File was not exist")
	}
}

func TestFileCreationWithTimeout(t *testing.T) {
	fileName := "timeoutResources.jpeg"

	err := MakeSnapshot(NewSnapshotConfig(asset.RenderContent(), fileName, func(config *SnapshotConfig) {
		config.Timeout = time.Nanosecond
	}))

	if err != context.DeadlineExceeded {
		t.Fatalf("Can not cancel for creating file with timeout config: %s", fileName)
	}
}
