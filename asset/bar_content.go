package asset

import (
	_ "embed"
)

//go:embed bar.html
var barHTML []byte

func RenderContent() []byte {
	return barHTML
}
