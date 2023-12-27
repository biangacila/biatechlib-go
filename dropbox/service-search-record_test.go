package dropbox

import (
	"github.com/biangacila/luvungula-go/global"
	"testing"
)

func TestServiceSearchRecord_Run(t *testing.T) {
	hub := ServiceSearchRecord{}
	hub.KeyQuery = "0676865372"
	hub.Params = []string{"C100003"}
	hub.Run()
	global.DisplayObject("LsResult", hub.Results)
}
