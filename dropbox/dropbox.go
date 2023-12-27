package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	//"io"
	"io"
	"log"
	"os"
)

// Dropbox client
type Dropbox struct {
	Debug         bool   // bool to dump request/response data
	Locale        string // Locale sent to the API to translate/format messages
	Token         *Token
	mediaURL      string
	listFolderURL string
	searchURL     string
	oAuth2Handler *OAuth2Handler
}

// NewDropbox returns a new Dropbox instance
func NewDropbox() *Dropbox {
	return &Dropbox{
		Locale:        "en",
		mediaURL:      MEDIA_URL,
		listFolderURL: LIST_FOLDER_URL,
		searchURL:     SEARCH_URL,
		oAuth2Handler: newOAuth2Handler(),
	}
}

// SetAppInfo sets app_key & app_secret from your Dropbox app
func (db *Dropbox) SetAppInfo(appKey string, appSecret string, redirectURL string) {
	db.oAuth2Handler.setAppKeys(appKey, appSecret, redirectURL)
}

// SetAccessToken sets access token
func (db *Dropbox) SetAccessToken(accessToken string) {
	db.Token = &Token{
		Token: accessToken,
	}
}

func (db *Dropbox) GetAccessToken() string {
	return db.Token.Token
}

func (db *Dropbox) GetAuthURL() string {
	return db.oAuth2Handler.authCodeURL()
}

func (db *Dropbox) ExchangeToken(code string) (*Token, error) {
	return db.oAuth2Handler.tokenExchange(code)
}

// Shares a file for streaming (direct access)
func (db *Dropbox) GetMediaURL(file string) (*SharedURL, *DropboxError) {

	client := &http.Client{}

	data := mediaParameters{
		Locale: db.Locale,
	}

	encoded, _ := json.Marshal(data)

	request, _ := http.NewRequest(POST, db.mediaURL+file, bytes.NewBuffer(encoded))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	if db.Debug {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Println(string(dump))
	}

	response, err := client.Do(request)

	if err != nil {
		return nil, &DropboxError{
			StatusCode:   http.StatusBadRequest,
			ErrorSummary: err.Error(),
		}
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 401:
		return nil, &DropboxError{
			StatusCode:   http.StatusUnauthorized,
			ErrorSummary: INVALID_ACCESS_TOKEN,
		}
	default:
		dumpData, _ := ioutil.ReadAll(response.Body)

		if db.Debug {
			fmt.Printf("%s\n", string(dumpData))
		}

		var mediaURL SharedURL

		err := json.Unmarshal(dumpData, &mediaURL)

		if err != nil {
			return nil, &DropboxError{
				StatusCode:   http.StatusServiceUnavailable,
				ErrorSummary: err.Error(),
			}
		}

		return &mediaURL, nil
	}
}

func (db *Dropbox) ListFolder() (*Folder, *DropboxError) {
	client := &http.Client{}

	data := listFolderParameters{
		Path:             "",
		Recursive:        true,
		IncludeMediaInfo: false,
		IncludeDeleted:   false,
	}

	encoded, _ := json.Marshal(data)

	req, _ := http.NewRequest(POST, db.listFolderURL, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	if db.Debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, &DropboxError{
			StatusCode:   http.StatusBadRequest,
			ErrorSummary: err.Error(),
		}
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 401:
		return nil, &DropboxError{
			StatusCode:   http.StatusUnauthorized,
			ErrorSummary: INVALID_ACCESS_TOKEN,
		}
	default:
		dumpData, _ := ioutil.ReadAll(res.Body)

		if db.Debug {
			fmt.Printf("%s\n", string(dumpData))
		}

		var metadata Folder

		err := json.Unmarshal(dumpData, &metadata)

		if err != nil {
			return nil, &DropboxError{
				StatusCode:   http.StatusServiceUnavailable,
				ErrorSummary: err.Error(),
			}
		}

		return &metadata, nil
	}
}
func (db *Dropbox) SearchByPrefix(searchKey string) (*Search, *DropboxError) {
	client := &http.Client{}

	data := searchParameters{
		Path:       "",
		Query:      searchKey,
		Mode:       &searchMode{Tag: "filename"},
		MaxResults: 1000,
	}

	encoded, _ := json.Marshal(data)

	req, _ := http.NewRequest(POST, db.searchURL, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	if db.Debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, &DropboxError{
			StatusCode:   http.StatusBadRequest,
			ErrorSummary: err.Error(),
		}
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 401:
		return nil, &DropboxError{
			StatusCode:   http.StatusUnauthorized,
			ErrorSummary: INVALID_ACCESS_TOKEN,
		}
	default:
		dumpData, _ := ioutil.ReadAll(res.Body)

		if db.Debug {
			fmt.Printf("%s\n", string(dumpData))
		}

		var metadata Search

		err := json.Unmarshal(dumpData, &metadata)

		if err != nil {
			return nil, &DropboxError{
				StatusCode:   http.StatusServiceUnavailable,
				ErrorSummary: err.Error(),
			}
		}

		return &metadata, nil
	}
}
func (db *Dropbox) SearchMusic() (*Search, *DropboxError) {
	client := &http.Client{}

	data := searchParameters{
		Path:       "",
		Query:      ".mp3",
		Mode:       &searchMode{Tag: "filename"},
		MaxResults: 1000,
	}

	encoded, _ := json.Marshal(data)

	req, _ := http.NewRequest(POST, db.searchURL, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	if db.Debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, &DropboxError{
			StatusCode:   http.StatusBadRequest,
			ErrorSummary: err.Error(),
		}
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 401:
		return nil, &DropboxError{
			StatusCode:   http.StatusUnauthorized,
			ErrorSummary: INVALID_ACCESS_TOKEN,
		}
	default:
		dumpData, _ := ioutil.ReadAll(res.Body)

		if db.Debug {
			fmt.Printf("%s\n", string(dumpData))
		}

		var metadata Search

		err := json.Unmarshal(dumpData, &metadata)

		if err != nil {
			return nil, &DropboxError{
				StatusCode:   http.StatusServiceUnavailable,
				ErrorSummary: err.Error(),
			}
		}

		return &metadata, nil
	}
}
func (db *Dropbox) SearchQuery(path, query string) (*Search, *DropboxError) {
	client := &http.Client{}

	data := searchParameters{
		Path:       path,
		Query:      query,
		Mode:       &searchMode{Tag: "filename"},
		MaxResults: 1000,
	}

	encoded, _ := json.Marshal(data) /**/

	req, _ := http.NewRequest(POST, db.searchURL, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	if db.Debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, &DropboxError{
			StatusCode:   http.StatusBadRequest,
			ErrorSummary: err.Error(),
		}
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 401:
		return nil, &DropboxError{
			StatusCode:   http.StatusUnauthorized,
			ErrorSummary: INVALID_ACCESS_TOKEN,
		}
	default:
		dumpData, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("[dumpData] --> %s\n", string(dumpData))

		if db.Debug {
			fmt.Printf("%s\n", string(dumpData))
		}

		var metadata Search

		err := json.Unmarshal(dumpData, &metadata)

		if err != nil {
			return nil, &DropboxError{
				StatusCode:   http.StatusServiceUnavailable,
				ErrorSummary: err.Error(),
			}
		}

		return &metadata, nil
	}
}

func (db *Dropbox) UploadFile(src, dst string, overwrite bool, parentRev string) (*Entry, error) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["path"] = dst
	data["mode"] = "overwrite"
	data["autorename"] = false
	data["mute"] = true

	encoded, _ := json.Marshal(data)

	var fd *os.File
	fd, _ = os.Open(src)

	request, _ := http.NewRequest(POST, API_CONTENT_URL, fd)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))
	request.Header.Set("Dropbox-API-Arg", string(encoded))
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("data-binary", src)

	if db.Debug {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Println(string(dump))
	}

	var rv Entry
	var body []byte
	var response *http.Response
	var err error

	if response, err = client.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if body, err = getResponse(response); err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &rv)
	return &rv, err

}

func GetFileContentType(src string) string {

	f, _ := os.Stat(src)
	str, err := json.Marshal(f)

	log.Println("[FILE INTO] > ", err, string(str), " > ", src)
	return ""
}

// Download requests the file located at src, the specific revision may be given.
// offset is used in case the download was interrupted.
// A io.ReadCloser and the file size is returned.
func (db *Dropbox) DownloadToFile(dst, rev string, offset int64) (error, string) {
	entry := db.GetMetadataInfoByRev(rev)
	filename := dst + "" + entry.Name

	client := &http.Client{}
	data := make(map[string]interface{})
	data["path"] = "rev:" + rev

	st := fmt.Sprintf(`{"path": "rev:%s"}`, rev)

	//encoded, _ := json.Marshal(data)

	request, _ := http.NewRequest(POST, API_CONTENT_DOWNLOAD_URL, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))
	request.Header.Set("Dropbox-API-Arg", st)
	request.Header.Set("Content-Type", "application/octet-stream")

	if db.Debug {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Println(string(dump))
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("ERROR DOWNLOAD > ", err)
		//return nil, err
	}
	defer response.Body.Close()

	//TODO LET SAVE NO THE RESPONSE TO THE FILE
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusPartialContent {
		var input io.ReadCloser
		var fd *os.File
		var err error

		if fd, err = os.Create(filename); err != nil {
			log.Println("err os.Create > ", err, " > ", entry.Name)
			return err, entry.Name
		}
		defer fd.Close()

		input = response.Body

		x, err := io.Copy(fd, input)
		fmt.Println("+++++++ERROR > ", err, x, response.StatusCode)

		if err != nil {

			os.Remove(dst)
		}

		input.Close()

	}

	return nil, filename
}

func (db *Dropbox) GetMetadataInfoByRev(rev string) Entry {
	client := &http.Client{}
	myUrl := "https://api.dropboxapi.com/2/files/get_metadata"
	myMap := make(map[string]interface{})
	myMap["path"] = fmt.Sprintf("rev:%s", rev)
	myMap["include_media_info"] = false
	myMap["include_deleted"] = false
	myMap["include_has_explicit_shared_members"] = false

	encoded, _ := json.Marshal(myMap)

	//st :=fmt.Sprintf(`{"path": "rev:%s","include_media_info": false,"include_deleted": false,"include_has_explicit_shared_members": false}`,rev)
	request, _ := http.NewRequest(POST, myUrl, bytes.NewBuffer(encoded))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))
	//request.Header.Set("Content-Type", "application/json")
	if db.Debug {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Println(string(dump))
	}

	response, err := client.Do(request)

	if err != nil {
		log.Println("ERROR METADATA DOWNLOAD client.Do > ", err, string(encoded))
		return Entry{}
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)

	j := Entry{}
	json.Unmarshal(b, &j)

	fmt.Println("====]--> ", j.Name, " > ", j.Size)

	return j

}

func (db *Dropbox) client() *http.Client {
	client := &http.Client{}
	return client
}

/*
#-=---------------------------
#-------- OTHER FUNCTION -----
#-----------------------------
*/

// Format of reply when http error code is not 200.
// Format may be:
// {"error": "reason"}
// {"error": {"param": "reason"}}
type requestError struct {
	Error interface{} `json:"error"` // Description of this error.
}

// Error - all errors generated by HTTP transactions are of this type.
// Other error may be passed on from library functions though.
type Error struct {
	StatusCode int // HTTP status code
	Text       string
}

// Error satisfy the error interface.
func (e *Error) Error() string {
	return e.Text
}

// newError make a new error from a string.
func newError(StatusCode int, Text string) *Error {
	return &Error{
		StatusCode: StatusCode,
		Text:       Text,
	}
}

// newErrorf makes a new error from sprintf parameters.
func newErrorf(StatusCode int, Text string, Parameters ...interface{}) *Error {
	return newError(StatusCode, fmt.Sprintf(Text, Parameters...))
}

// urlEncode encodes s for url
func urlEncode(s string) string {
	// Would like to call url.escape(value, encodePath) here
	encoded := url.QueryEscape(s)
	encoded = strings.Replace(encoded, "+", "%20", -1)
	return encoded
}
func getResponse(r *http.Response) ([]byte, error) {
	var e requestError
	var b []byte
	var err error

	if b, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	}
	if r.StatusCode == http.StatusOK {
		return b, nil
	}
	if err = json.Unmarshal(b, &e); err == nil {
		switch v := e.Error.(type) {
		case string:
			return nil, newErrorf(r.StatusCode, "%s", v)
		case map[string]interface{}:
			for param, reason := range v {
				if reasonstr, ok := reason.(string); ok {
					return nil, newErrorf(r.StatusCode, "%s: %s", param, reasonstr)
				}
			}
			return nil, newErrorf(r.StatusCode, "wrong parameter")
		}
	}
	return nil, newErrorf(r.StatusCode, "unexpected HTTP status code %d", r.StatusCode)
}
