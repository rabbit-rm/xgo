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

func (g *GormDB) Connect() error {
	logLevel := logger.Silent
	if g.config.DebugSQL {
		logLevel = logger.Info
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(mysql.Open(g.config.DSN), config)
	if err != nil {
		return xerror.Wrap(err, "connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return xerror.Wrap(err, "underlying to sql.DB")
	}

	sqlDB.SetMaxOpenConns(g.config.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(g.config.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(g.config.MaxLifeConnections)

	g.db = db
	return nil
}

func (g *GormDB) Close() error {
	if g.db != nil {
		sqlDB, err := g.db.DB()
		if err != nil {
			return xerror.Wrap(err, "underlying to sql.DB")
		}
		return sqlDB.Close()
	}
	return nil
}

func (g *GormDB) Ping() error {
	if g.db == nil {
		return xerror.New("database connection not established")
	}
	sqlDB, err := g.db.DB()
	if err != nil {
		return xerror.Wrap(err, "underlying to sql.DB")
	}
	return sqlDB.Ping()
}

func (g *GormDB) DB() *gorm.DB {
	return g.db
}
