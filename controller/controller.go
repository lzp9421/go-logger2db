package controller

import (
	"net/http"
	"encoding/json"
	"time"
	"logger/model"
	"strings"
)

var AppList map[string]int = make(map[string]int);
var Level map[string]int = make(map[string]int);

func init()  {
	AppList["default"] = 0
	AppList["ishangqiang"] = 1

	// level
	Level["UNKNOWN"] = 0;
	Level["DEBUG"] = 1;
	Level["INFO"] = 2;
	Level["NOTICE"] = 3;
	Level["WARNING"] = 4;
	Level["ERROR"] = 5;
	Level["ERROR"] = 5;
	Level["CRITICAL"] = 6;
	Level["ALERT"] = 7;
	Level["EMERGENCY"] = 8;
}

func Log (res http.ResponseWriter, req *http.Request) {
	app, ok := AppList[req.PostFormValue("app")];
	if !ok {
		app = AppList["default"]
	}
	service := req.PostFormValue("service");
	if service == "" {
		service = "unknown"
	}
	level, ok := Level[strings.ToUpper(req.PostFormValue("level"))];
	if !ok {
		level = Level["UNKNOWN"]
	}
	entry := req.PostFormValue("entry")
	message := req.PostFormValue("message")
	trace := req.PostFormValue("trace")
	create_time := time.Now().Unix()
	data := map[string]interface{}{
		"app": app,
		"service": service,
		"level": level,
		"entry": entry,
		"message": message,
		"trace": trace,
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