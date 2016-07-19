package gsession

import (
	"time"
	"fmt"
	"strconv"
	"github.com/aidonggua/growing/gutils"
)

//会话持续时间
var Session_time time.Duration

func init() {
	//读取配置文件
	//在gopath下查找
	properties := gutils.GetIniProperties("./growing.ini")
	if m, ok := properties["session"]; ok {
		if v, ok := m["session_time"]; ok {
			m, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			Session_time = time.Duration(m) * time.Minute
			fmt.Println(Session_time)
		}
	}
}
