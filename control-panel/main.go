package main

import (
	"log"
	"platformlab/controlpanel/application"
	configurationmodule "platformlab/controlpanel/modules/configuration"
	"platformlab/controlpanel/server"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Println("trying to load application configuration")
	configService, err := configurationmodule.TryLoadApplicationConfigFromEnvironment()
	if err != nil {
		panic(err.Error())
	}

	db, err := gorm.Open(sqlite.Open(configService.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed connecting to database")
	}

	application.Setup(configService, db)
	server.StartServer(":8080", configService, db)
}
