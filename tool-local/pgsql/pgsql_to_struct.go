package pgsql

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"strings"
)

type Field struct {
	Name string
	Type string
	Tags []string
}

// todo 外联指针对象无法区分；json的omitempty标记不够精确（只能对class或指针）
func SQL2Struct(sqlFile, destFile string) {
	sqls, err := filekit.ReadString(sqlFile)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(sqls, "\n")
	var dest, temp string
	var fields []Field
	var table string
	for i := 0; i < len(parts); i++ {
		val := strings.TrimSpace(parts[i])
		if strings.Index(val, "-") == 0 || strings.Index(val, "create index") == 0 || strings.Index(val, "insert") == 0 || strings.Index(val, "update") == 0 {
			continue
		}
		if strings.Index(val, "create table") == 0 {
			table = val[strings.Index(val, "create")+13 : len(val)-1]
			temp = "type " + stringkit.CamelCase(table) + " struct{\n"
			fields = []Field{}
		} else if strings.Index(val, ")") == 0 {
			// end
			for _, f := range fields {
				temp += "\t" + stringkit.CamelCase(f.Name) + " " + f.Type + " `" + strings.Join(f.Tags, " ") + "`\n"
			}
			temp += "}\n\n"
			dest += temp
			temp = ""
		} else if temp != "" {
			es := stringkit.Split(val, "[ ,\t]+")
			if es[0] == "primary" {
				// todo 单独定义primary key 时
				continue
			}
			f := Field{Name: es[0]}
			switch es[1] {
			case "varchar", "text":
				f.Type = "class.String"
			case "serial":
				f.Type = "int32"
				f.Tags = append(f.Tags, "autoincrement:\"true\"")
			case "bigserial":
				f.Type = "int64"
				f.Tags = append(f.Tags, "autoincrement:\"true\"")
			case "int", "smallint":
				f.Type = "class.Int32"
				if arraykit.StringContains(es, "primary") {
					f.Type = "int32"
				}
			case "bigint":
				f.Type = "class.Int64"
				if arraykit.StringContains(es, "primary") {
					f.Type = "int64"
				}
			case "timestamp", "date":
				f.Type = "class.Time"
			case "jsonb":
				f.Type = "class.MapString"
			case "varchar[]", "text[]":
				f.Type = "class.ArrString"
			case "int[]":
				f.Type = "class.ArrInt"
			case "boolean":
				f.Type = "class.Bool"
			default:
				if strings.Index(es[1], "decimal") == 0 {
					f.Type = "class.Decimal"
				}
			}
			if strings.Index(f.Type, "class") >= 0 || strings.Index(f.Type, "*") >= 0 {
				f.Tags = append(f.Tags, fmt.Sprintf("json:\"%s,omitempty\" db:\"%s\"", es[0], strings.ToLower(es[0])))
			} else {
				f.Tags = append(f.Tags, fmt.Sprintf("json:\"%s\" db:\"%s\"", es[0], strings.ToLower(es[0])))
			}
			if arraykit.StringContains(es, "primary") {
				f.Tags = append(f.Tags, fmt.Sprintf("pk:\"true\" tablename:\"%s\"", table))
			}
			// 注释 -- 分隔
			commentIndex := strings.Index(val, "--")
			if commentIndex > 0 {
				comment := strings.TrimSpace(val[commentIndex+2:])
				if comment != "" {
					f.Tags = append(f.Tags, fmt.Sprintf(`description:"%s"`, comment))
				}
			}
			fields = append(fields, f)
		}
	}
	_ = filekit.WriteFile(destFile, []byte(dest))
}
