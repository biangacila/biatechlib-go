package pdfs

import (
	"fmt"
	"testing"
)

func TestGeneratePdfFromHtml_Run(t *testing.T) {
	var hub GeneratePdfFromHtml
	hub.HtmlBody = fmt.Sprintf("<html><body><table><tr><th>Rep Type</th></tr><tr><td>Monthly</td></tr></table></body></body></html>")
	hub.AppName = "pmis"
	hub.Org = "C100010"
	hub.Module = "WoodChipping"
	hub.Ref = "cemex-oct-2023"
	hub.ActionSaveToDropbox = true

	/**
	TODO clear before commit
	*/
	hub.DropboxClientId = "" // provide with your detail
	hub.DropboxSecret = ""   // provide with your detail
	hub.DropboxToken = ""    // provide with your detail

	hub.Run()

	fmt.Println("Download link > ", hub.Result.LinkDownload)

}
