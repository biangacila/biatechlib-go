package dbcassandra

import (
	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"os"
	"reflect"
)

type EntryTable struct {
	Table      string
	PrimaryKey string
	EntityStr  []StructType
}

type ServiceKeyspace struct {
	Tables   []EntryTable
	Host     string
	Keyspace string
	HasFile  bool

	queryTable []string
	queryDB    string
}

func (obj *ServiceKeyspace) Run() {
	obj.CreateKeyspace()
	for _, row := range obj.Tables {
		obj.CreateTable(row.Table, row.PrimaryKey, row.EntityStr)
	}

	if obj.HasFile {
		obj.createFile()
	}
}
func (obj *ServiceKeyspace) createFile() {
	fname := fmt.Sprintf("./keyspace-%v.sql", obj.Keyspace)
	os.Remove(fname)
	//global.CreateFolderUploadIfNotExist(fname)
	global.WriteNewLineToLogFile(obj.queryDB, fname)
	global.WriteNewLineToLogFile("\n", fname)
	for _, line := range obj.queryTable {
		global.WriteNewLineToLogFile(line, fname)
		global.WriteNewLineToLogFile("\n", fname)
	}
}
func (obj *ServiceKeyspace) CreateKeyspace() {
	qry := fmt.Sprintf("CREATE KEYSPACE %v WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1': '1'}  AND durable_writes = true;", obj.Keyspace)
	obj.queryDB = qry
}
func (obj *ServiceKeyspace) CreateTable(table, primaryKey string, sTypeList []StructType) {
	qry := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v.%v (\n", obj.Keyspace, table)

	var x = 0
	for _, row := range sTypeList {
		key := row.Key
		mp := obj.fnWitchDataType3(row.Type, key)
		qry = qry + fmt.Sprintf("\t %v %v, \n", key, mp)
		x++
		//todo delete after
	}
	qry = qry + fmt.Sprintf("\t PRIMARY KEY(%v) \n", primaryKey)
	qry = qry + fmt.Sprintf(");")
	obj.queryTable = append(obj.queryTable, qry)
}

type StructType struct {
	Field int
	Key   string
	Value string
	Type  string
}

func GetStructType(val reflect.Value) []StructType {
	var arr []StructType
	typeOfTstObj := val.Type()
	global.DisplayObject("TTTTT", typeOfTstObj)
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Field(i)
		o := StructType{}
		o.Field = i
		o.Key = typeOfTstObj.Field(i).Name
		o.Value = fmt.Sprintf("%v", fieldType.Interface())
		o.Type = fmt.Sprintf("%v", fieldType.Type())
		arr = append(arr, o)
	}
	return arr
}
func (obj *ServiceKeyspace) fnWitchDataType3(dataType, key string) string {

	switch dataType {
	case "int":
		return "text"
	case "int64":
		return "int"
	case "int32":
		return "int"
	case "float64":
		return "float"
	case "float32":
		return "float"
	case "string":
		return "text"
	case "bool":
		return "boolean"
	case "map[string]float64":
		return "map<text,float>"
	case "map[string]string":
		return "map<text,text>"
	case "map[string]interface {}":
		return "map<text,text>"
	default:
		return fmt.Sprintf("frozen<%v >", key)
	}
}
func (obj *ServiceKeyspace) fnWitchDataType2(key string, list []StructType) string {
	o := StructType{}
	for _, row := range list {
		if row.Key == key {
			o = row
		}
	}

	switch o.Type {
	case "int":
		return "text"
	case "int64":
		return "int"
	case "int32":
		return "int"
	case "float64":
		return "float"
	case "float32":
		return "float"
	case "string":
		return "text"
	case "bool":
		return "boolean"
	case "map[string]float64":
		return "map<text,float>"
	case "map[string]string":
		return "map<text,text>"
	case "map[string]interface {}":
		return "map<text,text>"
	default:
		return fmt.Sprintf("frozen<%v >", o.Key)
	}
}
func (obj *ServiceKeyspace) fnWitchDataType(myInterface interface{}) string {
	switch myInterface.(type) {

	case int:
		return "text"
	case int64:
		return "int"
	case int32:
		return "int"
	case float64:
		return "float"
	case float32:
		return "float"
	case string:
		return "text"
	case bool:
		return "boolean"
	case map[string]float64:
		return "map<text,float>"
	case map[string]string:
		return "map<text,text>"
	case map[string]interface{}:
		return "map<text,text>"
	default:
		return "map<text,text>"
	}
}
