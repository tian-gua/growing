package grouter

import (
	"net/http"
	"fmt"
	"regexp"
	"strings"
	"reflect"
	"strconv"
)

type handler interface{}


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

	if h, ok := c.getHandlers[TrimParameter(pattern)]; ok {
		if checkHandler(h) {
			do(rw, req, h)
		}
	} else {
		fmt.Println("未找到[" + pattern + "]对应的GET处理器!")
		rw.WriteHeader(http.StatusNotFound)

	}
}
//请求post方法
func (c *controller)  post(pattern string, rw http.ResponseWriter, req *http.Request) {
	if h, ok := c.postHandlers[TrimParameter(pattern)]; ok {
		if checkHandler(h) {
			do(rw, req, h)
		}
	} else {
		fmt.Println("未找到[" + pattern + "]对应的POST处理器!")
		rw.WriteHeader(http.StatusNotFound)
	}
}
//trim掉&后面的参数
func TrimParameter(url string) string {
	reg, err := regexp.Compile("\\?.*")
	if err != nil {
		fmt.Println(err)
		return url
	}
	return strings.TrimRight(url, reg.FindString(url))
}


//检查处理器是不是func
func checkHandler(h handler) bool {
	k := reflect.TypeOf(h).Kind()
	if reflect.Func == k {
		return true
	}
	return false
}

func do(rw http.ResponseWriter, req *http.Request, h handler) {
	funcT := reflect.TypeOf(h)
	funcV := reflect.ValueOf(h)
	numin := funcT.NumIn()
	vs := new([]reflect.Value)
	for i := 0; i < numin; i++ {
		if funcT.In(i).Kind() == reflect.Ptr {
			if "Request" == funcT.In(i).Elem().Name() {
				*vs = append(*vs, reflect.ValueOf(req))
			}
		} else {
			if "ResponseWriter" == funcT.In(i).Name() {
				*vs = append(*vs, reflect.ValueOf(rw))
			} else if reflect.Struct == funcT.In(i).Kind() {
				v := reflect.New(funcT.In(i)).Elem()
				t := reflect.TypeOf(v.Interface())
				//遍历属性,并赋值
				//值从表单里取,根据属性名字对应
				for j := 0; j < t.NumField(); j++ {
					sf := t.Field(j)
					sfv := v.Field(j)
					//注入不同类型的值
					switch sf.Type.Kind() {
					case reflect.String:
						sfv.SetString(req.FormValue(sf.Name))
					case reflect.Int:
						conv, err := strconv.Atoi(req.FormValue(sf.Name))
						if err != nil {
							fmt.Println(err)
						} else {
							sfv.SetInt(int64(conv))
						}
					default:
						fmt.Println("注入参数失败!")
					}
				}
				*vs = append(*vs, v)
			}
		}
	}
	funcV.Call(*vs)
}