package main

import (
	"english/config"
	"english/src/infrastructure/dbconn"
	"english/src/infrastructure/entity"
	"log"
)

func init() {
	config.SetLogger(config.LogFileName())
}

func main() {
	db := dbconn.NewDB()
	defer dbconn.Close(db)

	if err := db.AutoMigrate(&entity.UserEntity{}); err != nil {
		log.Fatalf("Failed to migrate: %s", err)
	}

	log.Println("Migrate successfully")
}
