package main

import (
	"fmt"
	"mizuki/project/core-kit/library/filekit"
	"mizuki/project/core-kit/library/stringkit"
	"strings"
)

type field struct {
	Name string
	Type string
	Tags []string
}

func SQL2Struct(sqlFile, destFile string) {
	sqls, err := filekit.ReadString(sqlFile)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(sqls, "\n")
	var dest, temp string
	var fields []field
	var table string
	for i := 0; i < len(parts); i++ {
		val := strings.TrimSpace(parts[i])
		if strings.Index(val, "-") == 0 || strings.Index(val, "create index") == 0 || strings.Index(val, "insert") == 0 || strings.Index(val, "update") == 0 {
			continue
		}
		if strings.Index(val, "create table") == 0 {
			table = val[strings.Index(val, "create")+13 : len(val)-1]
			temp = "type " + stringkit.CamelCase(table) + " struct{\n"
			fields = []field{}
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
				// todo
				continue
			}
			f := field{Name: es[0]}
			f.Tags = append(f.Tags, fmt.Sprintf("json:\"%s\" db:\"%s\"", es[0], es[0]))
			if stringkit.ArrayContains(es, "primary") {
				f.Tags = append(f.Tags, fmt.Sprintf("pk:\"true\" tablename:\"%s\"", table))
			}
			switch es[1] {
			case "varchar", "text":
				f.Type = "class.String"
			case "serial":
				f.Type = "int32"
			case "bigserial":
				f.Type = "int64"
			case "int", "smallint":
				f.Type = "class.Int32"
			case "bigint":
				f.Type = "class.Int64"
			case "timestamp":
				f.Type = "class.Time"
			case "jsonb":
				f.Type = "class.MapString"
			case "varchar[]", "text[]":
				f.Type = "class.ArrString"
			case "int[]":
				f.Type = "class.ArrInt"
			case "boolean":
				f.Type = "bool"
			case "decimal":
				f.Type = "float64"
			}
			fields = append(fields, f)
		}
	}
	_ = filekit.WriteFile(destFile, []byte(dest))
}
