package service

import (
	"fmt"
	"platformlab/controlpanel/model"

	"gorm.io/gorm"
)

type Project struct {
	Db *gorm.DB
}

func (p *Project) FindAll() *[]model.Project {
	var projects []model.Project

	result := p.Db.Find(&projects)
	if result.Error != nil {
		panic(fmt.Sprintf("unable to query database: %s", result.Error.Error()))
	}

	return &projects
}

func (p *Project) Create(project *model.Project) (*model.Project, error) {
	var result *gorm.DB
	var maybeExisting *model.Project

	result = p.Db.Where("acronym = ?", project.Acronym).First(&maybeExisting)
	if result.Error == nil {
		return nil, &model.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s already exists", project.Acronym),
		}
	}

	result = p.Db.Create(project)
	if result.Error != nil {
		return nil, result.Error
	}

	var created *model.Project
	result = p.Db.Where("acronym = ?", project.Acronym).First(&created)
	if result.Error != nil {
		return nil, result.Error
	}
	if created == nil {
		panic("created item in database, but it was not found after insertion")
	}

	return created, nil
}
