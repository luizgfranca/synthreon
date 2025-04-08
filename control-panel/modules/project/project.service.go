package projectmodule

import (
	"fmt"
	"log"
	commonmodule "synthreon/modules/common"
	toolmodule "synthreon/modules/tool"

	"gorm.io/gorm"
)

type ProjectService struct {
	Db *gorm.DB
}

func (p *ProjectService) FindAll() *[]Project {
	var projects []Project

	result := p.Db.Find(&projects)
	if result.Error != nil {
		panic(fmt.Sprintf("unable to query database: %s", result.Error.Error()))
	}

	return &projects
}

func (p *ProjectService) FindById(id uint) (*Project, error) {
	var maybeProject *Project

	result := p.Db.Where("id = ?", id).First(&maybeProject)
	if result.Error != nil {
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("element with id %d not found", id),
		}
	}

	return maybeProject, nil
}

func (p *ProjectService) FindByAcronym(acronym string) (*Project, error) {
	var maybeProject *Project

	log.Println("[ProjectService] looking for project with acronym:", acronym)
	result := p.Db.Where("acronym = ?", acronym).First(&maybeProject)
	if result.Error != nil {
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s not found", acronym),
		}
	}

	log.Println("[ProjectService] project found:", maybeProject.Acronym, maybeProject)

	return maybeProject, nil
}

func (p *ProjectService) Create(project *Project) (*Project, error) {
	var result *gorm.DB

	_, err := p.FindByAcronym(project.Acronym)
	if err == nil {
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s already exists", project.Acronym),
		}
	}

	result = p.Db.Create(project)
	if result.Error != nil {
		return nil, result.Error
	}

	var created *Project
	result = p.Db.Where("acronym = ?", project.Acronym).First(&created)
	if result.Error != nil {
		return nil, result.Error
	}
	if created == nil {
		panic("created item in database, but it was not found after insertion")
	}

	return created, nil
}

func (t *ProjectService) FindTools(project *Project) *[]toolmodule.Tool {
	var tools []toolmodule.Tool

	log.Println("[ProjectService] finding tools for project", project)

	result := t.Db.Where("project_id = ?", project.ID).Find(&tools)
	if result.Error != nil {
		panic(fmt.Sprintf("unable to query database: %s", result.Error.Error()))
	}

	return &tools
}

func (t *ProjectService) FindToolByAcronym(project *Project, acronym string) (*toolmodule.Tool, error) {
	var tool *toolmodule.Tool

	result := t.Db.Where("project_id = ? and acronym = ?", project.ID, acronym).First(&tool)
	if result.Error != nil {
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s not found", acronym),
		}
	}

	return tool, nil
}
