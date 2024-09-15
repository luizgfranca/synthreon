package model

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
