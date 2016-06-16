package grouter

import "net/http"


//定义控制器的接口,所有接口必须实现 Get Post 方法
type Controller interface {
	Get(rw http.ResponseWriter, req *http.Request)
	Post(rw http.ResponseWriter, req *http.Request)
}
