package otsorm

import "reflect"

type StructScanner map[string][]int

func (_ StructScanner) New(typ reflect.Type) reflect.Value {
	return reflect.New(typ).Elem()
}

func (sc StructScanner) Scan(pks PrimaryKeyCols, cols AttributeCols, dest reflect.Value) error {
	setField := func(index []int, val interface{}) {
		field := dest.FieldByIndex(index)
		switch v := val.(type) {
		case bool:
			field.SetBool(v)
		case int64:
			field.SetInt(v)
		case string:
			field.SetString(v)
		case float64:
			field.SetFloat(v)
		}
	}

	for _, col := range pks {
		if vs, has := sc[col.ColumnName]; has {
			setField(vs, col.Value)
		}
	}

	for _, col := range cols {
		if vs, has := sc[col.ColumnName]; has {
			setField(vs, col.Value)
		}
	}

	return nil
}

func parseStruct(typ reflect.Type) map[string][]int {
	tags := make(map[string][]int, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		switch f.Type.Kind() {
		case reflect.Struct:
			vs := parseStruct(typ.Field(i).Type)
			for k, v := range vs {
				slice := make([]int, len(v)+1)
				slice[0] = i
				tags[k] = slice
				copy(slice[1:], v)
			}
		default:
			if tag := f.Tag.Get("ots"); tag != "" {
				tags[tag] = []int{i}
			}
		}
	}

	return tags
}
