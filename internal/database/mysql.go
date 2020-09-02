package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //You could import dialect
	"github.com/rahulshibu/go-echo-boilerplate/internal/configs"

	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

//ConnectSQL - connect to sql server
func ConnectSQL() *gorm.DB {

	var err error

	var mysqlHost = fmt.Sprint(configs.AppConfig.Database.User, ":", configs.AppConfig.Database.Password, "@(", configs.AppConfig.Database.Host, ")/", configs.AppConfig.Database.Name, "?parseTime=true")
	// log.Info(mysqlHost)
	db, err = gorm.Open("mysql", mysqlHost)

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Panic(err.Error())
	}

	//set limit
	//db.DB().SetConnMaxLifetime(5 * time.Minute)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(40)

	return db
}

//GetSharedConnection return the database connection
func GetSharedConnection() *gorm.DB {
	return db
}
