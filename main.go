package main

import (
	"file_server/handles"
	"file_server/log"
	"github.com/astaxie/beego/logs"
	"net/http"
)

func init() {
	e := log.Init("conf/log.json")
	if e == nil {
		logs.Info("init log successfully.")
	} else {
		logs.Error("init log failed, error = %v", e.Error())
	}
}

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/upload_picture", http.HandlerFunc(handles.HandleNcPostUploadPicture))
	serveMux.Handle("/get_picture", http.HandlerFunc(handles.HandleNcGetGetPicture))
	e := http.ListenAndServe(":9653", serveMux)
	if e != nil {
		logs.Error("http.ListenAndServe failed, error = %v", e.Error())
		return
	}
}
