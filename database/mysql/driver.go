package mysql

import (
	"fmt"
	"os"

	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DriverName = "mysql"

type mysqlDb struct {
	database.BaseDatabaseImplement
	database.BaseDatabaseImplementV2

	initialized bool
}

func (s *mysqlDb) Init(options database.Options) error {
	// 防止重复初始化
	if s.initialized {
		return nil
	}
	s.initialized = true

	// 初始化日志器
	if options.Logger == nil {
		options.Logger = logger.Default()
	}

	s.SetLogger(options.Logger)
	options.Logger = logger.Default()
	options.Logger.Info(logger.NewFields().WithMessage("start open mysqlDb database").WithData(options.DataSource))

	// 连接数据库
	dataSource := options.DataSource
	db, openErr := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if openErr != nil {
		return fmt.Errorf("open mysqlDb database error: %w", openErr)
	}
	db.Logger = database.NewDBLogger(options.Logger)

	// 设置数据库连接池
	sqlDb, dbe := db.DB()
	if dbe != nil {
		return fmt.Errorf("get mysqlDb database error: %w", dbe)
	}
	s.BaseDatabaseImplement.ParseDatabaseOptions(sqlDb, options)
	s.BaseDatabaseImplement.SetRandCommand("rand()")
	s.BaseDatabaseImplement.SetDriverName(DriverName)

	// 连接成功
	s.BaseDatabaseImplement.Db, s.BaseDatabaseImplementV2.Db = db, db
	s.Logger.Info(logger.NewFields().WithMessage("successfully open mysqlDb database").WithData(dataSource))

	// 注册退出事件
	exit.RegisterExitEvent(func(_ os.Signal) {
		_ = sqlDb.Close()
		fmt.Println("closed mysql database")
	}, "CLOSE_MYSQL_DB_CONN")
	return nil
}

// NewMysqlDb creates a new mysql database instance.
//
// Deprecated: Use NewMySQLv2 instead.
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

// NewMySQLv2 creates a new mysql database instance.
func NewMySQLv2(config Config, models ...any) (db database.DatabaseV2, err error) {
	mysqlDb := &mysqlDb{}
	if initErr := mysqlDb.Init(convertConfigToOptions(config)); initErr != nil {
		return nil, fmt.Errorf("init mysqlDb database error: %w", initErr)
	}

	if migrateErr := mysqlDb.Migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate mysqlDb database error: %w", migrateErr)
	}

	return mysqlDb, nil
}

func NewWithLogger(config Config, logger logger.Logger, models ...any) (db database.DatabaseV2, err error) {
	mysqlDb := &mysqlDb{}
	opts := convertConfigToOptions(config)
	opts.Logger = logger
	if initErr := mysqlDb.Init(opts); initErr != nil {
		return nil, fmt.Errorf("init mysqlDb database error: %w", initErr)
	}

	if migrateErr := mysqlDb.Migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate mysqlDb database error: %w", migrateErr)
	}

	return mysqlDb, nil
}
