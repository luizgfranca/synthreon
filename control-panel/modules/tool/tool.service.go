package toolmodule

import (
	"fmt"
	genericmodule "platformlab/controlpanel/modules/commonmodule"

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
	var result *gorm.DB
	var maybeExisting *Tool

	genericmodule.Probe(tool.Acronym)

	result = t.Db.Where("acronym = ?", tool.Acronym).First(&maybeExisting)
	if result.Error == nil {
		return nil, &genericmodule.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s already exists", tool.Acronym),
		}
	}

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
