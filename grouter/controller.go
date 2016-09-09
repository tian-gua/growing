package grouter

import (
	"net/http"
	"fmt"
	"regexp"
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
	reg, err := regexp.Compile("[/\\w]*")
	if err != nil {
		fmt.Println(err)
		return url
	}
	return reg.FindString(url)
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
	//获取函数的反射type和value
	funcType := reflect.TypeOf(h)
	funcValue := reflect.ValueOf(h)
	//获取函数的参数个数
	numin := funcType.NumIn()
	paramSlice := make([]reflect.Value, 0)
	//遍历参数
	for i := 0; i < numin; i++ {
		//获取参数的反射Type
		paramType := funcType.In(i)
		//处理参数
		//如果参数是指针并且是指向http.Request结构体则不处理(注入参数)，把参数放入参数切片里
		//如果参数是http.ResponseWriter也不处理，加入参数切片
		//如果是其他结构体，则遍历属性，把表单的值注入进去。
		if paramType.Kind() == reflect.Ptr {
			if "Request" == paramType.Elem().Name() {
				paramSlice = append(paramSlice, reflect.ValueOf(req))
			}
		} else {
			if "ResponseWriter" == paramType.Name() {
				paramSlice = append(paramSlice, reflect.ValueOf(rw))
			} else if reflect.Struct == paramType.Kind() {
				//通过参数反射类型新建一个参数Value
				newStructParam := reflect.New(paramType).Elem()
				//遍历结构体属性,并赋值
				//值从表单里取,根据属性名字对应
				for j := 0; j < paramType.NumField(); j++ {
					//获取结构体字段类型反射
					structField := paramType.Field(j)
					//获取结构体字段值反射
					structFieldValue := newStructParam.Field(j)

					//获取表单的key
					//之前是用的结构体字段的name，但是大小写很麻烦
					//现在是使用字段的tag，这里就定tag为"key"
					key := structField.Tag.Get("key")

					//注入不同类型的值
					switch structField.Type.Kind() {
					case reflect.String:
						//修改字段的字符串值
						structFieldValue.SetString(req.FormValue(key))
					case reflect.Int:
						//获取表单的值，并转换成int类型，并赋值给结构体参数的字段
						if formv := req.FormValue(key); len(formv) != 0 {
							conv, err := strconv.Atoi(formv)
							if err != nil {
								fmt.Println(err)
							} else {
								structFieldValue.SetInt(int64(conv))
							}
						}
					default:
						fmt.Println("注入参数失败!")
					}
				}
				//将结构体参数放入参数切片
				paramSlice = append(paramSlice, newStructParam)
			}
		}
	}
	//调用handler，并将上面构造的结构体参数传入
	funcValue.Call(paramSlice)
}