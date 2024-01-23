package md

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	rhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/mermaid"
	"mvdan.cc/xurls/v2"
)

type CustomHeadingRenderer struct{}

func (r *CustomHeadingRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHeading, r.renderHeading)
}

func (r *CustomHeadingRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		heading := node.(*ast.Heading)
		id := node.Attributes()[0].Value
		fmt.Fprintf(w, "<h%d id=\"%s\"><a href=\"#%s\" class=\"header-anchor\">%s", heading.Level, id, id, node.Text(source))
		fmt.Fprintf(w, "</a></h%d>", heading.Level)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func ToHTML(source []byte) []byte {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
					extension.LeftSingleQuote:  []byte("&sbquo;"),
					extension.RightSingleQuote: nil,
				}),
			),
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true),
				),
			),
			extension.NewLinkify(
				extension.WithLinkifyAllowedProtocols([][]byte{
					[]byte("http:"),
					[]byte("https:"),
				}),
				extension.WithLinkifyURLRegexp(
					xurls.Strict(),
				),
			),
			&mermaid.Extender{},
			&hashtag.Extender{},
			emoji.Emoji,
			mathjax.MathJax,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(&CustomHeadingRenderer{}, 500),
			),
			rhtml.WithHardWraps(),
			rhtml.WithXHTML(),
			rhtml.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
