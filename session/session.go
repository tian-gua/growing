package gsession

import (
	"github.com/aidonggua/growing/cache"
	"net/http"
)

//session对象
type Session struct {
	attributes map[string]interface{}
}

//获取session的值
func (s *Session) Get(key string) interface{} {
	if k, ok := s.attributes[key]; ok {
		return k
	}
	return nil
}

//存放值到session对象里
func (s *Session) Put(key string, value interface{}) {
	s.attributes[key] = value
}

//获取Session对象
func GetSession(rw http.ResponseWriter, req *http.Request) (*Session, error) {
	//获得sessionid的cookie
	cookie, err := req.Cookie("gsessionid")
	if err != nil {
		//如果请求不存在 gsessionid 则添加一个sessionid过去
		if err == http.ErrNoCookie {
			return newSession(rw)
		}
	}
	session := gcache.Get(cookie.Value)
	//如果session为空证明session过期了
	if session == nil {
		return newSession(rw)
	}
	return session.(*Session), nil
}

//新建一个session
func newSession(rw http.ResponseWriter) (*Session, error) {
	sessionId := getSessionId()
	ck := &http.Cookie{Name: "gsessionid", Value: sessionId, MaxAge: 0, Path: "/", HttpOnly: true}
	http.SetCookie(rw, ck)
	//新建一个session对象
	newSession := &Session{attributes: make(map[string]interface{})}
	//新session放缓存里
	gcache.Put(sessionId, newSession, gcache.IdleMode, Session_time)
	return newSession, nil
}
