package clickhouse

import (
	"context"
	"fmt"
	"strings"
	"time"

	clickhousego "github.com/ClickHouse/clickhouse-go/v2"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseClickhouse struct {
	db *gorm.DB
}

var (
	dbStorage *gorm.DB
)

func ConnectClickhouse(ctx context.Context, cfg *ClickhouseConfig) (*DatabaseClickhouse, error) {
	if dbStorage != nil {
		return &DatabaseClickhouse{db: dbStorage}, nil
	}

	if strings.HasPrefix(cfg.Host, "https://") {
		cfg.Host = cfg.Host[8:]
	}

	if strings.HasPrefix(cfg.Host, "http://") {
		cfg.Host = cfg.Host[7:]
	}

	if cfg.Host[len(cfg.Host)-1:] == "/" {
		cfg.Host = cfg.Host[:len(cfg.Host)-1]
	}

	sqlDB := clickhousego.OpenDB(&clickhousego.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: clickhousego.Auth{
			Database: cfg.DBName,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
	sqlDB.SetConnMaxIdleTime(2 * time.Hour)

	db, err := gorm.Open(clickhouse.New(clickhouse.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(func() logger.LogLevel {
			switch cfg.LogLevel {
			case int(logger.Silent):
				return logger.Silent
			case int(logger.Error):
				return logger.Error
			case int(logger.Warn):
				return logger.Warn
			case int(logger.Info):
				return logger.Info
			default:
				return logger.Silent
			}
		}()),
	})

	if err != nil {
		return nil, err
	}

	dbStorage = db

	return &DatabaseClickhouse{db: dbStorage}, nil
}

func (d *DatabaseClickhouse) GetDB() *gorm.DB {
	return dbStorage
}
