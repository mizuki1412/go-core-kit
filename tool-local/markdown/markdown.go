package markdown

import (
	"bytes"
	"fmt"
	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/service/pdfkit"
	"github.com/russross/blackfriday/v2"
	"os"
	"path/filepath"
	"strings"
)

func Run(title, dest string) {
	genHtml(title, dest, integrate())
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

func genHtml(title, dest, content string) {
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
body{
	font-size: 24px;
}
.flex {
  display: -webkit-box; /* wkhtmltopdf uses this one */
  -webkit-box-orient: horizontal;
  /* display: flex; */
}
.inline {
  display: inline-flex;
}

.center {
  -webkit-box-pack: center; /* wkhtmltopdf uses this one */
  -webkit-box-align: center;
}

.around {
  -webkit-box-pack: justify;
}

.start {
  -webkit-box-pack: start;
}

.end {
  -webkit-box-pack: end;
}

.align-center {
  -webkit-box-pack: center;
}

.content-center {
  -webkit-box-pack: center;
}

.title{
	font-size: 1.3rem;
	text-align: center;
	padding: 10px 0 20px 0;
}
// 主体内容部分
.content{
	padding: 0 1rem;
}
.chroma{
	box-shadow: 0px 0px 5px black;
    border-radius: 8px;
    padding: 12px 0;
}
p{
	padding-left: 1rem;
}

h1{
  	counter-increment: h1counter;
  	counter-reset: h2counter;
}
h1:before {
	content: counter(h1counter) " ";
}
h2{
	color: midnightblue;
  	counter-increment: h2counter;
	counter-reset: h3counter;
}
h2:before {
  content: counter(h1counter) "." counter(h2counter) " ";
}
h3{
	color: brown;
	counter-increment: h3counter;
}
h3:before
{
  content: counter(h1counter) "." counter(h2counter) "." counter(h3counter) " ";
}
`
	fin := fmt.Sprintf(`
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no" />
		<style>%s</style>
		<style>%s</style>
	</head>
	<body>
		<div class="title">%s</div>
		<div class="content">%s</div>
	</body>
</html>`, buffer.String(), css, title, string(data))
	pdfkit.Gen2File(fin, dest+"/"+title+".pdf")
	//_ = filekit.WriteFile("/Users/ycj/Downloads/doc.html", []byte(fin))
}
