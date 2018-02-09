package utils

import (
	"io/ioutil"
	"os"
	"strings"
)

/*
IsDirExist function
判断某个目录是否存在
*/
func IsDirExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

/*
ArrayContainer function
检测数组里面是否包含某个字符串
*/
func ArrayContainer(array []string, found string) bool {
	for _, v := range array {
		if v == found {
			return true
		}
	}

	return false
}

// ListDir 列出所有的目录
func ListDir(path string, suffix string) (files []string, err error) {
	if !IsDirExist(path) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			return nil, err
		}
	}
	files = []string{}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	suffix = strings.ToUpper(suffix)

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, fi.Name())
		}
	}

	return files, nil
}
