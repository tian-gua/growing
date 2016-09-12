package gutils

import (
	"os"
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"io"
	"path/filepath"
	"runtime"
)


//获取ini文件的配置
//[default]
//key=value
func GetIniProperties(path string) (map[string]map[string]string, error) {
	properties := make(map[string]map[string]string)
	//获取调用此函数的文件的绝对路径
	_, filename, _, _ := runtime.Caller(1)
	//拼接ini文件的路径
	file, err := os.Open(filepath.Join(filepath.Dir(filename), path))
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	var title string = "default"
	for {
		if line, err := reader.ReadString('\n'); err == nil || len(line) != 0 {
			//如果是标题
			if t, ok := isTitle(line); ok {
				title = t
				properties[title] = make(map[string]string)
			} else if kv, ok := isKV(line); ok {
				properties[title][kv[0]] = kv[1]
			} else {
				fmt.Println("未识别的配置:" + line)
			}
		} else {
			if err != io.EOF {
				return nil, err
			}
			break
		}
	}
	return properties, nil
}


//是否为标题
//[xxx]
func isTitle(line string) (string, bool) {
	var ok bool = false
	reg, _ := regexp.Compile("\\[\\w*\\]")
	result := reg.FindString(line)
	if len(result) != 0 {
		ok = true
	}
	return strings.TrimRight(strings.TrimLeft(result, "["), "]"), ok
}

//是否是属性
//k=v
func isKV(line string) ([]string, bool) {
	var ok bool = false
	kv := strings.Split(strings.Trim(line, "\n"), "=")
	if len(kv) == 2 {
		ok = true
	}
	return kv, ok
}