package template

import "io/fs"

//go:embed *.html.tmpl
var FS fs.FS
