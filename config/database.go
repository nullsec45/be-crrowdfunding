package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql" 
	"gorm.io/gorm"
)

type Mysql struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionMysql() (*Mysql, error) {
	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.DBName)

	db, err := gorm.Open(mysql.Open(dbConnString), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("[ConnectionMysql-1] Failed to connect to database " + cfg.Mysql.Host)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionMysql-2] Failed to get database connection")
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Mysql.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Mysql.DBMaxIdle)

	return &Mysql{DB: db}, nil
}