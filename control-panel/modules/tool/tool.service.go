package toolmodule

import (
	"fmt"
	"log"

	commonmodule "synthreon/modules/common"

	"gorm.io/gorm"
)

type ToolService struct {
	Db *gorm.DB
}

func (t *ToolService) FindAll() *[]Tool {
	var tools []Tool

	result := t.Db.Find(&tools)
	if result.Error != nil {
		panic(fmt.Sprintf("unable to query database: %s", result.Error.Error()))
	}

	return &tools
}

func (t *ToolService) Create(tool *Tool) (*Tool, error) {
	log.Println("toolservice.create")
	var result *gorm.DB
	var maybeExisting *Tool

	// TODO: adding project verification here for now just to
	// 		 avoid false positives, but this verification
	// 		 that requires awareness of project should be
	// 		 transfered to ProjectService in the future, and
	// 		 this function whould keep just the creation logic itself
	result = t.Db.
		Where("acronym = ? and project_id = ?", tool.Acronym, tool.ProjectId).
		First(&maybeExisting)

	if result.Error == nil {
		log.Println("already exists", maybeExisting)
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s already exists", tool.Acronym),
		}
	}

	log.Println("creating tool", tool)
	result = t.Db.Create(tool)
	if result.Error != nil {
		return nil, result.Error
	}

	var created *Tool
	result = t.Db.Where("acronym = ?", tool.Acronym).First(&created)
	if result.Error != nil {
		return nil, result.Error
	}
	if created == nil {
		panic("created item in database, but it was not found after insertion")
	}

	return created, nil
}

func (t *ToolService) FindByAcronym(acronym string) (*Tool, error) {
	var maybeTool *Tool

	result := t.Db.Where("acronym = ?", acronym).First(&maybeTool)
	if result.Error != nil {
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s not found", acronym),
		}
	}

	return maybeTool, nil
}
