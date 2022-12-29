package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"medicinal_share/gen/out/dao"
	"sync"
	"time"
)

var db *gorm.DB

var lock = sync.Mutex{}

const (
	dsnMaster = "root:CZLczl@20010821@tcp(localhost:3316)/medicinal_share?parseTime=true&loc=Local"
	dsnSlave1 = "root:CZLczl@20010821@tcp(localhost:3317)/medicinal_share?parseTime=true&loc=Local"
	dsnSlave2 = "root:CZLczl@20010821@tcp(localhost:3318)/medicinal_share?parseTime=true&loc=Local"
)

// CountPolicy 实现数据库负载均衡使用轮询的策略
type CountPolicy struct {
	now int
}

func init() {

	db = newDb()
	dao.SetDefault(db)

}

func (cp *CountPolicy) Resolve(pool []gorm.ConnPool) gorm.ConnPool {
	cp.now++
	cp.now %= len(pool)
	con := pool[cp.now]
	return con
}

func GetConnect() *gorm.DB {
	return db
}

func newDb() *gorm.DB {
	DB, err := gorm.Open(mysql.Open(dsnMaster), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	DB.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{DB.Dialector},
		Replicas: []gorm.Dialector{
			mysql.Open(dsnSlave1),
			mysql.Open(dsnSlave2),
		},
		Policy: &CountPolicy{
			now: 0,
		},
	}))

	DB.Use(dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(5 * time.Minute).
		SetMaxOpenConns(20).
		SetMaxIdleConns(10))
	return DB
}
