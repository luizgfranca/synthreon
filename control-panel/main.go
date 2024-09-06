package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"platformlab/controlpanel/component"
	"platformlab/controlpanel/model"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ErrorMessage struct {
	Message string
}

type TableColumn struct {
	Name string
	Type string
}

func CreateMockProjects(db *gorm.DB) {
	db.AutoMigrate(&model.Project{})

	p := []model.Project{
		{Acronym: "dcc", Name: "DCC"},
		{Acronym: "dsi", Name: "DSI"},
		{Acronym: "customer-identity", Name: "Customer Identity"},
	}

	for _, it := range p {
		db.Create(it)
	}
}

func GetAllProjects(db *gorm.DB) *[]model.Project {
	var projects []model.Project

	db.Find(&projects)

	return &projects
}

func GetDatabaseTables() ([]string, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("select name from sqlite_master where type = 'table'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []string{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		list = append(list, name)
	}

	return list, nil
}

func GetTableColumns(table string) ([]TableColumn, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("select name from sqlite_master where type = 'table'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []string{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		list = append(list, name)
	}

	return list, nil
}

func main() {
	router := mux.NewRouter()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed connecting to database")
	}

	CreateMockProjects(db)

	router.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
		var projects = GetAllProjects(db)
		component.ProjectList(*projects).Render(context.Background(), w)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := Hello("platformlab")
		component.Render(context.Background(), w)
	})

	router.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		tables, err := GetDatabaseTables()
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{err.Error()})
		}
		json.NewEncoder(w).Encode(tables)
	})

	http.ListenAndServe(":8080", router)
}
