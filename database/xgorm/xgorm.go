package xgorm

import (
	"github.com/rabbit-rm/xgo/database"
	"github.com/rabbit-rm/xgo/xerror"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	config *database.Config
	db     *gorm.DB
}

func NewGormDB(config *database.Config) *GormDB {
	return &GormDB{
		config: config,
	}
}

func (gdb *GormDB) Connect() error {
	logLevel := logger.Silent
	if gdb.config.DebugSQL {
		logLevel = logger.Info
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(mysql.Open(gdb.config.DSN), config)
	if err != nil {
		return xerror.Wrap(err, "connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return xerror.Wrap(err, "underlying to sql.DB")
	}

	sqlDB.SetMaxOpenConns(gdb.config.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(gdb.config.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(gdb.config.MaxLifeConnections)

	gdb.db = db
	return nil
}

func (gdb *GormDB) Close() error {
	if gdb.db != nil {
		sqlDB, err := gdb.db.DB()
		if err != nil {
			return xerror.Wrap(err, "underlying to sql.DB")
		}
		return sqlDB.Close()
	}
	return nil
}

func (gdb *GormDB) Ping() error {
	if gdb.db == nil {
		return xerror.New("database connection not established")
	}
	sqlDB, err := gdb.db.DB()
	if err != nil {
		return xerror.Wrap(err, "underlying to sql.DB")
	}
	return sqlDB.Ping()
}

func (gdb *GormDB) DB() *gorm.DB {
	return gdb.db
}
