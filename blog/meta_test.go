package blog

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

func TestMeta(t *testing.T) {
	markdown := goldmark.New(goldmark.WithExtensions, Meta)
	source := `# Example post
## This is a subtitle
# en, test,random
# 2022-08-21

This is the body
`
	var buf bytes.Buffer
	ctx := parser.NewContext()
	err := markdown.Convert([]byte(source), &buf, parser.WithContext(ctx))
	if err != nil {
		t.Fatalf("Could not parse test data", err)
	}
	metadata := Get(context)
	// Assert title
	// Assert tags
	// Assert date
}
