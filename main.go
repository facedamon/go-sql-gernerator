package main

import (
	"flag"
	"fmt"
	"github.com/facedamon/go-sql-generator/conf"
	"github.com/facedamon/go-sql-generator/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var d *sqlx.DB

func init() {
	var err error
	v := conf.Config().Db
	d, err = sqlx.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
			v.User, v.Pwd, v.Ip, v.Port, v.Schema))
	if err != nil {
		panic(err)
	}
	err = d.Ping()
	if err != nil {
		panic(err)
	}
	d.SetMaxIdleConns(v.MaxIdle)
	d.SetMaxOpenConns(v.MaxConn)

}

func genModelFile(render *template.Template, packageName, tableName string) (err error) {
	var tableSchema []model.TableSchema

	err = d.Select(&tableSchema,
		fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE,COLUMN_KEY,"+
			"COLUMN_COMMENT from information_schema.COLUMNS "+
			"WHERE TABLE_NAME='%s' and TABLE_SCHEMA='%s'", tableName, conf.Config().Db.Schema))

	if err != nil {
		return
	}

	if len(tableSchema) <= 0 {
		return fmt.Errorf("TableSchema is nil")
	}

	fileName := *modelFolder + model.FormatTableName(tableName) + ".go"

	os.Remove(fileName)
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &model.Info{
		PackageName: packageName,
		ConnName:    *connName,
		TableName:   tableName,
		ModelName:   tableName,
		TableSchema: &tableSchema}

	if err := render.Execute(f, model); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	return
}

var h = flag.Bool("h", false, "this help")
var tplFile = flag.String("t", "./model.tpl", "the path of tpl file")
var modelFolder = flag.String("m", "./db/", "the path for folder of model files")
var genTable = flag.String("tabs", "", "the name of table to be generated, split with ','")
var packageName = flag.String("pkg", "db", "packageName")
var connName = flag.String("c", "d", "the name of sqlx.DB")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(io.MultiWriter(os.Stderr), `go-sql-generator version: v1.0
Usage: go-sql-generator [-h help] [-t tpl] [-tabs] [-pkg]

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()
	flag.Usage()
	if *h {
		return
	}

	logDir, _ := filepath.Abs(*modelFolder)
	if _, err := os.Stat(logDir); err != nil {
		os.Mkdir(logDir, os.ModePerm)
	}

	data, err := ioutil.ReadFile(*tplFile)
	if nil != err {
		fmt.Printf("Read tpl err: %s\n", err.Error())
		return
	}

	render := template.Must(template.New("model").
		Funcs(template.FuncMap{
			"FirstCharUpper":       model.FirstCharUpper,
			"TypeConvert":          model.TypeConvert,
			"Tags":                 model.Tags,
			"ExportColumn":         model.ExportColumn,
			"FormatTableName": 		model.FormatTableName,
			"Join":                 model.Join,
			"MakeQuestionMarkList": model.MakeQuestionMarkList,
			"ColumnAndType":        model.ColumnAndType,
			"ColumnWithPostfix":    model.ColumnWithPostfix,
		}).
		Parse(string(data)))

	var tablaNames []string
	if len(*genTable) > 0 {
		tablaNames = strings.Split(*genTable, ",")
	} else {
		err = d.Select(&tablaNames,
			fmt.Sprintf("SELECT table_name from information_schema.tables where table_schema='%s'",
				conf.Config().Db.Schema))
		if err != nil {
			fmt.Printf("Query table record error err: %s\n", err.Error())
			return
		}
	}

	for _, table := range tablaNames {
		err := genModelFile(render, *packageName, table)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
