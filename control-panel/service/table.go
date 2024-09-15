package service

import (
	"database/sql"
	"platformlab/controlpanel/model"

	"gorm.io/gorm"
)

type ErrorMessage struct {
	Message string
}

type Table struct {
	Db *gorm.DB
}

func (*Table) GetDatabaseTables() ([]string, error) {
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

func (*Table) GetTableColumns(table string) ([]model.TableColumn, error) {
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

	list := []model.TableColumn{}
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

		column := model.TableColumn{
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

func (t *Table) GetTablesMetadata() (*[]model.TableInfo, error) {
	tables, err := t.GetDatabaseTables()
	if err != nil {
		return nil, err
	}

	tableInfoList := []model.TableInfo{}
	for _, table := range tables {
		columns, err := t.GetTableColumns(table)
		if err != nil {
			return nil, err
		}

		info := model.TableInfo{
			Name:       table,
			ColumnInfo: columns,
		}

		tableInfoList = append(tableInfoList, info)
	}

	return &tableInfoList, nil
}
