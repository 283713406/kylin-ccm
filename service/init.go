package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

var (
	o                 orm.Ormer
	tablePrefix       string             // 表前缀
	ClusterService    *clusterService       // 用户服务
)

func Init() {
	dbHost := beego.AppConfig.String("db.host")
	dbPort := beego.AppConfig.String("db.port")
	dbUser := beego.AppConfig.String("db.user")
	dbPassword := beego.AppConfig.String("db.password")
	dbName := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	tablePrefix = beego.AppConfig.String("db.prefix")

	if dbPort == "" {
		dbPort = "3306"
	}
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)

	orm.RegisterModelWithPrefix(tablePrefix,
		new(Cluster),
	)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	o = orm.NewOrm()
	orm.RunCommand()

	// 初始化服务对象
	initService()
}

func initService() {
	ClusterService = &clusterService{}
}

// 返回真实表名
func tableName(name string) string {
	return tablePrefix + name
}
