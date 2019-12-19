package infra

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"reflect"
	"sort"
)

const (
	KeyProps = "_conf"
)

//默认的空实现,方便资源启动器的实现
type BaseStarter struct {
}

//资源启动器上下文，
// 用来在服务资源初始化、安装、启动和停止的生命周期中变量和对象的传递
//是一个map
type StarterContext map[string]interface{}

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没有被初始化")
	}
	//类型转化 	***
	return p.(kvs.ConfigSource) //只能接口才能这样,进行类型转化
}

//设置配置   **
func (s StarterContext) SetProps(conf kvs.ConfigSource) {
	s[KeyProps] = conf
}

//定义的启动器的接口，需要实现的方法  																****

//资源启动器，每个应用少不了依赖其他资源，比如数据库，缓存，消息中间件等等服务
//启动器实现类，不需要实现所有方法，只需要实现对应的阶段方法即可，可以嵌入@BaseStarter
//通过实现资源启动器接口和资源启动注册器，友好的管理这些资源的初始化、安装、启动和停止。
//Starter对象注册器，所有需要在系统启动时需要实例化和运行的逻辑，都可以实现此接口
//注意只有Start方法才能被阻塞，如果是阻塞Start()，同时StartBlocking()要返回true
type Starter interface {
	//资源初始化和，通常把一些准备资源放在这里运行
	Init(StarterContext)
	//资源的安装，所有启动需要的具备条件，使得资源达到可以启动的就备状态    **
	Setup(StarterContext)
	//启动资源，达到可以使用的状态
	Start(StarterContext)
	//说明该资源启动器开始启动服务时，是否会阻塞
	//如果存在多个阻塞启动器时，只有最后一个阻塞，之前的会通过goroutine来异步启动
	//所以，需要规划好启动器注册顺序
	StartBlocking() bool
	//资源停止：
	// 通常在启动时遇到异常时或者启用远程管理时，用于释放资源和终止资源的使用，
	// 通常要优雅的释放，等待正在进行的任务继续，但不再接受新的任务
	Stop(StarterContext)
	PriorityGroup() PriorityGroup
	Priority() int
}

//服务启动注册器  								***
//不用需外部构造，全局只有一个
//所有的starter都会注册到这里面   ，会把starter存储在切片里，按照切片的数据一个一个的初始化
type starterRegister struct {
	nonBlockingStarters []Starter
	blockingStarters    []Starter
}

//已经初始化好的一个starter注册器 	**
var StarterRegister *starterRegister = &starterRegister{}

//返回所有的启动器，所有的阻塞非阻塞的都加进去															***
func (r *starterRegister) AllStarters() []Starter {
	starters := make([]Starter, 0)
	starters = append(starters, r.nonBlockingStarters...)
	starters = append(starters, r.blockingStarters...)
	return starters
}

//注册starter 																						***
func (r *starterRegister) Register(starter Starter) {
	if starter.StartBlocking() {
		r.blockingStarters = append(r.blockingStarters, starter)
	} else {
		r.nonBlockingStarters = append(r.nonBlockingStarters, starter)
	}
	typ := reflect.TypeOf(starter)
	log.Infof("Register starter: %s", typ.String())
}

//启动器切片
type Starters []Starter

//进行排序 									**
func (s Starters) Len() int {
	return len(s)
}

func (s Starters) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Starters) Less(i, j int) bool {
	return s[i].PriorityGroup() > s[j].PriorityGroup() && s[i].Priority() > s[j].Priority()
}

//注册starter   ***
func Register(starter Starter) {
	StarterRegister.Register(starter)
}

//对starters 进行排序
func SortStarters() {
	sort.Sort(Starters(StarterRegister.AllStarters()))
}

//获取所有注册的starter
//所有的starter是一个切片
func GetStarters() []Starter {
	return StarterRegister.AllStarters()
}

type PriorityGroup int

const (
	SystemGroup         PriorityGroup = 30
	BasicResourcesGroup PriorityGroup = 20
	AppGroup            PriorityGroup = 10

	INT_MAX          = int(^uint(0) >> 1)
	DEFAULT_PRIORITY = 10000
)

//每个启动器都包含了BaseStarter																	*****
//默认的空实现,方便资源启动器的实现
func (s *BaseStarter) Init(ctx StarterContext)  {}
func (s *BaseStarter) Setup(ctx StarterContext) {}
func (s *BaseStarter) Start(ctx StarterContext) {}
func (s *BaseStarter) Stop(ctx StarterContext)  {}

func (s *BaseStarter) StartBlocking() bool          { return false }
func (s *BaseStarter) PriorityGroup() PriorityGroup { return BasicResourcesGroup }
func (s *BaseStarter) Priority() int                { return DEFAULT_PRIORITY }
