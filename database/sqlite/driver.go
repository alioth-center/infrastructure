package sqlite

import (
	"fmt"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type sqliteDb struct {
	database.BaseDatabaseImplement
}

func (s *sqliteDb) Init(options database.Options) error {
	// 初始化日志
	s.BaseDatabaseImplement.ParseLoggerOptions(options)
	dataSource := options.DataSource
	s.Logger.Info(logger.NewFields().WithMessage("start open sqliteDb database").WithData(dataSource))

	// 连接数据库
	db, openErr := gorm.Open(sqlite.Open(dataSource))
	if openErr != nil {
		if !os.IsNotExist(openErr) {
			err := fmt.Errorf("open sqliteDb database error: %w", openErr)
			s.Logger.Error(logger.NewFields().WithMessage(err.Error()).WithData(dataSource))
			return err
		} else if _, createErr := os.Create(filepath.Join("data", dataSource)); createErr != nil {
			err := fmt.Errorf("create sqliteDb database error: %w", createErr)
			s.Logger.Error(logger.NewFields().WithMessage(err.Error()).WithData(dataSource))
			return err
		}
	}

	// 设置数据库连接池
	sqlDb, dbe := db.DB()
	if dbe != nil {
		return fmt.Errorf("get sqliteDb database error: %w", dbe)
	}
	s.BaseDatabaseImplement.ParseDatabaseOptions(sqlDb, options)
	s.BaseDatabaseImplement.SetRandCommand("random()")

	// 连接成功
	s.Db = db
	s.Logger.Info(logger.NewFields().WithMessage("successfully open sqliteDb database").WithData(dataSource))

	// 注册退出事件
	exit.Register(func(_ string) string {
		_ = sqlDb.Close()
		return "closed sqlite database"
	}, "sqlite database")
	return nil
}

func NewSqliteDb(config Config, models ...any) (db database.Database, err error) {
	sqliteDb := &sqliteDb{}
	if initErr := sqliteDb.Init(convertConfigToOptions(config)); initErr != nil {
		return nil, fmt.Errorf("init sqliteDb database error: %w", initErr)
	} else if migrateErr := sqliteDb.Migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate sqliteDb database error: %w", migrateErr)
	} else {
		return sqliteDb, nil
	}
}
