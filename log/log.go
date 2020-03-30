package log

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

// AdapterConfig adapter配置
type AdapterConfig struct {
	Name       string
	FileName   string
	Level      int
	MaxLines   int64
	MaxSize    int
	Daily      bool
	MaxDays    int
	Rotate     bool
	Perm       string
	RotatePerm string
}

// Config 日志配置
type Config struct {
	FuncCall      bool
	FuncCallDepth int
	Async         bool
	AsyncChan     int64
	Adapters      []*AdapterConfig
}

var cfg Config

func newConfig() {
	cfg = Config{
		FuncCall:      false,
		FuncCallDepth: 2,
		Async:         false,
		AsyncChan:     1e3,
		Adapters:      make([]*AdapterConfig, 0),
	}
}

// Init 初始化
func Init(filename string) (e error) {
	newConfig()

	if filename == "" {
		filename = "./conf/log.json"
	}

	e = LoadJsonToObject(filename, &cfg)
	if e != nil {
		return e
	}

	if cfg.FuncCall {
		logs.EnableFuncCallDepth(cfg.FuncCall)
		logs.SetLogFuncCallDepth(cfg.FuncCallDepth)
	}

	if cfg.Async {
		logs.Async(cfg.AsyncChan)
	}

	for _, adapter := range cfg.Adapters {
		e = logs.SetLogger(adapter.Name, fmt.Sprintf(`{
			"filename":"%v",
			"level":%v,
			"maxlines":%v,
			"maxsize":%v,
			"daily":%v,
			"maxdays":%v,
			"rotate":%v,
			"perm":"%v",
			"rotateperm":"%v"}`,
			adapter.FileName,
			adapter.Level,
			adapter.MaxLines,
			adapter.MaxSize,
			adapter.Daily,
			adapter.MaxDays,
			adapter.Rotate,
			adapter.Perm,
			adapter.RotatePerm,
		))
		if e != nil {
			return e
		}
		logs.Info("adapter:%v, filename:%v", adapter.Name, adapter.FileName)
	}
	return nil
}
