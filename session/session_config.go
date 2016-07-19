package gsession

import (
	"time"
	"gutils"

	"fmt"
)

//会话持续时间
var Session_time time.Duration

func init() {
	//读取配置文件
	//在gopath下查找
	properties := gutils.GetIniProperties("./growing.ini")
	if m, ok := properties["session"]; ok {
		if v, ok := m["session_time"]; ok {
			Session_time = time.Duration(v)
			fmt.Println(Session_time)
		}
	}
}
