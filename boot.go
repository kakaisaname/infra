package infra

import (
	"github.com/kakaisaname/props/kvs"
	log "github.com/sirupsen/logrus"
	"reflect"
)

//应用程序   管理所有程序启动、加载的一个生命周期   **
//这个结构体 有配置文件，资源启动器上下文，
//资源启动器上下文， 用来在服务资源初始化、安装、启动和停止的生命周期中变量和对象的传递
type BootApplication struct {
	IsTest     bool
	conf       kvs.ConfigSource
	StarterCtx StarterContext
}

//构造系统  初始化配置文件
func New(conf kvs.ConfigSource) *BootApplication {
	e := &BootApplication{conf: conf, StarterCtx: StarterContext{}} //初始化配置文件 和 资源启动器上下文
	e.StarterCtx.SetProps(conf)
	return e
}

//项目启动 									***
func (b *BootApplication) Start() {
	//1.初始化 starter  (所有的)
	b.init()
	//2. 安装starter
	b.setup()

	//3. 启动starter
	b.start()
}

//程序初始化
func (e *BootApplication) init() {
	log.Info("Initializing starters...")
	for _, v := range GetStarters() {
		//获取启动器的类型
		typ := reflect.TypeOf(v)
		log.Debugf("Initializing: PriorityGroup=%d,Priority=%d,type=%s", v.PriorityGroup(), v.Priority(), typ.String())
		//对每个starter进行初始化 	**
		v.Init(e.StarterCtx)
	}
}

//程序安装
func (e *BootApplication) setup() {
	log.Info("Setup starters...")
	for _, v := range GetStarters() {
		typ := reflect.TypeOf(v)
		log.Debug("Setup: ", typ.String())
		v.Setup(e.StarterCtx)
	}
}

//starter启动																									***
//需要判断启动器是否是阻塞的，不是阻塞的直接启动，阻塞的，
// 判断是不是最后一个，是最后一个，直接启动，不是最后一个，开go协程去启动
func (e *BootApplication) start() {
	log.Info("Starting starters...")
	for i, v := range GetStarters() {
		typ := reflect.TypeOf(v)
		log.Debug("Starting: ", typ.String())
		if e.StarterCtx.Props().GetBoolDefault("testing", false) {
			go v.Start(e.StarterCtx)
			continue
		}

		if v.StartBlocking() {
			//如果是最后一个可阻塞的，直接启动并阻塞
			if i+1 == len(GetStarters()) {
				v.Start(e.StarterCtx)
			} else {
				//如果不是，使用goroutine来异步起动
				//防止阻塞后面starter
				go v.Start(e.StarterCtx)
			}
		} else {
			v.Start(e.StarterCtx)
		}
	}
}
