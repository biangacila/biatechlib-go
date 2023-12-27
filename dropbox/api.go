package dropbox

import (
	"encoding/json"
	"github.com/biangacila/luvungula-go/global"
	"net/http"
)

func WsLeadSearch(w http.ResponseWriter, r *http.Request) {
	_, dataString := global.GetPostedDataMapAndString(r)
	hub := ServiceSearchRecord{}
	_ = json.Unmarshal([]byte(dataString), &hub)
	hub.Run()
	ls := hub.Results
	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = hub
	my["RESULT"] = ls
	global.PublishToReact(w, r, my, 200)

}
