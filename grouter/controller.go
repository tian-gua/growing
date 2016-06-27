package grouter

import "net/http"


//定义控制器的接口,所有接口必须实现 Get Post 方法
type Controller interface {
	Get(rw http.ResponseWriter, req *http.Request)
	Post(rw http.ResponseWriter, req *http.Request)
}



//默认控制器,实现Get和Post方法,用于子控制器继承
//继承后,子控制器可以只实现其中一个方法
type BaseController struct {

}

func (this *BaseController) Get(rw http.ResponseWriter, req *http.Request) {

}

func (this *BaseController) Post(rw http.ResponseWriter, req *http.Request) {

}