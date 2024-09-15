package service

import (
	"platformlab/controlpanel/model"

	"gorm.io/gorm"
)

type Project struct {
	Db *gorm.DB
}

func (p *Project) FindAll() *[]model.Project {
	var projects []model.Project

	p.Db.Find(&projects)

	return &projects
}
