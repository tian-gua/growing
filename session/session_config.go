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
	properties, err := gutils.GetIniProperties("../growing.ini")
	if err == nil {
		if m, ok := properties["session"]; ok {
			if v, ok := m["session_time"]; ok {
				m, err := strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				Session_time = time.Duration(m) * time.Minute
				fmt.Printf("session闲置时间为:%s\n", Session_time)
			}
			return
		}
	}
	Session_time = 30 * time.Minute
	fmt.Printf("session闲置时间为默认时间:%s\n", Session_time)

}
