package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
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

func OpenMysqlDefault() (*gorm.DB, error) {
	dsn := &DSNConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_ROOT_PASSWORD"),
		Net:      "tcp",
		Addr:     os.Getenv("DATABASE_URL"),
		DBName:   os.Getenv("DATABASE_SCHEMA"),
	}
	return gorm.Open(mysql.Open(dsn.FormatDSN()), &gorm.Config{})
}

func OpenMysql(config DSNConfig) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(config.FormatDSN()), &gorm.Config{})
}
