package core

import (
	"reflect"
	"strings"
)

// database table
type Table struct {
	Name          string
	Type          reflect.Type
	columnsSeq    []string
	columnsMap    map[string][]*Column
	columns       []*Column
	Indexes       map[string]*Index
	PrimaryKeys   []string
	AutoIncrement string
	Created       map[string]bool
	Updated       string
	Version       string
	Cacher        Cacher
	StoreEngine   string
	Charset       string
	checkedName   string
}

func (table *Table) CheckedName(d Dialect) string {
	if len(table.checkedName) == 0 {
		table.checkedName = d.CheckedQuote(table.Name)
	}
	return table.checkedName
}

func (table *Table) Columns() []*Column {
	return table.columns
}

func (table *Table) ColumnsSeq() []string {
	return table.columnsSeq
}

func NewEmptyTable() *Table {
	return newTable("", nil)
}

func newTable(name string, t reflect.Type) *Table {
	return &Table{Name: name, Type: t,
		columnsSeq:  make([]string, 0),
		columns:     make([]*Column, 0),
		columnsMap:  make(map[string][]*Column),
		Indexes:     make(map[string]*Index),
		Created:     make(map[string]bool),
		PrimaryKeys: make([]string, 0),
	}
}

func (table *Table) GetColumn(name string) *Column {
	if c, ok := table.columnsMap[strings.ToLower(name)]; ok {
		return c[0]
	}
	return nil
}

func (table *Table) GetColumnIdx(name string, idx int) *Column {
	if c, ok := table.columnsMap[strings.ToLower(name)]; ok {
		if idx < len(c) {
			return c[idx]
		}
	}
	return nil
}

// if has primary key, return column
func (table *Table) PKColumns() []*Column {
	columns := make([]*Column, 0)
	for _, name := range table.PrimaryKeys {
		columns = append(columns, table.GetColumn(name))
	}
	return columns
}

func (table *Table) AutoIncrColumn() *Column {
	return table.GetColumn(table.AutoIncrement)
}

func (table *Table) VersionColumn() *Column {
	return table.GetColumn(table.Version)
}

func (table *Table) UpdatedColumn() *Column {
	return table.GetColumn(table.Updated)
}

// add a column to table
func (table *Table) AddColumn(col *Column) {
	table.columnsSeq = append(table.columnsSeq, col.Name)
	table.columns = append(table.columns, col)
	colName := strings.ToLower(col.Name)
	if c, ok := table.columnsMap[colName]; ok {
		table.columnsMap[colName] = append(c, col)
	} else {
		table.columnsMap[colName] = []*Column{col}
	}

	if col.IsPrimaryKey {
		table.PrimaryKeys = append(table.PrimaryKeys, col.Name)
	}
	if col.IsAutoIncrement {
		table.AutoIncrement = col.Name
	}
	if col.IsCreated {
		table.Created[col.Name] = true
	}
	if col.IsUpdated {
		table.Updated = col.Name
	}
	if col.IsVersion {
		table.Version = col.Name
	}
}

// add an index or an unique to table
func (table *Table) AddIndex(index *Index) {
	table.Indexes[index.Name] = index
}
