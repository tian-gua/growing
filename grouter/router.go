package grouter

import "strings"

//将url 和 处理器 映射到一起
func Route(requestMapping string, h handler, method ...string) {

	//如果不指定方法,get和post请求都进行处理
	if len(method) == 0 {
		baseController.addGet(requestMapping, h)
		baseController.addPost(requestMapping, h)
	}else {
		m := method[0]
		if strings.ToUpper(m) == "POST" {
			baseController.addPost(requestMapping, h)

		} else if strings.ToUpper(m) == "GET" {
			baseController.addGet(requestMapping, h)
		}
	}

}


