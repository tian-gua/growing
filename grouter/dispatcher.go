package grouter

import (
	"net/http"
)

func init() {
	//所有请求先走dispath方法
	http.HandleFunc("/", dispath)
	http.HandleFunc("/static/", http.FileServer(http.Dir("./template")).ServeHTTP)
}

//根据 url 和 method 查找 对应的 处理方法
func dispath(rw http.ResponseWriter, req *http.Request) {
	//获得请求地址和方法
	url := req.RequestURI
	method := req.Method

	//根据http请求的url和Mehtod 执行相应的方法
	switch method {
	case "GET":
		baseController.get(url, rw, req)
	case "POST":
		baseController.post(url, rw, req)
	}

}
