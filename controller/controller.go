package controller

import (
	"net/http"
	"encoding/json"
	"time"
	"logger/model"
)

var AppList map[string]int = make(map[string]int);
var ServiceList map[string]int = make(map[string]int);
var Level map[string]int = make(map[string]int);

func init()  {
	AppList["ishangqiang"] = 1

	ServiceList["mysql"] = 1

	// level
	Level["DEBUG"] = 1;
	Level["WARNING"] = 2;
}

func Log (res http.ResponseWriter, req *http.Request) {
	app, ok := AppList[req.PostFormValue("app")];
	if !ok {
		res.Write(response(1, "没有指定应用或应用不存在", make(map[string]interface{})))
		return
	}
	service, ok := ServiceList[req.PostFormValue("service")];
	if !ok {
		res.Write(response(2, "没有指定业务或业务不存在", make(map[string]interface{})))
		return
	}
	level, ok := Level[req.PostFormValue("level")];
	if !ok {
		res.Write(response(3, "级别不存在", make(map[string]interface{})))
		return
	}
	entry := req.PostFormValue("entry")
	title := req.PostFormValue("title")
	body := req.PostFormValue("body")
	create_time := time.Now().Unix()
	data := map[string]interface{}{
		"app": app,
		"service": service,
		"level": level,
		"entry": entry,
		"title": title,
		"body": body,
		"create_time": create_time,
	}
	go func() {
		mdlLog := model.Log{}
		mdlLog.Insert(data)
	}()
	res.Write(response(0, "success", make(map[string]interface{})))
	return
}

func response (code int, message string, data map[string]interface{}) []byte {
	result := map[string]interface{}{
		"code": code,
		"message": message,
		"data": data,
	}
	json_data, err := json.Marshal(result)
	if err != nil {
		json_data = nil
	}
	return json_data
}