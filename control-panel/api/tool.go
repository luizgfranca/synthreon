package api

import (
	"encoding/json"
	"net/http"
	"platformlab/controlpanel/model"

	"gorm.io/gorm"
)

type Tool struct {
}

func (t *Tool) GetEventRresponseTEST() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input model.ToolEvent
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		if !input.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "invalid request data"})
			return
		}

		json.NewEncoder(w).Encode(input)
	}
}

func ToolRestAPI(db *gorm.DB) *Tool {
	return &Tool{}
}
