package markdown

import (
	"bytes"
	"fmt"
	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/russross/blackfriday/v2"
)

func Html() {
	content, _ := filekit.ReadString("./README.md")

	//render:=blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
	//	Flags: blackfriday.CommonHTMLFlags|blackfriday.TOC,
	//})
	render := bfchroma.NewRenderer(
		//bfchroma.WithoutAutodetect(),
		bfchroma.ChromaOptions(
			html.WithLineNumbers(true),
			html.WithClasses(true),
		),
		bfchroma.Extend(
			blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				Flags: blackfriday.CommonHTMLFlags | blackfriday.TOC,
			}),
		),
	)
	extensions := blackfriday.WithExtensions(blackfriday.CommonExtensions)
	data := blackfriday.Run([]byte(content), blackfriday.WithRenderer(render), extensions)
	// html
	// code css
	styles := make([]byte, 0)
	buffer := bytes.NewBuffer(styles)
	_ = render.ChromaCSS(buffer)
	// custom css
	// todo 加入css
	css := ""

	fin := fmt.Sprintf(`<html>
	<head>
		<style>%s</style>
		<style>%s</style>
	</head>
	<body>%s</body></html>`, buffer.String(), css, string(data))
	_ = filekit.WriteFile("/Users/ycj/Downloads/test.html", []byte(fin))
}
