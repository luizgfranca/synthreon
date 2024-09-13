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
	Id           int
	Name         string
	Type         string
	NotNull      bool
	DefaultValue string
	IsPrimaryKey bool
}

type TableInfo struct {
	Name       string
	ColumnInfo []TableColumn
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

	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []TableColumn{}
	var cid int
	var cname string
	var ctype string
	var cnotnull bool
	var dfltVal sql.NullString
	var primary bool

	for rows.Next() {
		err = rows.Scan(
			&cid,
			&cname,
			&ctype,
			&cnotnull,
			&dfltVal,
			&primary)

		if err != nil {
			return nil, err
		}

		column := TableColumn{
			Id:           cid,
			Name:         cname,
			Type:         ctype,
			NotNull:      cnotnull,
			DefaultValue: dfltVal.String,
			IsPrimaryKey: primary,
		}

		list = append(list, column)
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
		http.Redirect(w, r, "/control-panel", http.StatusFound)
	})

	router.PathPrefix("/control-panel").Handler(http.StripPrefix("/control-panel", http.FileServer(http.Dir("./web/dist"))))
	router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./web/dist")))
	// router.Handle("/assets", http.FileServer(http.Dir("./web/dist/assets")))

	router.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		tables, err := GetDatabaseTables()
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{err.Error()})
			return
		}

		tableInfoList := []TableInfo{}
		for _, table := range tables {
			columns, err := GetTableColumns(table)
			if err != nil {
				json.NewEncoder(w).Encode(ErrorMessage{err.Error()})
				return
			}

			info := TableInfo{
				Name:       table,
				ColumnInfo: columns,
			}

			tableInfoList = append(tableInfoList, info)
		}

		json.NewEncoder(w).Encode(tableInfoList)
	})

	http.ListenAndServe(":8080", router)
}
