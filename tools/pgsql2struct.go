package main

import (
	"mizuki/project/core-kit/library/filekit"
	"mizuki/project/core-kit/library/stringkit"
	"strings"
)

func SQL2Struct(sqlFile, destFile string) {
	sqls, err := filekit.ReadString(sqlFile)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(sqls, "\n")
	var dest string
	temp := ""
	for i := 0; i < len(parts); i++ {
		val := strings.ToLower(strings.TrimSpace(parts[i]))
		if strings.Index(val, "-") == 0 || strings.Index(val, "create index") == 0 || strings.Index(val, "insert") == 0 || strings.Index(val, "update") == 0 {
			continue
		}
		if strings.Index(val, "create table") == 0 {
			table := val[strings.Index(val, "create")+13 : len(val)-1]
			temp = "type " + stringkit.CamelCase(table) + " struct{\n"
		} else if strings.Index(val, ")") == 0 {
			// end
			temp += "}\n\n"
			dest += temp
			temp = ""
		} else if temp != "" {
			es := stringkit.Split(val, "[ ,\t]+")
			if es[0] == "primary" {
				continue
			}
			temp += "\t" + stringkit.CamelCase(es[0]) + " "
			switch es[1] {
			case "varchar", "text":
				temp += "string "
			case "serial":
				temp += "int "
			case "bigserial":
				temp += "int64 "
			case "int", "smallint":
				temp += "class.Int32 "
			case "bigint":
				temp += "class.Int64 "
			case "timestamp":
				temp += "class.Time "
			case "jsonb":
				temp += "class.MapString "
			case "varchar[]", "text[]":
				temp += "class.ArrString "
			case "int[]":
				temp += "class.ArrInt "
			case "boolean":
				temp += "bool "
			case "decimal":
				temp += "float64 "
			}
			temp += " `json:\"" + es[0] + "\"`\n"
		}
	}
	_ = filekit.WriteFile(destFile, []byte(dest))
}
