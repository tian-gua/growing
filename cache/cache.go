package cache

import (
	"time"
	"errors"
)


//缓存
type cache struct {
	key      string        //key 冗余一个字典的索引
	value    interface{}   //value 使用接口实现泛型
	putTime  time.Time     //添加时间
	lifetime time.Duration //存活时间
}

//缓存接口
type Icache interface {
	get(key string) interface{} //获得数据
	put(key string, value interface{}, lifetime  time.Time) //添加数据
	isExpire(v cache) bool //是否过期
}

//定义一个字典来存放key 和 cache
var cacheManager = make(map[string]cache)


//实现Icache接口的get方法
func (c *cache) get() (interface{}, error) {

	//如果值过期了 则 new一个cache结构体,并返回error
	if err := c.isExpire(); err != nil {
		return nil, err
	}
	//如果没过期,返回value,error返回nil
	return c.value, nil

}

//实现Icache接口的put方法,存数据
func (c *cache) put(key string, value interface{}, lifetime  time.Duration) {
	c.key = key
	c.value = value
	c.lifetime = lifetime
	c.putTime = time.Now()
}


//实现Icache接口的isExpire方发
func (c *cache) isExpire() error {

	//获得cache对象过期时间
	t := c.putTime.Add(c.lifetime)
	//获取当前时间
	now := time.Now()
	//判断是否过期
	if now.After(t) {
		return errors.New("expire")
	}
	return nil

}



//实现cache包的公共方法,用于获取数据
func Get(key string) (value interface{}, err error) {

	//从字典子查找是否存在cache对象,如果存在返回cache对象的value属性
	if v, ok := cacheManager[key]; ok {
		value, err := v.get()
		return value, err
	}
	//如果字典找不到cache对象,返回错误 not found
	return nil, errors.New("NOF")
}


//实现cache包的公共方法,用于存放数据
func Put(key string, value interface{}, long  time.Duration) {

	//新建一个cache对象
	newCache := &cache{}
	//初始化cache对象
	newCache.put(key, value, long)
	//讲对象放到字典里
	cacheManager[key] = *newCache

}
