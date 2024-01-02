package otsorm

import (
	"errors"
	"fmt"
	"reflect"
)

type Scanner interface {
	New(reflect.Type) reflect.Value
	Scan(PrimaryKeyCols, AttributeCols, reflect.Value) error
}

func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func scanRow(rows IRow, dest reflect.Value) error {
	var scanner Scanner
	dest = indirect(dest)
	switch dest.Kind() {
	case reflect.Map:
		scanner = MapScanner{}
	case reflect.Struct:
		scanner = (StructScanner)(parseStruct(dest.Type()))
	default:
		return fmt.Errorf("unsupport destination type: %v", dest.Type())
	}

	rows.Reset()
	pks, cols, ok := rows.Next()
	if !ok {
		return errors.New("empty rows")
	}

	return scanner.Scan(pks, cols, dest)
}

func ScanRow(rows IRow, dest interface{}) error {
	if rows.Len() == 0 {
		return nil
	}

	return scanRow(rows, reflect.ValueOf(dest))
}

func scanRows(rows IRow, dest reflect.Value) error {
	dest = indirect(dest)
	elemType := dest.Type().Elem()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		elemType = elemType.Elem()
	}
	dest.Set(reflect.MakeSlice(dest.Type(), 0, rows.Len()))

	var scanner Scanner
	switch elemType.Kind() {
	case reflect.Map:
		scanner = MapScanner{}
	case reflect.Struct:
		scanner = (StructScanner)(parseStruct(elemType))
	}

	for i := 0; i < rows.Len(); i++ {
		pks, cols, ok := rows.Next()
		if !ok {
			continue
		}

		obj := scanner.New(elemType)
		err := scanner.Scan(pks, cols, obj)
		if err != nil {
			return err
		}

		if isPtr {
			dest.Set(reflect.Append(dest, obj.Addr()))
		} else {
			dest.Set(reflect.Append(dest, obj))
		}
	}

	return nil
}

func ScanRows(rows IRow, dest interface{}) (err error) {
	if rows.Len() == 0 {
		return
	}

	return scanRows(rows, reflect.ValueOf(dest))
}

func Unmarshal(rows IRow, dest interface{}) (err error) {
	if rows.Len() == 0 {
		return nil
	}

	vv := indirect(reflect.ValueOf(dest))
	switch vv.Kind() {
	case reflect.Slice:
		return scanRows(rows, vv)
	default:
		return scanRow(rows, vv)
	}
}
