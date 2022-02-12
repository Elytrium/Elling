package elling

import (
	"context"
	"github.com/Elytrium/elling/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var DB *gorm.DB

type DBModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func LoadDatabase() {
	var dialector gorm.Dialector
	switch config.AppConfig.DBType {
	case config.MySQL:
		dialector = mysql.Open(config.AppConfig.DSN)
	case config.PostgreSQL:
		dialector = postgres.Open(config.AppConfig.DSN)
	case config.SQLite:
		dialector = sqlite.Open(config.AppConfig.DSN)
	default:
		log.Fatal().Str("db_type", string(config.AppConfig.DBType)).Msg("Current database type is not supported")
	}

	var err error
	DB, err = gorm.Open(dialector)

	log.Err(err).Msg("Database initialization finished")

	DB.Logger = &Logger{}

	_ = DB.AutoMigrate(&Balance{})
	_ = DB.AutoMigrate(&Product{})
	_ = DB.AutoMigrate(&User{})
}

type Logger struct{}

func (l *Logger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l *Logger) Info(_ context.Context, msg string, data ...interface{}) {
	log.Info().Interface("data", data).Msg(msg)
}

func (l *Logger) Warn(_ context.Context, msg string, data ...interface{}) {
	log.Warn().Interface("data", data).Msg(msg)
}

func (l *Logger) Error(_ context.Context, msg string, data ...interface{}) {
	log.Error().Interface("data", data).Msg(msg)
}

func (l *Logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), _ error) {
	elapsed := time.Since(begin).Milliseconds()

	if elapsed > config.AppConfig.SlowDBThreshold {
		sql, rows := fc()
		log.Warn().Interface("sql", sql).Interface("rows", rows).Int64("elapsed", elapsed).Msg("Slow SQL")
	}
}
