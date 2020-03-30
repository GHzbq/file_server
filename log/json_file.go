package log

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// LoadJsonToObject 读取json文件
func LoadJsonToObject(filename string, t interface{}) error {

	buf, e := loadFile(filename)

	if buf == nil {
		return e
	}

	if e != nil {
		return e
	}

	err := json.Unmarshal(buf, &t)

	if err != nil {
		return err
	}

	return nil
}


// loadFile 读取文件
func loadFile(filename string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return buf, nil
}