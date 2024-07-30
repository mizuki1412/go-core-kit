package excelkit

import (
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/stringkit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"strings"
)

// Param 注意，excel的格式必须是第一行title，后面k-v结构
type Param struct {
	Title string
	Sheet string
	// export时，key:name:width; load时，key:name
	Keys []string
	Data []map[string]any
	File class.File
	// 文件来源：本地文件路径
	FileOriginPath string
	// 导出的文件名
	FileName string
	// 导出存入本地文件地址
	SaveToFile string
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
	//err = f.MergeCell(param.Sheet, "A1", string(rune('A'+len(param.Keys)-1))+"1")
	//err = f.SetCellStyle(param.Sheet, "A1", string(rune('A'+len(param.Keys)-1))+"1", titleStyle)
	err = f.MergeCell(param.Sheet, "A1", baseConversion(len(param.Keys)-1)+"1")
	err = f.SetCellStyle(param.Sheet, "A1", baseConversion(len(param.Keys)-1)+"1", titleStyle)
	err = f.SetCellValue(param.Sheet, "A1", param.Title)
	// key title
	for _, v := range keyMap {
		//cell := string(rune('A'+v.Index)) + "2"
		cell := baseConversion(v.Index) + "2"
		err = f.SetCellStyle(param.Sheet, cell, cell, cellStyle)
		err = f.SetCellValue(param.Sheet, cell, v.Name)
		if v.Width > 0 {
			//err = f.SetColWidth(param.Sheet, string(rune('A'+v.Index)), string(rune('A'+v.Index)), v.Width)
			err = f.SetColWidth(param.Sheet, baseConversion(v.Index), baseConversion(v.Index), v.Width)
		}
	}
	// data
	for i, data := range param.Data {
		index := i + 3
		// 每个cell加style
		for j := range param.Keys {
			//cell := string(rune('A'+j)) + cast.ToString(index)
			cell := baseConversion(j) + cast.ToString(index)
			err = f.SetCellStyle(param.Sheet, cell, cell, cellStyle)
		}
		for k, v := range data {
			if _, ok := keyMap[k]; !ok {
				continue
			}
			//cell := string(rune('A'+keyMap[k].Index)) + cast.ToString(index)
			cell := baseConversion(keyMap[k].Index) + cast.ToString(index)
			err = f.SetCellValue(param.Sheet, cell, v)
		}
	}
	// 发送至web stream
	if param.FileName == "" {
		param.FileName = "export.xlsx"
	}

	if param.SaveToFile != "" {
		err = f.SaveAs(param.SaveToFile)
		if err != nil {
			panic(exception.New("excel export error: " + err.Error()))
		}
	} else {
		ctx.SetFileHeader(param.FileName)
		err = f.Write(ctx.Proxy.Writer)
		if err != nil {
			panic(exception.New("excel export error: " + err.Error()))
		}
	}
}

// Load key:name
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
	} else if param.FileOriginPath != "" {
		f, err = excelize.OpenFile(param.FileOriginPath)
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
			// excelize.Options{RawCellValue: true}
			values, _ := rows.Columns()
			for i, v := range values {
				if names != nil && len(names) > i && nameMap[names[i]] != "" {
					m[nameMap[names[i]]] = strings.TrimSpace(v)
				}
			}
			if len(m) > 0 {
				res = append(res, m)
			}
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

var base = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func baseConversion(number int) string {
	var ret string
	var m []rune
	length := len(base)
	flag := false
	for number >= 0 {
		if flag {
			break
		}
		index := number % length
		m = append(m, base[index])
		temp := number / length
		if temp == 0 {
			flag = true
		}
		number = temp - 1

	}
	for i := len(m) - 1; i >= 0; i-- {
		ret += string(m[i])
	}
	return ret
}
