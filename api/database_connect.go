package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
)

func dbConnect(retriesLeft int) error {
	con := os.Getenv("DATABASE_URL")
	dbType := os.Getenv("DATABASE_TYPE")
	logger.Infof("opening connection to database: %s", con)
	var err error
	db, err = gorm.Open(dbType, con)
	if err != nil {
		logger.Errorf("cannot open connection to database: %v", err)
		return err
	}
	err = db.DB().Ping()
	if err != nil {
		if retriesLeft > 0 {
			logger.Errorf("cannot talk to database, retrying in 10 seconds (%d attempts left): %v", retriesLeft-1, err)
			time.Sleep(10 * time.Second)
			return dbConnect(retriesLeft - 1)
		} else {
			logger.Errorf("cannot talk to database, last attempt failed: %v", err)
			return err
		}
	}
	initDb()
	maxIdleConnections, err := strconv.Atoi(os.Getenv("MAX_IDLE_DB_CONNECTIONS"))
	if err != nil {
		logger.Warningf("cannot parse COMMENTO_MAX_IDLE_DB_CONNECTIONS: %v", err)
		maxIdleConnections = 50
	}
	db.DB().SetMaxIdleConns(maxIdleConnections)
	return nil
}
