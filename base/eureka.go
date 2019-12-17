package base

import (
	"github.com/kakaisaname/infra"
	"github.com/kataras/iris"
	"github.com/tietang/go-eureka-client/eureka"
	"time"
)

var eurekaClient *eureka.Client

func EurekaClient() *eureka.Client {
	Check(eurekaClient)
	return eurekaClient
}

type EurekaStarter struct {
	infra.BaseStarter
	client *eureka.Client
}

//下面注释的代码后面rpc 会放开
func (e *EurekaStarter) Init(ctx infra.StarterContext) {
	e.client = eureka.NewClient(ctx.Props())
	//rpcPort := ctx.Props().GetDefault("app.rpc.port", "18082")
	//e.client.InstanceInfo.Metadata.Map["rpcPort"] = rpcPort
	e.client.Start()
	//e.client.Applications, _ = e.client.GetApplications()
	//eurekaClient = e.client
}

//http://127.0.0.1:18080/health
//http://127.0.0.1:18080/info

func (e *EurekaStarter) Start(ctx infra.StarterContext) {
	info := make(map[string]interface{})
	info["startTime"] = time.Now()
	info["appName"] = ctx.Props().GetDefault("app.name", "resk") //应用名称
	Iris().Get("/info", func(context iris.Context) {
		context.JSON(info)
	})
	Iris().Get("/health", func(context iris.Context) {
		health := eureka.Health{
			Details: make(map[string]interface{}),
		}
		health.Status = eureka.StatusUp
		context.JSON(health)

	})
}
