package template

import (
	_ "embed"
	"io/fs"
)

//go:embed *.html.tmpl
var FS fs.FS
