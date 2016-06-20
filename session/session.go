package gsession

import (
	"net/http"
	"time"
	"errors"
	"cache"
)

//session对象
type Session struct {
	list map[string]interface{}
}

//获取session的值
func (s *Session) Get(key string) interface{} {

	if k, ok := s.list[key]; ok {
		return k
	}
	return nil
}

//存放值到session对象里
func (s *Session) Put(key string, value interface{}) {

	s.list[key] = value

}


//默认缓存存在30分钟
var expires = 1000 * 60 * 30

//获取Session对象
func GetSession(req http.Request) (*Session, error) {
	//获得sessionid的cookie
	cookie, err := req.Cookie("gsessionid")
	if err != nil {
		return nil, err
	}
	//如果请求不存在 gsessionid 则添加一个sessionid过去
	if cookie == nil {
		ck := &http.Cookie{Name:"gsessionid", Value:"gsession", Expires:time.Now().Add(expires) }
		req.AddCookie(ck)
		//新建一个session对象
		newSession := &Session{list:make(map[string]interface{})}
		//新session放缓存里
		gcache.Put("gsession", newSession, 888888)
		return newSession, nil

	} else {
		//判断gsession是否过期
		if time.Now().After(cookie.Expires) {
			return errors.New("expired")
		}

		return gcache.Get(cookie.Value)
	}

}



