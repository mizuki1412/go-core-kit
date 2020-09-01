package markdown

import (
	"bytes"
	"fmt"
	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/russross/blackfriday/v2"
	"os"
	"path/filepath"
	"strings"
)

func Run(title string) {
	genHtml(title, integrate())
}

type mdData struct {
	Title   string
	Content string
}

// 整合所有的readme
func integrate() string {
	content, _ := filekit.ReadString("./README.md")
	// 提取第一层
	// [ { title:xxx, content:xxx} ]
	list := make([]*mdData, 0, 5)
	// title:mdData
	cache := map[string]*mdData{}
	// 解析
	lines := strings.Split(content, "\n")
	// 代码块标识
	flagCode := false
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.Index(l, "```") > -1 {
			flagCode = !flagCode
		}
		if !flagCode && len(l) > 1 && l[0] == '#' && l[1] != '#' {
			cache[l] = &mdData{
				Title:   l,
				Content: "",
			}
			list = append(list, cache[l])
		} else if len(list) > 0 {
			list[len(list)-1].Content += l + "\n"
		}
	}
	// 遍历目录下面
	var files []string
	_ = filepath.Walk("./", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if path[0] != '.' && strings.LastIndex(path, "README.md") > -1 && path != "README.md" {
			files = append(files, path)
		}
		return nil
	})
	for _, path := range files {
		content2, _ := filekit.ReadString(path)
		lines2 := strings.Split(content2, "\n")
		var title string
		for _, l := range lines2 {
			l = strings.TrimSpace(l)
			if len(l) > 1 && l[0] == '#' && l[1] != '#' {
				title = l
				if data, ok := cache[title]; ok {
					data.Content += "\n"
				}
			} else if data, ok := cache[title]; ok {
				data.Content += l + "\n"
			}
		}
	}
	var ret string
	for _, e := range list {
		ret += e.Title + "\n\n" + e.Content + "\n"
	}
	return ret
}

func genHtml(title, content string) {
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
	css := `
.title{
	font-size: 1.3rem;
	text-align: center;
	padding: 10px 0 20px 0;
}
.content{
	padding: 0 1rem;
}
p{
	padding-left: 1rem;
}
h2{
	color: midnightblue;
}
h3{
	color: brown;
}
`

	fin := fmt.Sprintf(`
<html>
	<head>
		<style>%s</style>
		<style>%s</style>
	</head>
	<body>
		<div class="title">%s</div>
		<div class="content">%s</div>
	</body>
</html>`, buffer.String(), css, title, string(data))
	_ = filekit.WriteFile("/Users/ycj/Downloads/doc.html", []byte(fin))
}
