package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kakaisaname/infra"
	"github.com/kakaisaname/infra/logrus"
	"github.com/kakaisaname/props/kvs"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

//dbx 数据库实例
var database *dbx.Database

func DbxDatabase() *dbx.Database {
	Check(database)
	return database
}

//dbx 数据库starter，并且设置为全局   实现了starter接口
type DbxDatabaseStarter struct {
	infra.BaseStarter
}

//数据库的初始化   数据库的连接启动稍晚一点，所以放在setUp阶段
func (s *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	//拿到配置文件
	conf := ctx.Props()
	//数据库配置
	settings := dbx.Settings{}
	//kvs.Unmarshal 直接解析conf下的内容到结构体里
	err := kvs.Unmarshal(conf, &settings, "mysql")
	if err != nil {
		panic(err)
	}
	log.Info("mysql.conn url:", settings.ShortDataSourceName())
	db, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	log.Info(db.Ping())
	db.SetLogger(logrus.NewUpperLogrusLogger())
	database = db
}
