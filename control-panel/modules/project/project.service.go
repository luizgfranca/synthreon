package projectmodule

import (
	"fmt"
	commonmodule "platformlab/controlpanel/modules/common"
	toolmodule "platformlab/controlpanel/modules/tool"

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

func (p *ProjectService) FindByAcronym(acronym string) (*Project, error) {
	var maybeProject *Project

	result := p.Db.Where("acronym = ?", acronym).First(&maybeProject)
	if result.Error != nil {
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s not found", acronym),
		}
	}

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

	result := t.Db.Find(&tools).Where("project_id = ?", project.ID)
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
