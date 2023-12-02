package dbcassandra

import (
	"encoding/json"
	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"github.com/gocql/gocql"
	"github.com/pborman/uuid"
	"log"
	"strings"
)

type Where struct {
	Key  string
	Val  interface{}
	Type string
}

func RunQueryCass2(session *gocql.Session, qry string) (string, error) {
	defer global.RecoverMe("crum queries RunQueryCass2")
	qResult := "[]"
	iter := session.Query(qry).Iter()
	myrow, err := iter.SliceMap()
	if err != nil {
		log.Println("RunQueryCass Cassandra  session.Query Error --->>> ", err, " > ", qry)
		return qResult, err
	}
	str, _ := json.Marshal(myrow)
	qResult = string(str)
	return qResult, nil
}
func RunQueryCass3(session *gocql.Session, qry string) ([]byte, error) {
	defer global.RecoverMe("crum queries RunQueryCass2")
	qResult := "[]"
	iter := session.Query(qry).Iter()
	myrow, err := iter.SliceMap()
	if err != nil {
		log.Println("RunQueryCass Cassandra  session.Query Error --->>> ", err, " > ", qry)
		return []byte(qResult), err
	}
	str, _ := json.Marshal(myrow)
	qResult = string(str)
	return []byte(qResult), nil
}
func RunQueryCass(session *gocql.Session, qry string, fack []string) (string, error) {
	qResult := "[]"
	iter := session.Query(qry).Iter()
	myrow, err := iter.SliceMap()
	if err != nil {
		log.Println("RunQueryCass Cassandra  session.Query Error --->>> ", err, " > ", qry)
		return qResult, err
	}
	str, _ := json.Marshal(myrow)
	qResult = string(str)
	return qResult, nil
}
func FetchEntityFromDbAll(session *gocql.Session, entity, dbName, appName string) []byte {
	qry := fmt.Sprintf("select * from %v.%v where appname='%v'",
		dbName, entity, appName)
	if appName == "" {
		qry = fmt.Sprintf("select * from %v.%v ",
			dbName, entity)
	}
	res, _ := RunQueryCass2(session, qry)
	return []byte(res)
}
func FetchEntityFromDbAllSpecificField(session *gocql.Session, entity, dbName, appName string, fields []string) []byte {
	selectionField := buildFieldList(fields)
	qry := fmt.Sprintf("select %v from %v.%v where appname='%v'",
		selectionField, dbName, entity, appName)
	if appName == "" {
		qry = fmt.Sprintf("select * from %v.%v ",
			dbName, entity)
	}
	res, _ := RunQueryCass2(session, qry)
	return []byte(res)
}
func UpdateQuery(session *gocql.Session, dbName, table string, params []Where, setParams []Where) string {
	qry := fmt.Sprintf("update %v.%v ", dbName, table)
	if len(setParams) > 0 {
		var x = 0
		for _, row := range setParams {
			if x == 0 {
				innerQry := fmt.Sprintf(" set %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf(" set %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			} else {
				innerQry := fmt.Sprintf(" , %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf(" , %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			}
			x++
		}

	}
	if len(params) > 0 {
		var x = 0
		for _, row := range params {
			if x == 0 {
				innerQry := fmt.Sprintf("where %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf("where %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			} else {
				innerQry := fmt.Sprintf(" and %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf(" and %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			}
			x++
		}
	}
	RunQueryCass2(session, qry)
	return "OK"
}
func DeleteQuery(session *gocql.Session, dbName, table string, params []Where) string {
	qry := fmt.Sprintf("delete  from %v.%v ", dbName, table)
	if len(params) > 0 {
		var x = 0
		for _, row := range params {
			if x == 0 {
				innerQry := fmt.Sprintf("where %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf("where %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			} else {
				innerQry := fmt.Sprintf(" and %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf(" and %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			}
			x++
		}
	}
	RunQueryCass2(session, qry)
	return "OK"
}
func SelectQuery(session *gocql.Session, dbName, table string, params []Where) []byte {
	qry := fmt.Sprintf("select * from %v.%v ", dbName, table)
	if len(params) > 0 {
		var x = 0
		for _, row := range params {
			if x == 0 {
				innerQry := fmt.Sprintf("where %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf("where %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			} else {
				innerQry := fmt.Sprintf(" and %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf(" and %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			}
			x++
		}
	}
	var ls []interface{}
	fmt.Println("LIST-QUERY ):(--> ", qry)
	res, _ := RunQueryCass2(session, qry+" ALLOW FILTERING")
	err := json.Unmarshal([]byte(res), &ls)
	str, err := json.Marshal(ls)
	if err != nil {
		fmt.Println("Error Unmarshal fetch data: ", err, qry)
		/*do nothing*/
	}

	return str
}
func InsertTable(session *gocql.Session, dbName, table string, objIn interface{}, queryOnly bool) string {
	in := make(map[string]interface{})
	str1, err := json.Marshal(objIn)
	if err != nil {
	}
	err = json.Unmarshal(str1, &in)
	if err != nil {
	}

	dt, hr := global.GetDateAndTimeString()
	if profile, ok := in["Profile"]; ok {
		if profile == nil {
			in["Profile"] = make(map[string]interface{})
		}
	}

	if orgDateTime, ok := in["OrgDateTime"]; ok {
		if fmt.Sprintf("%v", orgDateTime) == "" {
			in["OrgDateTime"] = fmt.Sprintf("%v %v", dt, hr)
		}
	}
	if Status, ok := in["Status"]; ok {
		if fmt.Sprintf("%v", Status) == "" {
			in["Status"] = "active"
		}
	}

	if _, ok := in["Date"]; ok {
		if isValueEmpty(in["Date"]) {
			in["Date"] = dt
		}
	}
	if _, ok := in["Time"]; ok {
		if isValueEmpty(in["Time"]) {
			in["Time"] = hr
		}
	}
	if _, ok := in["Id"]; ok {
		if isValueEmpty(in["Id"]) {
			in["Id"] = uuid.New()
		}

	}

	str, _ := json.Marshal(in)
	qry := fmt.Sprintf("insert into %v.%v  JSON '%v' ", dbName, table, string(str))
	if !queryOnly {
		_, err = RunQueryCass2(session, qry)
		return "OK"
	}

	return qry
}
func InsertTablePure(session *gocql.Session, dbName, table string, objIn interface{}) string {
	str, _ := json.Marshal(objIn)
	qry := fmt.Sprintf("insert into %v.%v  JSON '%v' ", dbName, table, string(str))
	_, _ = RunQueryCass2(session, qry)
	return "OK"
}
func isValueEmpty(in interface{}) bool {
	str := fmt.Sprintf("%v", in)
	if str == "" {
		return true
	}
	return false
}
func buildFieldList(fields []string) string {
	var line string
	for _, item := range fields {
		line += fmt.Sprintf("%v ,", item)
	}
	line = strings.Trim(line, ",")
	return line
}

func SelectQueryWithFieldList(session *gocql.Session, dbName, table string, params []Where, fieldList []string) []byte {
	qry := fmt.Sprintf("select * from %v.%v ", dbName, table)
	if len(fieldList) > 0 {
		fields := buildFieldList(fieldList)
		qry = fmt.Sprintf("select %v from %v.%v ", fields, dbName, table)
	}
	if len(params) > 0 {
		var x = 0
		for _, row := range params {
			if x == 0 {
				innerQry := fmt.Sprintf("where %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf("where %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			} else {
				innerQry := fmt.Sprintf(" and %v='%v' ", row.Key, row.Val)
				if row.Type != "string" {
					innerQry = fmt.Sprintf(" and %v=%v ", row.Key, row.Val)
				}
				qry = qry + innerQry
			}
			x++
		}
	}
	var ls []interface{}
	fmt.Println("LIST-QUERY ):(--> ", qry)
	res, _ := RunQueryCass2(session, qry+" ALLOW FILTERING")
	err := json.Unmarshal([]byte(res), &ls)
	str, err := json.Marshal(ls)
	if err != nil {
		fmt.Println("Error Unmarshal fetch data: ", err, qry)
		/*do nothing*/
	}

	return str
}
