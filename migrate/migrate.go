package main

import (
	"english/cmd/infrastructure/dbconn"
	"english/cmd/infrastructure/entity"
	"english/config"
	"log"
)

func init() {
	if config.GoEnv() == "dev" {
		config.SetLogger(config.LogFileName())
	}
}

func main() {
	db := dbconn.NewDB()
	defer dbconn.Close(db)

	if err := db.AutoMigrate(&entity.UserEntity{}, &entity.EnglishItemEntity{}, &entity.ExampleEntity{}, &entity.ImgEntity{}, &entity.OutputEntity{}); err != nil {
		log.Fatalf("Failed to migrate: %s", err)
	}

	log.Println("Migrate successfully")
}
