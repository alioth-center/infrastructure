package postgres

import (
	"fmt"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDb struct {
	database.BaseDatabaseImplement
}

func (s *postgresDb) Init(options database.Options) error {
	// 初始化日志
	s.BaseDatabaseImplement.ParseLoggerOptions(options)
	s.Logger.Info(logger.NewFields().WithMessage("start open postgresDb database").WithData(options.DataSource))

	// 连接数据库
	dataSource := options.DataSource
	db, openErr := gorm.Open(postgres.Open(dataSource), &gorm.Config{})
	if openErr != nil {
		return fmt.Errorf("open postgresDb database error: %w", openErr)
	}

	// 设置数据库连接池
	sqlDb, dbe := db.DB()
	if dbe != nil {
		return fmt.Errorf("get postgresDb database error: %w", dbe)
	}
	s.BaseDatabaseImplement.ParseDatabaseOptions(sqlDb, options)
	s.BaseDatabaseImplement.SetRandCommand("random()")

	// 连接成功
	s.Db = db
	s.Logger.Info(logger.NewFields().WithMessage("successfully open postgresDb database").WithData(dataSource))

	// 注册退出事件
	exit.Register(func(_ string) string {
		_ = sqlDb.Close()
		return "closed postgres database"
	}, "postgres database")
	return nil
}

func NewPostgresDb(config Config, models ...any) (db database.Database, err error) {
	postgresDb := &postgresDb{}
	if initErr := postgresDb.Init(convertConfigToOptions(config)); initErr != nil {
		return nil, fmt.Errorf("init postgresDb database error: %w", initErr)
	} else if migrateErr := postgresDb.Migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate postgresDb database error: %w", migrateErr)
	} else {
		return postgresDb, nil
	}
}
