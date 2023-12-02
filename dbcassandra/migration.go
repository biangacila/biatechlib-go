package dbcassandra

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"sort"
	"strings"
)

type Migration struct {
	OriginHost      string
	DestinationHost string

	OriginDatabase      string
	DestinationDatabase string

	ReplicationFactor int

	tables                map[string]string
	tableCreateQuery      map[string]string
	SystemSchemaColumns   map[string]map[string]map[string]SystemSchemaColumns // db, table, field , type
	SystemSchemaKeyspaces map[string]SystemSchemaKeyspaces

	connOrigin  *gocql.Session
	connDest    *gocql.Session
	qryDbCreate string
	queries     map[string][]string
}

func (obj *Migration) Run() {
	obj.tables = make(map[string]string)
	obj.tableCreateQuery = make(map[string]string)
	obj.SystemSchemaKeyspaces = make(map[string]SystemSchemaKeyspaces)
	obj.SystemSchemaColumns = make(map[string]map[string]map[string]SystemSchemaColumns)
	obj.queries = make(map[string][]string)

	obj.initialDbConn()
	fmt.Println("1 (--> finish initialDbConn")
	//todo get tables name from db
	_ = obj.getAllTableNameFromDb()
	fmt.Println("2 (--> finish getAllTableNameFromDb")
	//todo get the table create query
	obj.getTableCreateQuery()
	fmt.Println("3 (--> finish getTableCreateQuery")
	fmt.Println("total queries > ", len(obj.queries))
	_ = obj.getSchemaInfoDatabase()
	fmt.Println("4 (--> finish getSchemaInfoDatabase")

	//todo let create now the table structure in destination host
	obj.createTablesStructureInDest()
	fmt.Println("5 (--> finish createTablesStructureInDest")

	//todo let fetch our data from origin database
	obj.FetchDataFromOrigin()
	fmt.Println("6 (--> finish FetchDataFromOrigin")

	//todo let insert data into destination server
	obj.insertDataIntoDestinationDb()
	fmt.Println("7 (--> finish insertDataIntoDestinationDb")

}
func (obj *Migration) insertDataIntoDestinationDb() {
	y := 1
	for table, queries := range obj.queries {
		fmt.Println(y, "--(::::Inserting into| ", table, " | total: ", len(queries))
		InsertBatchQueries(queries, obj.DestinationHost, obj.DestinationDatabase)
		lastIndex := len(queries) - 1
		lastQry := queries[lastIndex]
		_, _ = RunQueryCass2(obj.connDest, lastQry)
		y++
	}

}
func (obj *Migration) FetchDataFromOrigin() {
	var lat string
	for _, table := range obj.tables {
		var ls []interface{}
		//qry := "select * from system_schema.columns"
		b := FetchEntityFromDbAll(obj.connOrigin, table, obj.OriginDatabase, "")
		if err := json.Unmarshal(b, &ls); err != nil {
			fmt.Println("error > ", err)
			continue
		}
		for _, info := range ls {
			str, _ := json.Marshal(info)
			qry := fmt.Sprintf("insert into %v.%v  JSON '%v' ", obj.DestinationDatabase, table, string(str))
			obj.queries[table] = append(obj.queries[table], qry)
			lat = qry
		}
	}

	fmt.Println("LAst query> ", lat)
}
func (obj *Migration) createTablesStructureInDest() {
	_, _ = RunQueryCass2(obj.connDest, obj.qryDbCreate)
	for _, qry := range obj.tableCreateQuery {
		_, _ = RunQueryCass2(obj.connDest, qry)
	}
}
func (obj *Migration) getSchemaInfoDatabase() error {
	var ls []SystemSchemaKeyspaces
	/*qry := "select * from system_schema.keyspaces;"
	rs, _ := io.RunQueryCass2(qry)*/
	b := FetchEntityFromDbAll(obj.connOrigin, "keyspaces", "system_schema", "")

	if err := json.Unmarshal(b, &ls); err != nil {
		return err
	}

	for _, row := range ls {
		obj.SystemSchemaKeyspaces[row.Keyspace_name] = row
	}

	if row, ok := obj.SystemSchemaKeyspaces[obj.OriginDatabase]; ok {
		wDbase := obj.DestinationDatabase
		mClass := row.Replication["class"]
		//mFactor := row.Replication["replication_factor"]
		//mDurable := row.Durable_writes
		arr := strings.Split(mClass, ".")
		mClass = arr[(len(arr) - 1)]
		cqry := fmt.Sprintf(`CREATE KEYSPACE %v WITH replication = {'class': '%v', 'dc1': '%v'}  AND durable_writes = true`, wDbase, mClass, obj.ReplicationFactor)
		obj.qryDbCreate = cqry
	}
	return nil
}
func (obj *Migration) initialDbConn() {
	obj.connOrigin = CreateConnectionWithAuth(obj.OriginHost)
	obj.connDest = CreateConnectionWithAuth(obj.DestinationHost)
}
func (obj *Migration) getTableCreateQuery() {
	for tname, rowTable := range obj.SystemSchemaColumns[obj.OriginDatabase] {
		qry := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s(", obj.DestinationDatabase, tname)
		x := 0
		for field, row := range rowTable {
			comma := " ,"
			if x == 0 {
				comma = " "
			}
			inner := fmt.Sprintf("%s %s %s ", comma, field, row.Type)
			qry = qry + inner
			x++
		}
		strPrimaryKey, cout := obj.findClusteringTable(obj.OriginDatabase, tname)
		if cout == 0 {

		} else {
			strPrimaryKey = fmt.Sprintf(",PRIMARY KEY %v", strPrimaryKey)
		}
		qry = fmt.Sprintf("%v  %v );", qry, strPrimaryKey)
		obj.tableCreateQuery[tname] = qry
	}
}
func (obj *Migration) getAllTableNameFromDb() error {
	var ls []SystemSchemaColumns
	//qry := "select * from system_schema.columns"
	b := FetchEntityFromDbAll(obj.connOrigin, "columns", "system_schema", "")
	if err := json.Unmarshal(b, &ls); err != nil {
		return err
	}
	for _, row := range ls {
		if row.Keyspace_name == obj.OriginDatabase {
			if _, ok := obj.SystemSchemaColumns[row.Keyspace_name]; !ok {
				obj.SystemSchemaColumns[row.Keyspace_name] = make(map[string]map[string]SystemSchemaColumns)
			}
			if _, ok := obj.SystemSchemaColumns[row.Keyspace_name][row.Table_name]; !ok {
				obj.SystemSchemaColumns[row.Keyspace_name][row.Table_name] = make(map[string]SystemSchemaColumns)
			}
			obj.SystemSchemaColumns[row.Keyspace_name][row.Table_name][row.Column_name] = row
			obj.tables[row.Table_name] = row.Table_name
		}
	}
	return nil
}
func (obj *Migration) findClusteringTable(dbname, tablename string) (string, int) {
	count := 0
	mymap := make(map[int]string)
	arr := []int{}
	partition_key := ""
	for _, row := range obj.SystemSchemaColumns[dbname][tablename] {
		if strings.ToLower(row.Kind) == "clustering" {
			mymap[row.Position] = row.Column_name
			arr = append(arr, row.Position)
			count++
		}
		if strings.ToLower(row.Kind) == "partition_key" {
			partition_key = row.Column_name
		}
	}
	str := ""
	//let build our result
	if count > 0 || partition_key != "" {
		sort.Ints(arr)
		x := 0
		str = "("
		if partition_key != "" {
			str = str + " " + partition_key
			x++
		}
		for _, index := range arr {
			key, _ := mymap[index]
			if x == 0 {
				str = str + " " + key
			} else {
				str = str + ", " + key
			}
			x++
		}
		str = str + ")"
	}

	return str, count
}
