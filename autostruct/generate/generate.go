package generate

import (
	"MyProject/autostruct/conf"
	"MyProject/autostruct/dbtools"
	"fmt"
	"io"
	"os"
	"strings"
)

func Camel_Case(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCII_Lower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCII_Digit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCII_Lower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCII_Lower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

//重构源码此方法
func isASCII_Lower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

//重构源码此方法
func isASCII_Digit(c byte) bool {
	return '0' <= c && c <= '9'
}

func Generate(tableNames ...string) {
	tableNamesStr := ""
	for _, name := range tableNames {
		if tableNamesStr != "" {
			tableNamesStr += ","
		}
		tableNamesStr += "'" + name + "'"
	}
	tables := getTables(tableNamesStr) //生成所有表信息
	//tables := getTables("admin_info","video_info") //生成指定表信息，可变参数可传入过个表名
	for _, table := range tables {
		fields := getFields(table.Name)
		generateModel(table, fields)
	}
}

//获取表信息
func getTables(tableNames string) []Table {
	db := dbtools.GetMysqlDb()
	var tables []Table
	if tableNames == "" {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + conf.MasterDbConfig.DbName + "';").Find(&tables)
	} else {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (" + tableNames + ") AND table_schema='" + conf.MasterDbConfig.DbName + "';").Find(&tables)
	}
	return tables
}

//获取所有字段信息
func getFields(tableName string) []Field {
	db := dbtools.GetMysqlDb()
	var fields []Field
	db.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

//生成Model
func generateModel(table Table, fields []Field) {
	content := "package models\n\n"
	//表注释
	if len(table.Comment) > 0 {
		content += "// " + table.Comment + "\n"
	}
	content += "type " + Camel_Case(table.Name) + " struct {\n"
	//生成字段
	for _, field := range fields {
		fieldName := Camel_Case(field.Field)
		fieldJson := getFieldJson(field)
		fieldType := getFiledType(field)
		fieldComment := getFieldComment(field)
		content += "	" + fieldName + " " + fieldType + " `" + fieldJson + "` " + fieldComment + "\n"
	}
	content += "}"

	filename := conf.ModelPath + Camel_Case(table.Name) + ".go"
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		if !conf.ModelReplace {
			fmt.Println(Camel_Case(table.Name) + " 已存在，需删除才能重新生成...")
			return
		}
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(Camel_Case(table.Name) + " 已生成...")
	}
}

//获取字段类型
func getFiledType(field Field) string {
	typeArr := strings.Split(field.Type, "(")

	switch typeArr[0] {
	case "int":
		return "int"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float32"
	case "double":
		return "float32"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}

//获取字段json描述
func getFieldJson(field Field) string {
	return `json:"` + field.Field + `"`
}

//获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		return "// " + field.Comment
	}
	return ""
}

//检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
