package dropbox

import "os"

func CreateFolderUploadIfNotExist(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
}
