package main

import (
	"english/config"
	"english/src/infrastructure/db"
	"english/src/infrastructure/entity"
	"log"
)

func init() {
	config.SetLogger(config.LogFileName())
}

func main() {
	conn := db.NewDB()
	defer db.Close(conn)

	if err := conn.AutoMigrate(&entity.UserEntity{}); err != nil {
		log.Fatalf("Failed to migrate: %s", err)
	}

	log.Println("Migrate successfully")
}
