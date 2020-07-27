package {{.PackageName}}

import (
    "bytes"
    "database/sql"
    "fmt"
    "github.com/jmoiron/sqlx"
)
{{$exportModelName := .ModelName | FormatTableName}}

type {{$exportModelName}} struct {
{{range .TableSchema}} {{.ColumnName | ExportColumn}} {{.DataType | TypeConvert}} {{.ColumnName | Tags}} // {{.ColumnComment}}
{{end}}}

var Default{{$exportModelName}} = &{{$exportModelName}}{}

// transact 封装事务控制
// example:
// func (s service) dosomething() error {
//    return transact(s.db, func (tx *sql.Tx) error {
//        if _, err := tx.Exec(...); err != nil {
//            return err
//        }
//        if _, err := tx.Exec(...); err != nil {
//            return err
//        }
//        return nil
//    })
//}
func transact(txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := {{.ConnName}}.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}

{{if .HavePk}}
func (m *{{$exportModelName}}) GetByPK({{.PkColumnsSchema | ColumnAndType}}) (*{{$exportModelName}}, error) {
	obj := &{{$exportModelName}}{}
	sql := "select * from {{.TableName}} where {{ColumnWithPostfix .PkColumns "=?" " and "}}"
	err := {{.ConnName}}.Get(obj, sql,{{range $K:=.PkColumns}}{{$K}},{{end}})

	if err != nil {
	    if err == sql.ErrNoRows {
    			return nil, nil
    		}
		return nil, err
	}
	return obj, nil
}
{{end}}



func (m *{{$exportModelName}}) InsertTx() (int64, error) {
	sql := "insert into {{.TableName}}({{Join .ColumnNames ","}}) values({{.ColumnCount | MakeQuestionMarkList}})"
	var affected int64
	err := transact(func(tx *sqlx.Tx) error {
        result, err := {{.ConnName}}.Exec(sql,
                {{range .TableSchema}}m.{{.ColumnName | ExportColumn}},
                {{end}}
            )
            if err != nil {
            		return err
            	}
        affected, _ = result.RowsAffected()
        return nil
	})
	if err != nil {
	    return -1, err
	}

	return affected, nil
}

{{if .HavePk}}
func (m *{{$exportModelName}}) DeleteByPk() error {
	sql := `delete from {{.TableName}} where {{ColumnWithPostfix .PkColumns "=?" " and "}}`

	err := transact(func(tx *sqlx.Tx) error {
        _, err := {{.ConnName}}.Exec(sql,
                {{range .PkColumns}}m.{{. | ExportColumn}},
                {{end}}
            )
            if err != nil {
                return err
            }
            return nil
	})

	if err != nil {
	    return err
	}
	return nil
}
{{end}}

{{if .HavePk}}
func (m *{{$exportModelName}}) UpdateByPk() error {
	sql := `update {{.TableName}} set {{ColumnWithPostfix .NoPkColumns "=?" ","}} where {{ColumnWithPostfix .PkColumns "=?" " and "}}`
	err := transact(func(tx *sqlx.Tx) error {
        _, err := {{.ConnName}}.Exec(sql,
                {{range .NoPkColumns}}m.{{. | ExportColumn}},
                {{end}}{{range .PkColumns}}m.{{. | ExportColumn}},
                {{end}}
            )

            if err != nil {
                return err
            }
            return nil
	})
    if err != nil {
        return err
    }
    return nil
}
{{end}}

func (m *{{$exportModelName}}) QueryByMap(ma map[string]interface{}) ([]*{{$exportModelName}}, error) {
	var result []*{{$exportModelName}}
	var params []interface{}

	sql := bytes.NewBufferString("select * from {{.TableName}} where 1=1 ")
	for k, v := range ma {
		sql.WriteString(fmt.Sprintf(" and %s=? ", k))
		params = append(params, v)
	}
	err := {{.ConnName}}.Select(&result, sql.String(), params...)
	if err != nil {
	    if err == sql.ErrNoRows {
    			return nil, nil
    		}
		return nil, err
	}
	return result, nil
}

func (m *{{$exportModelName}}) SliceScanByMap(ma map[string]interface{}) ([][]interface{}, error) {
    var result [][]interface{}
    var params []interface{}

    sql := bytes.NewBufferString("select * from {{.TableName}} where 1=1 ")
    for k, v := range ma {
        sql.WriteString(fmt.Sprintf(" and %s=? ", k))
        params = append(params, v)
    }
    nStmt, err := {{.ConnName}}.Preparex(sql.String())
    if err != nil {
    		return nil, err
    	}
    rows, err := nStmt.Queryx(params...)
    if err != nil {
        if err == sql.ErrNoRows{
            return nil, nil
        }
        return nil, err
    }
    	for rows.Next() {
    		cols, err := rows.SliceScan()
    		if err != nil {
                return nil, err
    		}
    		result = append(result, cols)
    	}
    return result, nil
}

func(m *{{$exportModelName}}) SliceMapByMap(ma map[string]interface{}) ([]map[string]interface{}, error) {
    var result []map[string]interface{}
    var params []interface{}

    sql := bytes.NewBufferString("select * from {{.TableName}} where 1=1 ")
    for k, v := range ma {
        sql.WriteString(fmt.Sprintf(" and %s=? ", k))
        params = append(params, v)
    }
    nStmt, err := {{.ConnName}}.Preparex(sql.String())
    if err != nil {
    		return nil, err
    	}
    rows, err := nStmt.Queryx(params...)
    if err != nil {
        if err == sql.ErrNoRows{
            return nil, nil
        }
        return nil, err
    }
    	for rows.Next() {
    		col := make(map[string]interface{})
            err := rows.MapScan(col)
    		if err != nil {
                return nil, err
    		}
    		for k, v := range col {
            			switch v.(type) {
            			case []uint8:
            				col[k] = string(v.([]uint8))
            			default:
            			}
            		}
    		result = append(result, col)
    	}
    return result, nil
}