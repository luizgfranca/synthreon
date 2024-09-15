package service

import (
	"database/sql"

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

func (*Table) GetTableColumns(table string) ([]TableColumn, error) {
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

func (t *Table) GetTablesMetadata() (*[]TableInfo, error) {
	tables, err := t.GetDatabaseTables()
	if err != nil {
		return nil, err
	}

	tableInfoList := []TableInfo{}
	for _, table := range tables {
		columns, err := t.GetTableColumns(table)
		if err != nil {
			return nil, err
		}

		info := TableInfo{
			Name:       table,
			ColumnInfo: columns,
		}

		tableInfoList = append(tableInfoList, info)
	}

	return &tableInfoList, nil
}
