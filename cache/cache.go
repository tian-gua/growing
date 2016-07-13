package gcache

import (
	"time"
	"errors"
	"fmt"
)

const (
	IdleMode = iota //闲置模式
	Expire        //过期模式
)

//缓存
type cache struct {
	key      string        //key 冗余一个字典的索引
	value    interface{}   //value 使用接口实现泛型
	putTime  time.Time     //添加时间
	lifetime time.Time     //存活时间
	mode     int           //缓存过期模式
	modetime time.Duration //过期时间
}

//缓存接口
//type Icache interface {
//	Get(key string) interface{} //获得数据
//	Put(key string, value interface{}, lifetime  time.Time) //添加数据
//	IsExpire(d time.Duration) bool //是否过期
//	IsIdle(d time.Duration) bool //是否闲置
//}

//定义一个字典来存放key 和 cache
var cacheManager = make(map[string]cache)


//实现Icache接口的get方法
func (this *cache) get() (interface{}, error) {
	switch this.mode {
	case IdleMode:
		if this.isIdle(this.modetime) {
			return nil, errors.New("cache is idled")
		} else {
			this.updateLifeTime(time.Now().Add(this.modetime))
		}
	case Expire:
		if this.isExpire(this.modetime) {
			return nil, errors.New("cache is expired")
		}
	}
	return this.value, nil

}

//实现Icache接口的put方法,存数据
func (c *cache) put(key string, value interface{}, lifetime  time.Time, mode int) {
	c.key = key
	c.value = value
	c.lifetime = lifetime
	c.putTime = time.Now()
	c.mode = mode
}

//判断时候闲置
func (this *cache) isIdle(d time.Duration) bool {
	return time.Now().After(this.lifetime.Add(d))
}

//实现Icache接口的isExpire方发
func (this *cache) isExpire(d time.Duration) bool {
	return time.Now().After(this.putTime.Add(d))
}

//更新存货时间
func (this *cache) updateLifeTime(t time.Time) {
	this.lifetime = t
}


//实现cache包的公共方法,用于获取数据
func Get(key string) interface{} {
	//从字典子查找是否存在cache对象,如果存在返回cache对象的value属性
	if v, ok := cacheManager[key]; ok {
		value, err := v.get()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return value
	}
	return nil
}


//实现cache包的公共方法,用于存放数据
func Put(key string, value interface{}, lifeTime  time.Time, mode int) {

	//新建一个cache对象
	newCache := &cache{}
	//初始化cache对象
	newCache.put(key, value, lifeTime, mode)
	//讲对象放到字典里
	cacheManager[key] = *newCache

}
