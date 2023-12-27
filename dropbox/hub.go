package dropbox

import (
	"fmt"
	"github.com/pborman/uuid"
	"log"
	"net/http"
	"net/http/httptest"
	//"reflect"
	//"encoding/json"
)

type SearchResult struct {
	Id       string
	Name     string
	Path     string
	Revision string
	Size     int
	Tag      string
	Modified string
	ParentId string
}

type Entry struct {
	UID      string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Path     string `json:"path_display,omitempty"`
	Revision string `json:"rev,omitempty"`
	Size     int    `json:"size,omitempty"`
	Tag      string `json:".tag"`
}

type HubDropbox struct {
	//todo User info
	Company string

	//todo connect portion
	Db          *Dropbox
	MyAppKey    string
	MyAppSecret string
	MyAppUrl    string
	MyAppToken  string

	//todo Search portion
	SearchPath   string
	SearchQuery  string
	SearchResult []SearchResult

	//todo Upload file portion
	UploadModule   string
	UploadBusket   string
	UploadDst      string
	UploadSrc      string
	UploadFilename string
	UploadResult   Entry

	//todo Download file portion
	DownloadRev    string
	DownloadFolder string
	DownloadLink   string
	DownloadError  error
	//todo Metadata info
	MetadataRev       string
	MetadataResultOne Entry

	//todo AUTHENTICATION
	DROPBOX_CLIENTID string
	DROPBOX_SECRET   string
	DROPBOX_TOKEN    string
}

func (obj *HubDropbox) Connect() {

	obj.MyAppKey = obj.DROPBOX_CLIENTID
	obj.MyAppSecret = obj.DROPBOX_SECRET
	obj.MyAppToken = obj.DROPBOX_TOKEN
	db := NewDropbox()

	var sampleResponse = fmt.Sprintf("{\"access_token\":\"%s\",\"token_type\":\"bearer\",\"uid\":\"%s\"}",
		obj.MyAppToken, uuid.New())
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, sampleResponse)
	}))

	db.SetAppInfo(obj.MyAppKey, obj.MyAppSecret, testServer.URL)
	db.oAuth2Handler.authURL = testServer.URL
	db.SetAccessToken(obj.MyAppToken)
	db.mediaURL = testServer.URL
	obj.Db = db
}
func (obj *HubDropbox) ListFolder() {
	obj.Connect()
	db := obj.Db

	folderMeta, dbErr := db.ListFolder()
	//g :=  folderMeta.Entries[0].Child
	//sharedURL, dbErr := db.GetMediaURL("fantastic_song")
	for _, row := range folderMeta.Entries {
		g := row.Name
		str := fmt.Sprintf(":) sharedURL  got %v, want %v", dbErr, g)
		log.Println(str)
	}
}
func RecoverMe(module string) {
	if r := recover(); r != nil {
		fmt.Println("-->> RECOVERD FROM "+module+" > , ", r)
	}
}

func (obj *HubDropbox) Search() []SearchResult {
	defer func() {
		RecoverMe("HubDropbox Search ")
	}()
	var path, query string
	path = obj.SearchPath
	query = obj.SearchQuery
	obj.Connect()
	db := obj.Db

	folderMeta, _ := db.SearchQuery(path, query)

	fmt.Println("++++++++++++++++++>>>> ", path, query)

	ls := []SearchResult{}
	for _, row := range folderMeta.Matches {
		g := row.Metadata
		//f, _ := json.Marshal(g)
		/*str := fmt.Sprintf("--:)  %v, want %v", dbErr, string(f))
		log.Println(str)*/

		rs := SearchResult{}
		rs.Id = g.UID
		rs.Name = g.Name
		rs.Path = g.Path
		rs.Modified = g.Modified
		rs.Revision = g.Revision
		rs.Size = g.Size
		rs.Tag = g.Tag
		rs.ParentId = g.ParentId

		ls = append(ls, rs)

	}

	obj.SearchResult = ls
	return ls
}
func (obj *HubDropbox) SearchByPrefix() []SearchResult {
	defer func() {
		RecoverMe("HubDropbox Search ")
	}()
	var path, query string
	path = obj.SearchPath
	query = obj.SearchQuery
	obj.Connect()
	db := obj.Db

	folderMeta, _ := db.SearchByPrefix(query)

	fmt.Println("++++++++++++++++++>>>> ", folderMeta, path, query)

	ls := []SearchResult{}
	for _, row := range folderMeta.Matches {
		g := row.Metadata
		//f, _ := json.Marshal(g)
		/*str := fmt.Sprintf("--:)  %v, want %v", dbErr, string(f))
		log.Println(str)*/

		rs := SearchResult{}
		rs.Id = g.UID
		rs.Name = g.Name
		rs.Path = g.Path
		rs.Modified = g.Modified
		rs.Revision = g.Revision
		rs.Size = g.Size
		rs.Tag = g.Tag
		rs.ParentId = g.ParentId

		ls = append(ls, rs)

	}

	obj.SearchResult = ls
	return ls
}

func (obj *HubDropbox) UploadFileCallcentre() Entry {
	obj.Connect()
	db := obj.Db
	src := obj.UploadSrc
	dst := fmt.Sprintf("/callcentre/%s/%s/%s",
		obj.Company,
		obj.UploadBusket,
		obj.UploadFilename)

	myentry, dbErr := db.UploadFile(src, dst, true, "")

	rs := Entry{}
	rs.Tag = myentry.Tag
	rs.Size = myentry.Size
	rs.Revision = myentry.Revision
	rs.Path = myentry.Path
	rs.Name = myentry.Name
	rs.UID = myentry.UID

	obj.UploadResult = rs
	str := fmt.Sprintf(":) UploadFile  got %v,  %v", dbErr, myentry)
	fmt.Println(str)

	return rs
}
func (obj *HubDropbox) Upload(src, dst string) Entry {
	obj.Connect()
	db := obj.Db
	dst = "/" + dst
	myentry, err := db.UploadFile(src, dst, true, "")
	if err != nil {
		log.Println("upload file 	error  response > ", err)
	}
	if myentry == nil {
		log.Println("UploadFile 	EMPTY RESP0NSE > ", myentry)
		return Entry{}
	}
	rs := Entry{}
	rs.Tag = myentry.Tag
	rs.Size = myentry.Size
	rs.Revision = myentry.Revision
	rs.Path = myentry.Path
	rs.Name = myentry.Name
	rs.UID = myentry.UID

	obj.UploadResult = rs

	return rs
}
func (obj *HubDropbox) Download() string {
	CreateFolderUploadIfNotExist(DIR_TEMP_DOWNLOAD)
	obj.Connect()
	obj.DownloadFolder = DIR_TEMP_DOWNLOAD
	db := obj.Db
	src := DIR_TEMP_DOWNLOAD
	err, saveLink := db.DownloadToFile(src, obj.DownloadRev, 0)
	obj.DownloadLink = saveLink
	obj.DownloadError = err
	fmt.Println("Download info ?///////> ", err, saveLink)
	return saveLink
}
func (obj *HubDropbox) MetadataInfo() {
	obj.Connect()
	db := obj.Db
	entry := db.GetMetadataInfoByRev(obj.MetadataRev)
	obj.MetadataResultOne = entry

}
