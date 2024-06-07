package database

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type DSNConfig struct {
	User     string
	Password string
	Net      string
	Addr     string
	DBName   string
}

func (dsnConfig *DSNConfig) FormatDSN() string {
	return fmt.Sprintf("%s:%s@%s(%s)/%s",
		dsnConfig.User,
		dsnConfig.Password,
		dsnConfig.Net,
		dsnConfig.Addr,
		dsnConfig.DBName,
	)
}

func ConnectMysql() (*gorm.DB, error) {
	dsn := &DSNConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_ROOT_PASSWORD"),
		Net:      "tcp",
		Addr:     os.Getenv("DATABASE_URL"),
		DBName:   os.Getenv("DATABASE_SCHEMA"),
	}
	return gorm.Open(mysql.Open(dsn.FormatDSN()), &gorm.Config{})
}

func ConnectRedis() (*redis.Client, error) {
	dbIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
		DB:   dbIndex,
	})
	return rdb, nil
}
