package connection

import (
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/HuangXiaoL/xiaoshuo/internal/pkg/config"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	db *sqlx.DB
)

func initDB() error {
	dbConfig := config.Get().Database
	logrus.WithField("user", dbConfig.User).Debug("open database")
	var err error
	dsn := dbConfig.User + ":" + dbConfig.Password + "@tcp(" + dbConfig.IP + dbConfig.Port + ")/" + dbConfig.Name + "?charset=utf8mb4&parseTime=True"
	if db, err = sqlx.Connect("mysql", dsn); err != nil {
		return err
	}

	db.SetMaxOpenConns(dbConfig.MaxConn)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(5 * time.Minute)

	return nil
}
