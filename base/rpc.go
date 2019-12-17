package base

import (
	"github.com/kakaisaname/infra"
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"reflect"
)

var rpcServer *rpc.Server

func RpcServer() *rpc.Server {
	Check(rpcServer)
	return rpcServer
}

//rpc暴露后需要一个注册函数 											***
func RpcRegister(ri interface{}) {
	typ := reflect.TypeOf(ri)
	log.Infof("goRPC Register: %s", typ.String())
	RpcServer().Register(ri)
}

type GoRPCStarter struct {
	infra.BaseStarter
	server *rpc.Server
}

//GoRPCStarter 要在 GoRpcApiStarter 之前							***
func (s *GoRPCStarter) Init(ctx infra.StarterContext) {
	s.server = rpc.NewServer() //先new好，得出rpcServer,然后才可以调用这个 RpcRegister		**
	rpcServer = s.server
}
func (s *GoRPCStarter) Start(ctx infra.StarterContext) {

	port := ctx.Props().GetDefault("app.rpc.port", "8082") //获取配置的端口
	//监听网络端口
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panic(err) //如果有错误，抛出错误
	}
	log.Info("tcp port listened for rpc:", port)
	//处理网络连接和请求
	go s.server.Accept(listener) //一个协程去处理
}
