package dropbox

// not exported

type listFolderParameters struct {
	Path             string `json:"path"`
	Recursive        bool   `json:"recursive"`
	IncludeMediaInfo bool   `json:"include_media_info"`
	IncludeDeleted   bool   `json:"include_deleted"`
}

type mediaParameters struct {
	Locale string `json:"locale"`
}

type searchMode struct {
	Tag string `json:".tag"`
}

type searchParameters struct {
	Path       string      `json:"path"`
	Query      string      `json:"query"`
	Mode       *searchMode `json:"mode"`
	MaxResults uint64      `json:"max_results"`
}

// exported

type Token struct {
	UID   string  `json:"uid"`
	Token string  `json:"access_token"`
	Error *string `json:"error"`
}

type Folder struct {
	Entries []*Entry `json:"entries"`
	Cursor  string   `json:"cursor"`
	HasMore bool     `json:"has_more"`
}

type SharedURL struct {
	URL     string `json:"url"`
	Expires string `json:"expires"`
}

type DropboxError struct {
	ErrorSummary string `json:"error_summary"`
	StatusCode   int
}

type MatchType struct {
	Tag string `json:".tag"`
}

type Metadata struct {
	UID      string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Path     string `json:"path_display,omitempty"`
	Revision string `json:"rev,omitempty"`
	Size     int    `json:"size,omitempty"`
	Tag      string `json:".tag"`
	Modified string `json:"client_modified"`
	ParentId string `json:"parent_shared_folder_id"`
}

type Match struct {
	MatchType *MatchType `json:"match_type"`
	Metadata  *Metadata  `json:"metadata"`
}

type Search struct {
	Matches []*Match `json:"matches"`
	Start   int      `json:"start"`
	More    bool     `json:"more"`
}
