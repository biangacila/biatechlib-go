package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/biangacila/biatechlib-go/dbcassandra"
	"github.com/biangacila/luvungula-go/global"
	"github.com/gocql/gocql"
	"net/http"
)

type PostRequest struct {
	AppName            string
	Org                string
	Conditions         map[string]interface{}
	ConditionFields    map[string]interface{}
	ConditionValues    map[string]interface{}
	ConditionFieldList []string
	DbName             string
	Data               any
}
type PostResponse struct {
	Data         any
	Errors       []string
	AffectedRows int
	HasData      bool
	HasError     bool
	Ref          string
}

func WsGenericRequestPostInsert[t any](w http.ResponseWriter, r *http.Request, session *gocql.Session, table string, outs t, prefixField, prefixCode string, prefixInitial int, hasPrefix, queryOnly bool) {
	_, dataString := global.GetPostedDataMapAndString(r)
	var o PostRequest
	err := json.Unmarshal([]byte(dataString), &o)
	if err != nil {
		fmt.Println("Error WsGenericRequestPostInsert : ", err)
	}
	//let convert record
	var record = outs
	str, _ := json.Marshal(o.Data)
	_ = json.Unmarshal(str, &record)
	code, inserted, emptyError, returnQuery := dbcassandra.GenericDbInsert(session, o.AppName, o.DbName, table, record, o.ConditionFieldList, prefixField, prefixCode, prefixInitial, hasPrefix, queryOnly)
	var response PostResponse
	response.Data = returnQuery
	response.HasError = len(emptyError) > 0
	response.Ref = code
	response.HasData = inserted

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = o
	my["RESULT"] = response
	global.PublishToReact(w, r, my, 200)
}
func WsGenericRequestPostDelete[t any](w http.ResponseWriter, r *http.Request, session *gocql.Session, table string, outs t) {
	_, dataString := global.GetPostedDataMapAndString(r)
	var o PostRequest
	err := json.Unmarshal([]byte(dataString), &o)
	if err != nil {
		fmt.Println("Error WsGenericRequestPostDelete : ", err)
	}

	//let convert record
	var record = outs
	str, _ := json.Marshal(o.Data)
	_ = json.Unmarshal(str, &record)

	results := dbcassandra.GenericDbDeleteOne(session, o.DbName, table, o.AppName, o.Org, record, o.ConditionFieldList)
	var response PostResponse
	response.Data = results
	response.HasData = len(results) > 0
	if results != "OK" {
		response.HasError = true
		response.Errors = append(response.Errors, results)
	}
	response.HasError = results != "OK"

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = o
	my["RESULT"] = response
	global.PublishToReact(w, r, my, 200)
}
func WsGenericRequestPostList[t any](w http.ResponseWriter, r *http.Request, session *gocql.Session, dbName, table string, outs []t) {
	_, dataString := global.GetPostedDataMapAndString(r)
	var o PostRequest
	err := json.Unmarshal([]byte(dataString), &o)
	if err != nil {
		fmt.Println("Error WsGenericFetchDataList : ", err)
	}
	global.DisplayObject("validation", o)
	results := dbcassandra.GenericDbRead2(session, dbName, table, dbName, o.Org, outs, o.Conditions)
	var response PostResponse
	response.Data = results
	response.HasData = len(results) > 0
	response.HasError = false

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = o
	my["RESULT"] = response
	global.PublishToReact(w, r, my, 200)
}
func WsGenericRequestPostFind[t any](w http.ResponseWriter, r *http.Request, session *gocql.Session, dbName, table string, outs []t) {
	_, dataString := global.GetPostedDataMapAndString(r)
	var o PostRequest
	err := json.Unmarshal([]byte(dataString), &o)
	if err != nil {
		fmt.Println("Error WsGenericRequestPostFind : ", err)
	}

	results, boo := dbcassandra.GenericDbFindRecord(session, dbName, table, dbName, o.Org, outs, o.Conditions)
	var response PostResponse
	response.Data = results
	response.HasData = boo
	response.HasError = false

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = o
	my["RESULT"] = response
	global.PublishToReact(w, r, my, 200)
}
func WsGenericRequestPostSearch[t any](w http.ResponseWriter, r *http.Request, session *gocql.Session, dbName, table string, outs []t) {
	_, dataString := global.GetPostedDataMapAndString(r)
	var o PostRequest
	err := json.Unmarshal([]byte(dataString), &o)
	if err != nil {
		fmt.Println("Error WsGenericRequestPostSearch : ", err)
	}

	var searchKey = dbcassandra.ToString(o.Data)
	results, boo := dbcassandra.GenericDbSearchRecords(session, dbName, table, dbName, o.Org, outs, searchKey)
	var response PostResponse
	response.Data = results
	response.HasData = boo
	response.HasError = false

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = o
	my["RESULT"] = response
	global.PublishToReact(w, r, my, 200)
}
func WsGenericRequestPostUpdate[t any](w http.ResponseWriter, r *http.Request, session *gocql.Session, dbName, table string, outs []t) {
	_, dataString := global.GetPostedDataMapAndString(r)
	var o PostRequest
	err := json.Unmarshal([]byte(dataString), &o)
	if err != nil {
		fmt.Println("Error WsGenericRequestPostSearch : ", err)
	}
	results := dbcassandra.GenericDbUpdate(session, dbName, table, dbName, o.Org, outs, o.ConditionFields, o.ConditionValues)
	var response PostResponse
	response.Data = results
	response.HasData = len(results) > 0
	response.HasError = false

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = o
	my["RESULT"] = response
	global.PublishToReact(w, r, my, 200)
}
