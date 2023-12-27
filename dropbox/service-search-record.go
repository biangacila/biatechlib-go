package dropbox

import (
	"strings"
)

type ServiceSearchRecord struct {
	KeyQuery string
	Params   []string

	DROPBOX_CLIENTID string
	DROPBOX_TOKEN    string
	DROPBOX_SECRET   string

	rowData []SearchResult
	Results []SearchResult
}

func (obj *ServiceSearchRecord) Run() {
	obj.fetchData()
	obj.filterData()
}
func (obj *ServiceSearchRecord) filterData() {
	var ls []SearchResult
	for _, rec := range obj.rowData {
		boo := true
		pname := rec.Name
		pname = strings.ToLower(pname)
		pname = strings.Trim(pname, " ")
		for _, item := range obj.Params {
			item = strings.ToLower(item)
			item = strings.Trim(item, " ")
			if !strings.Contains(pname, item) {
				boo = false
			}
		}
		if boo {
			ls = append(ls, rec)
		}
	}
	obj.Results = ls
}
func (obj *ServiceSearchRecord) fetchData() {
	var hub HubDropbox
	hub.Company = ""
	hub.SearchQuery = obj.KeyQuery
	/*hub.DROPBOX_CLIENTID = obj.DROPBOX_CLIENTID
	hub.DROPBOX_TOKEN = constant.DROPBOX_TOKEN
	hub.DROPBOX_SECRET = constant.DROPBOX_SECRET*/
	ls := hub.Search()
	obj.rowData = ls
}

//CCIDIAL-IN--
