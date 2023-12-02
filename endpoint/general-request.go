package endpoint

import (
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"net/http"
)

func WsCrudRequest[t any](w http.ResponseWriter, r *http.Request, dbName, table string, hasPrefix bool, prefixCode string, prefixInitial int, prefixField string, outs []t, session *gocql.Session) {
	vars := mux.Vars(r)
	action, _ := vars["action"]
	if action == "list" {
		WsGenericRequestPostList(w, r, session, dbName, table, outs)
	} else if action == "insert" {
		WsGenericRequestPostInsert(w, r, session, table, outs,
			prefixField, prefixCode, prefixInitial, hasPrefix, false)
	} else if action == "delete" {
		WsGenericRequestPostDelete(w, r, session, table, outs)
	} else if action == "find" {
		WsGenericRequestPostFind(w, r, session, dbName, table, outs)
	} else if action == "search" {
		WsGenericRequestPostSearch(w, r, session, dbName, table, outs)
	} else if action == "update" {
		WsGenericRequestPostUpdate(w, r, session, dbName, table, outs)
	}
}
