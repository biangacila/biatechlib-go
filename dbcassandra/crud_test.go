package dbcassandra

import (
	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"github.com/gocql/gocql"
	"testing"
)

type Category struct {
	Name        string
	Description string
}
type User struct {
	AppName  string
	Id       string
	Org      string
	Username string
	Password string
	Role     string
	Name     string
	Surname  string
	FullName string
	Phone    string
	Email    string

	Profile     map[string]interface{}
	OrgDateTime string
	Date        string
	Time        string
	Status      string
}

func getSession(server ...string) *gocql.Session {
	var host = "voip.easipath.com"
	if len(server) > 0 {
		host = server[0]
	}
	session := CreateConnectionWithAuth(host)
	return session
}
func TestServiceGeneric_Find(t *testing.T) {
	session := getSession()

	table := "Accounts"
	dbName := "fxtrader"
	appName := "fxtrader"
	var outs []Category
	var conditions = make(map[string]interface{})
	conditions["Code"] = "ACC100100020"
	results, boo := GenericDbFindRecord(session, dbName, table, appName, "", outs, conditions)
	fmt.Println("----> ", boo)
	global.DisplayObject("results", results)
}
func TestServiceGeneric_Search(t *testing.T) {
	session := getSession("pmis.easipath.com")
	table := "User"
	dbName := "pmis"
	appName := "pmis"
	var outs []User
	results, boo := GenericDbSearchRecords(session, dbName, table, appName, "", outs, "peterd")
	fmt.Println("----> ", boo)
	global.DisplayObject("results", results)
}
func TestServiceGeneric_Insert(t *testing.T) {
	session := getSession()
	table := "DataType"
	dbName := "fxtrader"
	type DataType struct {
		AppName string
		Id      string
		Code    string
		Org     string
		Name    string
		Age     float64
	}
	var record DataType
	record.Age = 45
	//record.Org = "C103"
	record.Name = "Mr Bia"
	//record.Id = "ywywywyw"

	var validationEmpty = []string{"Name"}
	code, inserted, emptyError, returnQuery := GenericDbInsert(session, dbName, dbName, table, record, validationEmpty, "Code", "UC", 10012, true, true)
	fmt.Println("code :)-> ", code)
	fmt.Println("inserted :)-> ", inserted)
	fmt.Println("emptyError :)-> ", emptyError)
	fmt.Println("returnQuery :)-> ", returnQuery)
	return

}
func TestServiceGeneric_Delete(t *testing.T) {
	session := getSession()
	dbName := "fxtrader"
	appName := "fxtrader"
	type DataType struct {
		AppName string
		Id      string
		Org     string
		Name    string
		Age     float64
	}
	var record DataType
	record.Age = 45
	record.Org = "C103"
	record.Name = "Mr Bia"
	record.Id = "ywywywyw"

	var conditions = []string{"Org", "Name"}
	result := GenericDbDeleteOne(session, dbName, "Users", appName, "C103", record, conditions)
	fmt.Println(result)
	return

}
func TestServiceGeneric_List(t *testing.T) {
	session := getSession()
	dbName := "fxtrader"
	appName := "fxtrader"

	var conditions = make(map[string]interface{})
	table := "Category"
	org := ""
	var outs []Category

	results := GenericDbRead2(session, dbName, table, appName, org, outs, conditions)
	global.DisplayObject("result", results)
}

func TestServiceGeneric_Update(t *testing.T) {
	//[{"Username":"peterd@marginmentor.co.za","Name":"John","Surname":"Doe","Code":"","Email":"peterd@marginmentor.co.za","Org":"","Phone":"0720882572"}]
	//
	session := getSession("pmis.easipath.com")
	dbName := "fleetminder2"
	table := "user"
	appName := "pmis"
	org := ""
	password := global.GetMd5("phd@1964")
	username := "peterd@marginmentor.co.za"

	conditions := map[string]interface{}{"Username": "biangacila@gmail.com"}
	values := map[string]interface{}{"Password": password}
	GenericDbUpdate(session, dbName, table, appName, org, []User{}, conditions, values)

	loginString := fmt.Sprintf("update user set password='%v' where  username='%v'",
		password, username)

	fmt.Println(loginString)
}

//CREATE KEYSPACE eventsourcing WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1': '1'}  AND durable_writes = true;
