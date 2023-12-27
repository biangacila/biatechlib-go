package dropbox

import (
	"github.com/biangacila/luvungula-go/global"
	"testing"
)

func TestHubDropbox_Search(t *testing.T) {
	var hub HubDropbox
	hub.Company = ""
	hub.SearchQuery = "27832990519" //"27832990519"
	hub.DROPBOX_CLIENTID = ""       //constant.DROPBOX_CLIENTID
	hub.DROPBOX_TOKEN = ""          // constant.DROPBOX_TOKEN
	hub.DROPBOX_SECRET = ""         // constant.DROPBOX_SECRET
	ls := hub.Search()
	global.DisplayObject("LsResult", ls)

}
