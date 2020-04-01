package handles

import (
	"github.com/astaxie/beego/logs"
	imgType "github.com/shamsher31/goimgtype"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	//serverBase = "http://47.102.208.185:9653"
	serverBase = "http://127.0.0.1:9653"
	imagePath = "/home/zhangbaqing/all/project/file_server/img/"
)

var (
	atomicCount int64
)

func init() {
	// 存储图片的目录是否存在，不存在就创建，存在不做任何处理
	ok, e := PathExists(imagePath)
	if e != nil {
		logs.Error("unknown error, error = %v", e.Error())
		return
	}
	if ok == false {
		// 目录不存在，创建之
		e := os.MkdirAll(imagePath, os.ModePerm)
		if e != nil {
			logs.Error("os.MkdirAll failed, error = %v", e.Error())
			return
		}
		logs.Debug("os.MkdirAll succeed, path = %v", imagePath)
	}
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// HandleNcPostUploadPicture 处理上传图片逻辑
func HandleNcPostUploadPicture(writer http.ResponseWriter, request *http.Request) {
	// 解决跨域问题
	writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	logs.Debug("get a new request, request.Method = %v", request.Method)
	e := request.ParseForm()
	if e != nil {
		logs.Error("request.ParseForm failed, error = %v", e.Error())
		http.Error(writer, "parse form failed", http.StatusInternalServerError)
		return
	}
	logs.Debug("request.ParseForm succeed")
	if request.Method == "POST" {
		f, h, e := request.FormFile("filename")
		if e != nil {
			logs.Error("request.FormFile failed, error = %v", e.Error())
			http.Error(writer, "request.FormFile failed", http.StatusInternalServerError) // 500
			return
		}
		defer f.Close()
		logs.Debug("request.FormFile succeed, h.Filename = %v", h.Filename)
		// 限定文件格式
		logs.Debug("ready to get fileType")
		fileType, e := imgType.Get(h.Filename)
		if e != nil {
			logs.Error("not image format")
			http.Error(writer, "not image format", http.StatusBadRequest)
			return
		}
		logs.Debug("imgType.Get succeed, fileType = %v", fileType)

		count := atomic.AddInt64(&atomicCount, 1)
		fileName := time.Now().Format("20060102150405")
		fileName += "_" + strconv.FormatInt(count, 10) + fileType
		filePath := imagePath + fileName
		fileURL := serverBase + "/get_picture?id=" + fileName
		logs.Debug("filePath = %v, fileURL = %v", filePath, fileURL)

		file, e := os.Create(filePath)
		if e != nil {
			logs.Error("os.Create failed, error = %v", e.Error())
			http.Error(writer, "create file failed", http.StatusInternalServerError)
			return
		} else {
			logs.Debug("os.Create succeed, filePath = %v", filePath)
		}
		defer file.Close()
		_, e = io.Copy(file, f)
		if e != nil {
			logs.Error("io.Copy failed, error = %v", e.Error())
			http.Error(writer, "write file failed", http.StatusInternalServerError)
			return
		}
		_, e = io.WriteString(writer, fileURL)
		if e != nil {
			http.Error(writer, "write fileURL failed", http.StatusInternalServerError)
			return
		}
		http.Redirect(writer, request, fileURL, http.StatusFound) // 302
		return
	} else {
		logs.Debug("upload method is unsupported")
		http.Error(writer, "upload method is unsupported", http.StatusBadRequest) // 400
		return
	}
}

// HandleNcGetGetPicture 获取图片
func HandleNcGetGetPicture(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	logs.Debug("get a new request, method = %v", request.Method)
	e := request.ParseForm()
	if e != nil {
		logs.Error("request.ParseForm failed, error = %v", e.Error())
		http.Error(writer, "parse form failed", http.StatusInternalServerError)
		return
	}
	logs.Debug("request.ParseForm succeed")
	if request.Method == "GET" {
		fileName := request.FormValue("id")
		if fileName == "" {
			logs.Debug("id is unset")
			http.Error(writer, "id is unset", http.StatusBadRequest)
			return
		}
		filePath := imagePath + fileName
		logs.Debug("filePath = %v", filePath)

		ok, e := PathExists(filePath)
		if e != nil {
			logs.Error("unknown error, error = %v", e.Error())
			http.Error(writer, "unknown error", http.StatusInternalServerError)
			return
		}
		if ok == false {
			logs.Debug("file not exists, filePath = %v", filePath)
			http.NotFound(writer, request)
		}

		writer.Header().Set("Content-Type", "image")
		http.ServeFile(writer, request, imagePath)
	} else {
		logs.Debug("unsupported method, request.Method = %v", request.Method)
		http.Error(writer, "unsupported method", http.StatusBadRequest)
		return
	}
}
