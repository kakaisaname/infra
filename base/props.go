package base

import (
	"github.com/kakaisaname/infra"
	"github.com/prometheus/common/log"
	"github.com/tietang/props/kvs"
	"sync"
)

var props kvs.ConfigSource

//配置会贯穿整个应用程序，需要向外暴露，在这定义一个对外暴露的函数，返回的是一个接口
//在任何地方都可以调用这个函数获取配置实例
func Props() kvs.ConfigSource {
	Check(props)
	return props
}

//项目启动的时候，最先starter Props，然后就会执行下面的Init函数，初始化props
type PropsStarter struct {
	infra.BaseStarter
}

//配置文件的初始化，读取配置文件
func (p *PropsStarter) Init(ctx infra.StarterContext) {
	//初始化配置 	获取配置 **
	props = ctx.Props()
	log.Info("初始化配置.")
	GetSystemAccount() //发红包的时候必须要这个账户，娶不到这个账户，红包业务不能正常进行，初始化的时候就需要
}

type SystemAccount struct {
	AccountNo   string
	AccountName string
	UserId      string
	Username    string
}

//系统账户
var systemAccount *SystemAccount
var systemAccountOnce sync.Once

//获取系统的账户配置									***
func GetSystemAccount() *SystemAccount {
	//只进行一次 									***
	systemAccountOnce.Do(func() {
		systemAccount = new(SystemAccount)
		//这是解析配置中的 system.account 的配置														***
		err := kvs.Unmarshal(Props(), systemAccount, "system.account")
		if err != nil {
			panic(err)
		}
	})
	return systemAccount
}

//获取红包活动链接 								***
func GetEnvelopeActivityLink() string {
	link := Props().GetDefault("envelope.link", "/v1/envelope/link")
	return link
}

//获取红包域名的函数
func GetEnvelopeDomain() string {
	domain := Props().GetDefault("envelope.domain", "http://localhost")
	return domain
}
