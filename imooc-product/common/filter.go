package common

import (
	"net/http"
	"strings"
)

// 声明一个新的数据类型（函数类型）
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// 拦截器结构体
type Filter struct {
	// 用来存储需要拦截的URI
	filterMap map[string]FilterHandle
}

// Filter初始化函数
func NewFilter()*Filter  {
	return &Filter{filterMap:make(map[string]FilterHandle)}
}

// 注册拦截器
func (f *Filter) RegisterFilterUrl(url string, handler FilterHandle)  {
	f.filterMap[url] = handler
}

// 根据url获取对应的handle
func (f *Filter) GetFilterHandle(url string) FilterHandle  {
	return f.filterMap[url]
}

// 声明新的函数类型
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// 执行拦截器，返回函数类型
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, req *http.Request)  {
	return func(rw http.ResponseWriter, req *http.Request) {
		for path, handle := range f.filterMap {
			if strings.Contains(req.RequestURI, path) {
				// 执行拦截业务
				err := handle(rw,req)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				// 跳出循环
				break
			}
		}
		// 执行正常注册的函数
		webHandle(rw, req)
	}
}

