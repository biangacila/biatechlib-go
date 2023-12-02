package dbcassandra

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"strings"
)

func GenericDbSearchRecords[t any](session *gocql.Session, dbName, table string, appName, org string, outs []t, searchKey string) ([]t, bool) {
	var boo bool
	var conditions = make(map[string]interface{})
	if org == "" {
		//org = io.ORG_NAME
	}
	var clearRecord = func(valueIn t) t {
		var mapRecord = make(map[string]interface{})
		str, _ := json.Marshal(valueIn)
		_ = json.Unmarshal(str, &mapRecord)
		for key, _ := range mapRecord {
			if key == "AppName" || key == "Id" || key == "Date" || key == "OrgDateTime" || key == "Time" {
				mapRecord[key] = ""
			}
		}
		var rec t
		str, _ = json.Marshal(mapRecord)
		_ = json.Unmarshal(str, &rec)
		return rec
	}
	var records []t
	var results []t
	records = GenericDbRead2(session, dbName, table, appName, org, outs, conditions)
	for _, row := range records {
		b, _ := json.Marshal(clearRecord(row))
		strObj := string(b)
		if strings.Contains(strObj, searchKey) {
			results = append(results, row)
			boo = true
		}
	}
	return results, boo
}
func GenericDbFindRecord[t any](session *gocql.Session, dbName, table string, appName, org string, outs []t, conditions map[string]interface{}) (t, bool) {
	var findRecord t
	var boo bool

	var records []t
	records = GenericDbRead2(session, dbName, table, appName, org, outs, conditions)
	if len(records) > 0 {
		findRecord = records[0]
		boo = true
	}
	return findRecord, boo
}
func GenericDbInsert[t any](session *gocql.Session, appName, dbName, table string, recordIn t, validationEmpty []string,
	prefixField, prefixCode string, prefixInitial int, hasPrefix bool, queryOnly bool) (string, bool, []string, string) {
	var code = ""
	var emptyError []string
	var returnQuery string
	//let generate prefix code is set
	if hasPrefix {
		code = GenerateAutoNumber(session, dbName, appName, table, prefixField, prefixCode, prefixInitial)
	}
	var compareRecord = make(map[string]interface{})
	str, _ := json.Marshal(recordIn)
	_ = json.Unmarshal(str, &compareRecord)
	for key, val := range compareRecord {
		if hasPrefix {
			if key == prefixField {
				compareRecord[key] = code
			}
		}
		// let check for validation empty
		if IsInArrayString(validationEmpty, key) {
			if IsNumberValue(val) {
				if ToFloat64(val) == 0 {
					emptyError = append(emptyError, fmt.Sprintf("%v is zero value", key))
				}
			} else if IsStringValue(val) {
				if ToString(val) == "" {
					emptyError = append(emptyError, fmt.Sprintf("%v has empty value", key))
				}
			}
		}
	}
	if len(emptyError) > 0 {
		return code, false, emptyError, returnQuery
	}
	//let now insert into db
	if !queryOnly {
		InsertTable(session, dbName, table, compareRecord, false)
	} else {
		returnQuery = InsertTable(session, dbName, table, compareRecord, true)
	}
	return code, true, emptyError, returnQuery
}
func GenericDbDeleteOne[t any](session *gocql.Session, dbName, table, appName, org string, recordIn t, conditions []string) string {
	var qry = fmt.Sprintf("delete from %v.%v where appName='%v' ",
		dbName, table, appName)
	if org != "" {
		qry = fmt.Sprintf("delete from %v.%v where appName='%v' and org='%v'",
			dbName, table, appName, org)
	}
	mapFields := make(map[string]interface{})
	str, _ := json.Marshal(recordIn)
	_ = json.Unmarshal(str, &mapFields)

	for key, val := range mapFields {
		if !IsInArrayString(conditions, key) {
			continue
		}
		var myType = fmt.Sprintf("%T", val)
		if myType == "int" || myType == "int64" || myType == "int32" || myType == "float64" ||
			myType == "float34" {
			qry += fmt.Sprintf(" and %v=%v ", key, val)
		} else {
			qry += fmt.Sprintf(" and %v='%v' ", key, val)
		}
	}
	fmt.Println("---> Deleted qry :(( ", qry)
	// let run queries
	_, err := RunQueryCass2(session, qry)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("OK")
}
func GenericDbUpdatePureArray[t any](session *gocql.Session, dbName, table string, outs []t) string {
	var numberOfRowAffected int
	for _, record := range outs {
		InsertTablePure(session, dbName, table, record)
		numberOfRowAffected++
	}
	return fmt.Sprintf("OK rows affected %v", numberOfRowAffected)
}
func GenericDbUpdatePureOne[t any](session *gocql.Session, dbName, table string, record t) string {
	var numberOfRowAffected int
	InsertTablePure(session, dbName, table, record)
	numberOfRowAffected++
	return fmt.Sprintf("OK rows affected %v", numberOfRowAffected)
}
func GenericDbUpdate[t any](session *gocql.Session, dbName, table, appName, org string, outs []t, conditions map[string]interface{}, values map[string]interface{}) string {

	qry := fmt.Sprintf("select * from %v.%v %v", dbName, table, BuildWhereQuery(conditions))
	res, _ := RunQueryCass2(session, qry)
	var storeRecords []t
	_ = json.Unmarshal([]byte(res), &storeRecords)

	var numberOfRowAffected = len(storeRecords)

	var records []interface{}
	for _, rec := range storeRecords {
		var compareRecord = make(map[string]interface{})
		str, _ := json.Marshal(rec)
		_ = json.Unmarshal(str, &compareRecord)
		for k, v := range values {
			if _, ok := compareRecord[k]; ok {
				compareRecord[k] = v
			}
		}
		records = append(records, compareRecord)
		InsertTablePure(session, dbName, table, compareRecord)
	}
	return fmt.Sprintf("OK rows affected %v", numberOfRowAffected)
}
func GenericDbRead[t any](session *gocql.Session, dbName, table, appName, org string, outs []t, conditions map[string]interface{}) any {
	var strData = FetchEntityDataFromDb(session, dbName, table, appName, org)
	err := json.Unmarshal([]byte(strData), &outs)
	if err != nil {
		fmt.Println("error>> ", err)
	}
	var records []t
	for _, rec := range outs {
		var compareRecord = make(map[string]interface{})
		str, _ := json.Marshal(rec)
		_ = json.Unmarshal(str, &compareRecord)
		boo := true
		// let compare with any only
		for keyStore, valueStore := range conditions {
			recValue := compareRecord[keyStore]
			var val interface{}
			str, _ := json.Marshal(recValue)
			errV := json.Unmarshal(str, &val)
			if errV != nil {
				fmt.Println("errV: ", errV)
			}
			var myType = fmt.Sprintf("%T", valueStore)
			if myType == "int" || myType == "int64" || myType == "int32" || myType == "float64" ||
				myType == "float34" {
				boo = CompareEqualNumbers(val, valueStore)
			} else if myType == "bool" {
				boo = CompareEqualBool(val, valueStore)
			} else {
				if val != valueStore {
					boo = false
				}
			}
		}
		if !boo {
			continue
		}
		records = append(records, rec)
	}
	return records
}
func GenericDbRead2[t any](session *gocql.Session, dbName, table, appName, org string, outs []t, conditions map[string]interface{}) []t {
	var strData = FetchEntityDataFromDb(session, dbName, table, appName, org)
	err := json.Unmarshal([]byte(strData), &outs)
	if err != nil {
		fmt.Println("error>> ", err)
	}
	var records = GenericRecordsMatch(outs, conditions)
	return records
}
func GenericRecordsMatch[t any](outs []t, conditions map[string]interface{}) []t {
	var records []t
	for _, rec := range outs {
		var compareRecord = make(map[string]interface{})
		str, _ := json.Marshal(rec)
		_ = json.Unmarshal(str, &compareRecord)
		boo := true
		// let compare with any only
		for keyStore, valueStore := range conditions {
			recValue := compareRecord[keyStore]
			var val interface{}
			str, _ := json.Marshal(recValue)
			errV := json.Unmarshal(str, &val)
			if errV != nil {
				fmt.Println("errV: ", errV)
			}
			var myType = fmt.Sprintf("%T", valueStore)
			if myType == "int" || myType == "int64" || myType == "int32" || myType == "float64" ||
				myType == "float34" {
				boo = CompareEqualNumbers(val, valueStore)
			} else if myType == "bool" {
				boo = CompareEqualBool(val, valueStore)
			} else {
				if val != valueStore {
					boo = false
				}
			}
		}
		if !boo {
			continue
		}
		records = append(records, rec)
	}
	return records
}
