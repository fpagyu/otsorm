package otsorm

import "reflect"

type MapScanner struct {
}

func (_ MapScanner) New(typ reflect.Type) reflect.Value {
	return reflect.MakeMap(typ)
}

func (_ MapScanner) Scan(pks PrimaryKeyCols, cols AttributeCols, dest reflect.Value) error {
	for _, col := range pks {
		dest.SetMapIndex(reflect.ValueOf(col.ColumnName),
			reflect.ValueOf(col.Value))
	}

	for _, col := range cols {
		dest.SetMapIndex(reflect.ValueOf(col.ColumnName),
			reflect.ValueOf(col.Value))
	}
	return nil
}
