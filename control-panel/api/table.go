package api

import (
	"encoding/json"
	"net/http"
	"platformlab/controlpanel/service"

	"gorm.io/gorm"
)

type Table struct {
	service service.Table
}

func (t *Table) GetTablesMetadata() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metadata, err := t.service.GetTablesMetadata()
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
		}

		json.NewEncoder(w).Encode(metadata)
	}
}

func TableRESTApi(db *gorm.DB) *Table {
	service := service.Table{Db: db}
	return &Table{service}
}
