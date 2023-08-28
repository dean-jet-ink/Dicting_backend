package main

import (
	"english/cmd/infrastructure/dbconn"
	"english/cmd/infrastructure/entity"
	"english/config"
	"log"
)

func init() {
	config.SetLogger(config.LogFileName())
}

func main() {
	db := dbconn.NewDB()
	defer dbconn.Close(db)

	if err := db.AutoMigrate(&entity.UserEntity{}, &entity.EnglishItemEntity{}); err != nil {
		log.Fatalf("Failed to migrate: %s", err)
	}

	log.Println("Migrate successfully")
}
