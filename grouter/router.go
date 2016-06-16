package grouter

import "net/http"

//将url 和 控制器 映射到一起
func Route(requestMapping string, c Controller) {
	//将控制器和url通过map关联
	dpc.mappings[requestMapping] = c

}

func init() {
	//所有请求先走dispath方法
	http.HandleFunc("/", dpc.dispath)
}

