package excelkit

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type Param struct {
	Title string
	Sheet string
	// export时，key:name:width
	Keys []string
	Data []map[string]interface{}
	// load时, name:key，选填，自动生成
	//Names []string
	File class.File
	// 导出文件名
	FileName string
}

type KeyDef struct {
	Name  string
	Index int
	Width float64
}

// Export title/_sheet:x/key:name/key:name:width
// todo 其中的err仍未处理
func Export(param Param, ctx *context.Context) {
	if len(param.Keys) == 0 {
		panic(exception.New("excel keys empty"))
	}
	f := excelize.NewFile()
	if param.Sheet == "" {
		param.Sheet = f.GetSheetName(0)
	} else if f.GetSheetName(0) != param.Sheet {
		f.SetSheetName(f.GetSheetName(0), param.Sheet)
	}
	keyMap := map[string]KeyDef{}
	for i, key := range param.Keys {
		ts := stringkit.Split(key, ":")
		if len(ts) < 2 {
			panic(exception.New("excel keys param 语法错误"))
		}
		m := KeyDef{
			Name:  ts[1],
			Index: i,
		}
		if len(ts) > 2 {
			m.Width = cast.ToFloat64(ts[2])
		}
		keyMap[ts[0]] = m
	}
	// style title
	titleStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Size: 15,
		},
	})
	if err != nil {
		panic(exception.New(err.Error()))
	}
	cellStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			//Horizontal: "right",
			Vertical: "center",
			WrapText: true,
		},
		Font: &excelize.Font{
			Size: 12,
		},
		Border: BorderStyleDefault(),
	})
	if err != nil {
		panic(exception.New(err.Error()))
	}
	// title
	err = f.MergeCell(param.Sheet, "A1", string(rune('A'+len(param.Keys)-1))+"1")
	err = f.SetCellStyle(param.Sheet, "A1", string(rune('A'+len(param.Keys)-1))+"1", titleStyle)
	err = f.SetCellValue(param.Sheet, "A1", param.Title)
	// key title
	for _, v := range keyMap {
		cell := string(rune('A'+v.Index)) + "2"
		err = f.SetCellStyle(param.Sheet, cell, cell, cellStyle)
		err = f.SetCellValue(param.Sheet, cell, v.Name)
		if v.Width > 0 {
			err = f.SetColWidth(param.Sheet, string(rune('A'+v.Index)), string(rune('A'+v.Index)), v.Width)
		}
	}
	// data
	for i, data := range param.Data {
		index := i + 3
		// 每个cell加style
		for j := range param.Keys {
			cell := string(rune('A'+j)) + cast.ToString(index)
			err = f.SetCellStyle(param.Sheet, cell, cell, cellStyle)
		}
		for k, v := range data {
			cell := string(rune('A'+keyMap[k].Index)) + cast.ToString(index)
			err = f.SetCellValue(param.Sheet, cell, v)
		}
	}
	// 发送至web stream
	if param.FileName == "" {
		param.FileName = "export.xlsx"
	}
	//err = f.SaveAs("/Users/ycj/Downloads/test.xlsx")
	ctx.SetFileHeader(param.FileName)
	err = f.Write(ctx.Proxy.ResponseWriter())
	if err != nil {
		panic(exception.New("excel export error: " + err.Error()))
	}
}

// Load name(题头):key(map-key):type(number)
func Load(param Param) []map[string]string {
	if len(param.Keys) == 0 {
		panic(exception.New("excel names empty"))
	}
	nameMap := map[string]string{}
	for _, key := range param.Keys {
		ts := stringkit.Split(key, ":")
		if len(ts) < 2 {
			panic(exception.New("excel keys param 语法错误"))
		}
		nameMap[ts[1]] = ts[0]
	}
	var f *excelize.File
	var err error
	if param.File.File != nil {
		f, err = excelize.OpenReader(param.File.File)
	} else {
		panic(exception.New("file is nil"))
	}
	if err != nil {
		panic(exception.New(err.Error()))
	}
	rows, err := f.Rows(f.GetSheetName(0))
	if err != nil {
		panic(exception.New(err.Error()))
	}
	var res []map[string]string
	index := 1
	var names []string
	for rows.Next() {
		if index == 2 {
			names, _ = rows.Columns()
		} else if index > 2 {
			m := map[string]string{}
			values, _ := rows.Columns()
			for i, v := range values {
				if names != nil {
					m[nameMap[names[i]]] = v
				}
			}
			res = append(res, m)
		} else {
			// 需要取出，否则后续names重叠
			_, _ = rows.Columns()
		}
		index++
	}
	return res
}

func BorderStyleDefault() []excelize.Border {
	return []excelize.Border{
		{
			Type:  "left",
			Color: "000000",
			Style: 1,
		},
		{
			Type:  "right",
			Color: "000000",
			Style: 1,
		},
		{
			Type:  "top",
			Color: "000000",
			Style: 1,
		},
		{
			Type:  "bottom",
			Color: "000000",
			Style: 1,
		},
	}
}
