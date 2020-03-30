package handles

import (
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	imagePath = "/home/work/project/file_server/img/"
	serverBase = "http://47.102.208.185:9653"
)

var (
	atomicCount int64
)

// HandleNcPostUploadPicture 处理上传图片逻辑
func HandleNcPostUploadPicture(writer http.ResponseWriter, request *http.Request) {
	// 解决跨域问题
	writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	e := request.ParseForm()
	if e != nil {
		logs.Error("request.ParseForm failed, error = %v", e.Error())
		io.WriteString(writer, "parse form failed.")
		return
	}
	var pictureDesc string
	if request.Method == "GET" {
		pictureDesc = request.FormValue("picture_name")
	} else if request.Method == "POST" {
		body, e := ioutil.ReadAll(request.Body)
		if e != nil {
			logs.Error("ioutil.ReadAll failed, error = %v", e.Error())
			io.WriteString(writer, "read body failed")
			return
		}
		defer request.Body.Close()
		logs.Debug("body = %v", string(body))

	} else {
		io.WriteString(writer, "unsupported method")
		return
	}

	if pictureDesc == "" {
		io.WriteString(writer, "picture_name is nil")
		return
	}
	logs.Debug("pictureDesc = %v", pictureDesc)
	count := atomic.AddInt64(&atomicCount, 1)
	fileName := time.Now().Format("20060102150405")
	fileName += "_" + strconv.FormatInt(count, 10)
	filePath := imagePath + fileName
	pictureURL := serverBase + "/get_picture?id=" + fileName
	logs.Debug("filePath = %v, pictureUrl = %v", filePath, pictureURL)

	file, e := os.Create(filePath)
	if e != nil {
		logs.Error("os.Create failed, filePath = %v", filePath)
		io.WriteString(writer, "create file failed")
		return
	}

	_, e = file.WriteString(pictureDesc)
	if e != nil {
		logs.Error("file.WriteString failed, error = %v", e.Error())
		io.WriteString(writer, "write file failed")
		return
	}

	io.WriteString(writer, pictureURL)
}

// HandleNcGetGetPicture 获取图片
func HandleNcGetGetPicture(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	e := request.ParseForm()
	if e != nil {
		logs.Error("request.ParseForm failed, error = %v", e.Error())
		io.WriteString(writer, "parse form failed")
		return
	}
	if request.Method == "GET" {
		fileName := request.FormValue("id")
		if fileName == "" {
			io.WriteString(writer, "id is unset")
			return
		}
		filePath := imagePath + fileName
		logs.Debug("filePath = %v", filePath)

		file, e := os.OpenFile(filePath, os.O_RDONLY, 0666)
		if e != nil {
			logs.Error("os.OpenFile failed, error = %v", e.Error())
			io.WriteString(writer, "open file failed")
			return
		}
		defer file.Close()
		b, e := ioutil.ReadAll(file)
		if e != nil {
			logs.Error("ioutil.ReadAll failed, error = %v", e.Error())
			io.WriteString(writer, "read file failed")
			return
		}
		logs.Debug("b = %v", string(b))
		writer.Write(b)
	} else {
		io.WriteString(writer, "unsupported method")
		return
	}
}
