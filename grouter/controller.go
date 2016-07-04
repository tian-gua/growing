package grouter

import (
	"net/http"
	"fmt"
)

type handler func(rw http.ResponseWriter, req *http.Request)


//定义控制器的接口,所有接口必须实现 Get Post 方法
type controller struct {
	getHandlers  map[string]handler
	postHandlers map[string]handler
	getnofund    handler
}
//初始化控制器
var baseController = &controller{getHandlers: make(map[string]handler), postHandlers:make(map[string]handler)}

//添加get处理器
func (c *controller) addGet(pattern string, h handler) {
	c.getHandlers[pattern] = h
}
//添加post处理器
func (c *controller) addPost(pattern string, h handler) {
	c.postHandlers[pattern] = h
}

//请求get方法
func (c *controller) get(pattern string, rw http.ResponseWriter, req *http.Request) {

	if h, ok := c.getHandlers[pattern]; ok {
		h(rw, req)
	} else {
		fmt.Println("未找到[" + pattern + "]对应的GET处理器!")
		rw.WriteHeader(http.StatusNotFound)

	}
}
//请求post方法
func (c *controller)  post(pattern string, rw http.ResponseWriter, req *http.Request) {
	if h, ok := c.postHandlers[pattern]; ok {
		h(rw, req)
	} else {
		fmt.Println("未找到[" + pattern + "]对应的POST处理器!")
		rw.WriteHeader(http.StatusNotFound)
	}
}