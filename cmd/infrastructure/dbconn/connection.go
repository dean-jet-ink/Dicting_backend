package dbconn

import (
	"english/config"
	"log"

	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	myConf := mysql.Config{
		DBName:               config.MySQLDBName(),
		User:                 config.MySQLUser(),
		Passwd:               config.MySQLPass(),
		Addr:                 config.MySQLHost(),
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	dsn := myConf.FormatDSN()

	log.Println(dsn)

	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB: %s", err)
	}

	log.Println("DB connected")

	return db
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB: %s", err)
	}

	sqlDB.Close()
	log.Println("DB closed")
}
