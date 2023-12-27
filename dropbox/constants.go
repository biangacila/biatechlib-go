package dropbox

import "time"

const (
	// PollMinTimeout is the minimum timeout for longpoll.
	PollMinTimeout = 30
	// PollMaxTimeout is the maximum timeout for longpoll.
	PollMaxTimeout = 480
	// DefaultChunkSize is the maximum size of a file sendable using files_put.
	DefaultChunkSize = 4 * 1024 * 1024
	// MaxPutFileSize is the maximum size of a file sendable using files_put.
	MaxPutFileSize = 150 * 1024 * 1024
	// MetadataLimitMax is the maximum number of entries returned by metadata.
	MetadataLimitMax = 25000
	// MetadataLimitDefault is the default number of entries returned by metadata.
	MetadataLimitDefault = 10000
	// RevisionsLimitMax is the maximum number of revisions returned by revisions.
	RevisionsLimitMax = 1000
	// RevisionsLimitDefault is the default number of revisions returned by revisions.
	RevisionsLimitDefault = 10
	// SearchLimitMax is the maximum number of entries returned by search.
	SearchLimitMax = 1000
	// SearchLimitDefault is the default number of entries returned by search.
	SearchLimitDefault = 1000
	// DateFormat is the format to use when decoding a time.
	DateFormat = time.RFC1123Z
	//this is use to tempory download file
	DIR_TEMP_DOWNLOAD = "./tempdownload/"
)

const (
	POST                     = "POST"
	INVALID_ACCESS_TOKEN     = "invalid_access_token"
	MEDIA_URL                = "https://api.dropboxapi.com/1/media/auto/"
	LIST_FOLDER_URL          = "https://api.dropboxapi.com/2/files/list_folder"
	SEARCH_URL               = "https://api.dropboxapi.com/2/files/search"
	API_CONTENT_URL          = "https://content.dropboxapi.com/2/files/upload"
	API_CONTENT_DOWNLOAD_URL = "https://content.dropboxapi.com/2/files/download"
	ROOT_DIRECTORY           = "auto"
	API_NOTIFY_URL           = "https://api-notify.dropbox.com/2"
)
