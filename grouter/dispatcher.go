package grouter

import (
	"net/http"
)


//基础控制d器,用于子控制器继承
type dispatcher  struct {
	mappings map[string]Controller
}

//初始化dispatcher对象
var dpc = &dispatcher{make(map[string]Controller)}


//根据 url 和 method 查找 对应的 处理方法
func (dp *dispatcher) dispath(rw http.ResponseWriter, req *http.Request) {
	//获得请求地址
	url := req.RequestURI
	//查找相应的控制器
	if c, ok := dpc.mappings[url]; ok {
		//根据http请求的Mehtod 执行相应的方法
		switch req.Method  {
		case "GET":
			c.Get(rw, req)
		case "POST":
			c.Post(rw, req)
		}
	} else {
		//如果找不到映射,则返回一个404状态码
		rw.WriteHeader(http.StatusNotFound)
	}

}

