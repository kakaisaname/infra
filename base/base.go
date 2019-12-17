package base

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

//结构体指针检查验证，如果传入的interface为nil，就通过log.Panic函数抛出一个异常
//被用在starter中检查公共资源是否被实例化了
//check函数的作用是检查传入的参数是否为nil，如果是nil，会报错
func Check(a interface{}) {
	if a == nil {

		//skip 代表层级
		//层次为0的时候返回我们调用runtime.Caller的地方.为1的时候就是我们调用call函数的地方
		//2 3 是go源码的调用
		_, file, line, _ := runtime.Caller(1)
		strs := strings.Split(file, "/")
		size := len(strs)
		if size > 4 {
			size = 4
		}
		//长度大于4的设置为4，并去掉前4长度字符   如果没有大于4，则是全部返回
		file = filepath.Join(strs[len(strs)-size:]...)
		log.Panicf("object can't be nil, cause by: %s(%d)", file, line)
	}
}
