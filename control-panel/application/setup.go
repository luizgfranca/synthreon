package application

import (
	"log"

	"gorm.io/gorm"

	configurationmodule "platformlab/controlpanel/modules/configuration"
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	usermodule "platformlab/controlpanel/modules/user"
)

func doMigrations(db *gorm.DB) {
	db.AutoMigrate(&projectmodule.Project{})
	db.AutoMigrate(&toolmodule.Tool{})
	db.AutoMigrate(&usermodule.User{})
}

func createExampleProjectsIfNotExists(db *gorm.DB) {
	s := projectmodule.ProjectService{Db: db}

	testProjects := []projectmodule.Project{
		{Acronym: "sandbox", Name: "Sandbox", Description: "Sandbox project to test tool development."},
		{Acronym: "test", Name: "Test", Description: "Test project."},
	}

	for i := range testProjects {
		p := testProjects[i]

		dbProject, _ := s.FindByAcronym(p.Acronym)
		if dbProject == nil {
			log.Println("saving: ", p.Acronym)
			s.Create(&p)
		}
	}
}

func createExampleToolsIfNotExists(db *gorm.DB) {
	ps := projectmodule.ProjectService{Db: db}
	ts := toolmodule.ToolService{Db: db}
	exampleTools := []toolmodule.Tool{
		{ProjectId: 1, Acronym: "sandbox", Description: "Sandbox tool for development testing"},
		{ProjectId: 1, Acronym: "test-a", Description: "Tool A used for testing"},
		{ProjectId: 2, Acronym: "sandbox", Description: "Sandbox tool for test project"},
	}

	for _, t := range exampleTools {
		dbproject, _ := ps.FindById(t.ProjectId)
		if dbproject == nil {
			log.Fatalln("project for tool to be created should already exist")
		}

		dbtool, _ := ps.FindToolByAcronym(dbproject, t.Acronym)
		if dbtool == nil {
			log.Println("saving: ", t.Acronym)
			ts.Create(&t)
		}
	}
}

func createDefaultUserIfNotExists(db *gorm.DB, email string, password string) {
	s := usermodule.UserService{Db: db}

	defaultUser, err := usermodule.NewUser("root", email, password)
	if err != nil {
		panic(err.Error())
	}

	user, _ := s.FindByEmail(defaultUser.Email)
	if user != nil {
		return
	}

	_, err = s.Create(defaultUser)
	if err != nil {
		panic(err.Error())
	}
}

func Setup(configService *configurationmodule.ConfigurationService, db *gorm.DB) {
	log.Println("[Setup] doing database migrations...")
	doMigrations(db)

	log.Println("[Setup] creating example projects...")
	createExampleProjectsIfNotExists(db)

	log.Println("[Setup] creating example tools...")
	createExampleToolsIfNotExists(db)

	log.Println("[Setup] asseting creation of default user...")
	createDefaultUserIfNotExists(db, configService.RootUserEmail, configService.RootPassword)
}
