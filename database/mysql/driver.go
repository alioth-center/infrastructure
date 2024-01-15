package mysql

import (
	"fmt"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DriverName = "mysql"

type mysqlDb struct {
	database.BaseDatabaseImplement
}

func (s *mysqlDb) Init(options database.Options) error {
	// 初始化日志
	s.BaseDatabaseImplement.ParseLoggerOptions(options)
	s.Logger.Info(logger.NewFields().WithMessage("start open mysqlDb database").WithData(options.DataSource))

	// 连接数据库
	dataSource := options.DataSource
	db, openErr := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if openErr != nil {
		return fmt.Errorf("open mysqlDb database error: %w", openErr)
	}

	// 设置数据库连接池
	sqlDb, dbe := db.DB()
	if dbe != nil {
		return fmt.Errorf("get mysqlDb database error: %w", dbe)
	}
	s.BaseDatabaseImplement.ParseDatabaseOptions(sqlDb, options)
	s.BaseDatabaseImplement.SetRandCommand("rand()")
	s.BaseDatabaseImplement.SetDriverName(DriverName)

	// 连接成功
	s.Db = db
	s.Logger.Info(logger.NewFields().WithMessage("successfully open mysqlDb database").WithData(dataSource))

	// 注册退出事件
	exit.Register(func(_ string) string {
		_ = sqlDb.Close()
		return "closed mysql database"
	}, "mysql database")
	return nil
}

func NewMysqlDb(config Config, models ...any) (db database.Database, err error) {
	mysqlDb := &mysqlDb{}
	if initErr := mysqlDb.Init(convertConfigToOptions(config)); initErr != nil {
		return nil, fmt.Errorf("init mysqlDb database error: %w", initErr)
	} else if migrateErr := mysqlDb.Migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate mysqlDb database error: %w", migrateErr)
	} else {
		return mysqlDb, nil
	}
}
