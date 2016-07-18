package gutils

import (
	"os"
	"bufio"
	"fmt"
	"regexp"
	"strings"
)


//获取ini文件的配置
//[default]
//key=value
func GetIniProperties(path string) map[string]map[string]string {
	properties := make(map[string]map[string]string)
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	var title string = "default"
	for {
		if line, err := reader.ReadString('\n'); err == nil || len(line) != 0{
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
			fmt.Println(err)
			break
		}
	}
	return properties
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
	kv := strings.Split(strings.Trim(line,"\n"), "=")
	if len(kv) == 2 {
		ok = true
	}
	return kv, ok
}