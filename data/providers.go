package data

import (
	"Cookie_Session/data/common"
	"Cookie_Session/data/memory"
)

var (
	Provides map[string]common.IProvider
)

func init(){
	Provides = make(map[string]common.IProvider) //初始化数据存储接口
	//注册到memory存储接口
	pder := memory.NewProvider()
	Register("memory", pder)
}

//注册provider
func Register(name string, provider common.IProvider){
	if provider == nil {
		panic("session: Register provider is nil")
	}

	if _, dup := Provides[name]; dup {
		panic("session: Register called twice for provider" + name)
	}

	Provides[name] = provider
}
