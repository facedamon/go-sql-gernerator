package model

import (
	"html/template"
	"strings"
)

// model info
type Info struct {
	TableName    string
	PackageName  string
	ModelName    string
	ConnName	 string
	TableSchema  *[]TableSchema
}

// info_schema properties
type TableSchema struct {
	ColumnName    string `db:"COLUMN_NAME" json:"column_name"`
	DataType      string `db:"DATA_TYPE" json:"data_type"`
	ColumnKey     string `db:"COLUMN_KEY" json:"column_key"`
	ColumnComment string `db:"COLUMN_COMMENT" json:"COLUMN_COMMENT"`
}

// get info_schema all columns
func (m *Info) ColumnNames() []string {
	result := make([]string, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		result = append(result, t.ColumnName)
	}
	return result
}

// get info_schema count of columns
func (m *Info) ColumnCount() int {
	return len(*m.TableSchema)
}

// get info_schema pk of columns
func (m *Info) PkColumnsSchema() []TableSchema {
	result := make([]TableSchema, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.ColumnKey == "PRI" {
			result = append(result, t)
		}
	}
	return result
}

// is there have pk
func (m *Info) HavePk() bool {
	return len(m.PkColumnsSchema()) > 0
}

// get info_schema no pk of columns
func (m *Info) NoPkColumnsSchema() []TableSchema {
	result := make([]TableSchema, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.ColumnKey != "PRI" {
			result = append(result, t)
		}
	}
	return result
}

// get info_schema tableNames of no pk on all columns
func (m *Info) NoPkColumns() []string {
	noPkColumnsSchema := m.NoPkColumnsSchema()
	result := make([]string, 0, len(noPkColumnsSchema))
	for _, t := range noPkColumnsSchema {
		result = append(result, t.ColumnName)
	}
	return result
}

// get info_schema tableNames ok pk on all columns
func (m *Info) PkColumns() []string {
	pkColumnsSchema := m.PkColumnsSchema()
	result := make([]string, 0, len(pkColumnsSchema))
	for _, t := range pkColumnsSchema {
		result = append(result, t.ColumnName)
	}
	return result
}

// format title
func FirstCharLower(str string) string {
	if len(str) > 0 {
		return strings.ToLower(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

// format title
func FirstCharUpper(str string) string {
	if len(str) > 0 {
		return strings.ToUpper(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

// get title tags
func Tags(columnName string) template.HTML {
	return template.HTML("`db:" + `"` + columnName + `"` +
		" json:" + `"` + columnName + "\"`")
}

func Join(str []string, sep string) string {
	return strings.Join(str, sep)
}

// get export columns
func ExportColumn(columnName string) string {
	columnItems := strings.Split(columnName, "_")
	columnItems[0] = FirstCharUpper(columnItems[0])
	for i := 0; i < len(columnItems); i++ {
		item := strings.Title(columnItems[i])
		if strings.ToUpper(item) == "ID" {
			item = "ID"
		}
		columnItems[i] = item
	}
	return strings.Join(columnItems, "")
}

func FormatTableName(tableName string) string {
	ts := strings.Split(tableName, "_")
	ts[0] = FirstCharLower(ts[0])
	for i := 1; i < len(ts); i++ {
		ts[i] = strings.Title(ts[i])
	}
	return strings.Join(ts, "")
}

// type convert
func TypeConvert(str string) string {
	switch str {
	case "smallint", "tinyint":
		return "sql.NullInt32"

	case "varchar", "text", "longtext", "char":
		return "sql.NullString"

	case "date", "timestamp", "datetime":
		return "sql.NullTime"

	case "int", "bigint":
		return "sql.NullInt64"

	case "float", "double", "decimal":
		return "sql.NullFloat64"

	default:
		return "sql.NullString"
	}
}

// columnName type, columnName type, ...
func ColumnAndType(TableSchema []TableSchema) string {
	result := make([]string, 0, len(TableSchema))
	for _, t := range TableSchema {
		result = append(result, t.ColumnName+" "+TypeConvert(t.DataType))
	}
	return strings.Join(result, ",")
}

func ColumnWithPostfix(columns []string, Postfix, sep string) string {
	result := make([]string, 0, len(columns))
	for _, t := range columns {
		result = append(result, t+Postfix)
	}
	return strings.Join(result, sep)
}

func MakeQuestionMarkList(num int) string {
	a := strings.Repeat("?,", num)
	return a[:len(a)-1]
}
