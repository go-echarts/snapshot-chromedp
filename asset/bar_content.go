package asset

import (
	_ "embed"
)

//go:embed bar.html
var barHTML []byte

//go:embed page.html
var pageHTML []byte

func RenderContent() []byte {
	return barHTML
}

func RenderPageContent() []byte {
	return pageHTML
}
