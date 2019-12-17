package base

import (
	//"git.imooc.com/wendell1000/infra"
	"github.com/kakaisaname/infra"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

var callbacks []func()

func Register(fn func()) {
	callbacks = append(callbacks, fn)
}

//定义一个Hook
type HookStarter struct {
	infra.BaseStarter
}

func (s *HookStarter) Init(ctx infra.StarterContext) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM) //几个信号量 			**去监听
	go func() {
		for {
			c := <-sigs
			log.Info("notify: ", c)
			for _, fn := range callbacks {
				fn() //执行注册的回调方法
			}
			break
			os.Exit(0)
		}
	}()

}

func (s *HookStarter) Start(ctx infra.StarterContext) {
	starters := infra.GetStarters()

	for _, s := range starters { //获取所有的 starters
		typ := reflect.TypeOf(s)
		log.Infof("【Register Notify Stop】:%s.Stop()", typ.String())
		Register(func() { //注册方法
			s.Stop(ctx)
		})
	}

}
