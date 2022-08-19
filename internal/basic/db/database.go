package db

import (
	"context"
	"cron/internal/basic/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

type Database struct {
	Write *MyDB
	Read *MyDB
}

var (
	write *gorm.DB
	read *gorm.DB
	once sync.Once
)

// 连接数据库
func New(ctx context.Context)*Database{
	once.Do(func() {
		conf := config.DbConf()
		write = conn(conf["write"])
		read = conn(conf["read"])
	})

	// 根据实例,修改上下文
	return &Database{
		Write: &MyDB{read.WithContext(ctx)},
		Read:  &MyDB{write.WithContext(ctx)},
	}
}

func conn(conf config.DataBaseConf) *gorm.DB {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conf.Source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "", // 表前缀
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败 %w", err))
	}
	// 启用连接池
	if err = polling(db); err != nil {
		panic(fmt.Errorf("连接池设置异常 %w", err))
	}
	// 调试模式
	if conf.Debug {
		db = db.Debug()
	}
	return db
}


// 设置程序池;
// 这个有空要研究一下
func polling(_db *gorm.DB) error {
	sqlDb, err := _db.DB()
	if err != nil {
		return err
	}
	// TODO: 这些参数要迁移到配置文件中
	// 设置连接池;
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDb.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDb.SetConnMaxLifetime(time.Hour)
	return nil
}