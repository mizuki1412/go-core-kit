package markdown

import (
	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/russross/blackfriday/v2"
)

func Test() {
	content, _ := filekit.ReadString("./README.md")

	// todo 加入head+css
	//render:=blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
	//	Flags: blackfriday.CommonHTMLFlags|blackfriday.TOC,
	//})
	render := bfchroma.NewRenderer(
		bfchroma.WithoutAutodetect(),
		bfchroma.ChromaOptions(
			html.WithLineNumbers(true),
		),
		bfchroma.Extend(
			blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				Flags: blackfriday.CommonHTMLFlags | blackfriday.TOC,
			}),
		),
	)
	extensions := blackfriday.WithExtensions(blackfriday.CommonExtensions)

	_ = filekit.WriteFile("/Users/ycj/Downloads/test.html", blackfriday.Run([]byte(content), blackfriday.WithRenderer(render), extensions))
}
