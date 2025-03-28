package parser

import (
	"database/sql/driver"
	"reflect"
	"slices"
	"strings"
)

func StructsToTableName(rft reflect.Type) (tab string) {
	if field, ok := rft.FieldByName("TableName"); ok {
		if field.Tag.Get("db") != "" {
			tab = field.Tag.Get("db")
		}
	}
	if tab == "" {
		if tn := reflect.New(rft).Elem().MethodByName("TableName"); tn.IsValid() {
			tab = tn.Call(nil)[0].String()
		}
	}
	if tab == "" {
		tab = rft.Name()
	}
	return
}

func StructsParse(obj any) (FieldTag []string, FieldStruct []string, pkField string) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		return StructsTypeParse(rfv.Type())
	case reflect.Slice:
		return StructsTypeParse(rfv.Type())
	default:
		return
	}
}

func StructsTypeParse(rft reflect.Type) (fieldTag []string, fieldStruct []string, pkField string) {
	//rfv := reflect.Indirect(reflect.ValueOf(obj))
	if rft.Kind() == reflect.Slice {
		rft2 := rft.Elem()
		if rft2.Kind() == reflect.Struct {
			return StructsTypeParse(rft2)
		}
	} else {
		for i := 0; i < rft.NumField(); i++ {
			field := rft.Field(i)
			if field.Anonymous {
				continue
			}
			tag := field.Tag.Get("db")
			if tag == "-" || field.Name == "TableName" {
				continue
			}
			if tag == "" {
				//field.Tag = reflect.StructTag("db:" + field.Name)
				fieldStruct = append(fieldStruct, field.Name)
				tag = field.Name
			} else {
				if strings.Contains(tag, ",") {
					tags := strings.Split(tag, ",")
					if slices.Contains(tags, "pk") {
						pkField = field.Name
						tag = tags[0]
					}
				}
				fieldStruct = append(fieldStruct, field.Name)
			}
			//else {
			//	fieldStruct = append(fieldStruct, field.Tag.Get("db"))
			//}
			//if field.Tag.Get("pkField") == "true" {
			//	pkField = field.Name
			//	pkValue = rfv.FieldByName(field.Name)
			//}
			fieldTag = append(fieldTag, tag)
		}
	}
	return
}

//func StructsToSelects(obj any) []string {
//	tag, fieldStruct, _ := StructsParse(obj)
//	if len(tag) > 0 {
//		return tag
//	} else {
//		return fieldStruct
//	}
//}

func StructDataToMap(rfv reflect.Value, tags, fieldStruct []string, mustColumn ...string) (data map[string]any, err error) {
	data = make(map[string]any)
	for i, fieldName := range fieldStruct {
		field := rfv.FieldByName(fieldName)
		if (field.Kind() == reflect.Ptr && field.IsNil()) || (field.IsZero() && !slices.Contains(mustColumn, tags[i])) {
			continue
		}
		var rfvVal = field.Interface()
		if v, ok := rfvVal.(driver.Valuer); ok {
			var value driver.Value
			value, err = v.Value()
			if err != nil {
				return
			}
			data[tags[i]] = value
		} else {
			data[tags[i]] = rfvVal
		}
	}
	return
}

func structUpdateDataToMap(rfv reflect.Value, tags, fieldStruct []string, pkField string, mustColumn ...string) (data map[string]any, err error) {
	data = make(map[string]any)
	for i, fieldName := range fieldStruct {
		field := rfv.FieldByName(fieldName)
		if (field.Kind() == reflect.Ptr && field.IsNil()) || (field.IsZero() && !slices.Contains(mustColumn, tags[i])) || fieldName == pkField {
			continue
		}
		var rfvVal = field.Interface()
		if v, ok := rfvVal.(driver.Valuer); ok {
			var value driver.Value
			value, err = v.Value()
			if err != nil {
				return
			}
			data[tags[i]] = value
		} else {
			data[tags[i]] = rfvVal
		}
	}
	return
}

func StructToDelete(obj any, mustColumn ...string) (data map[string]any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	if rfv.Kind() == reflect.Struct {
		tag, fieldStruct, _ := StructsTypeParse(rfv.Type())
		data, err = StructDataToMap(rfv, tag, fieldStruct, mustColumn...)
	}
	return
}

func StructsToInsert(obj any, mustColumn ...string) (datas []map[string]any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		fieldTag, fieldStruct, _ := StructsTypeParse(rfv.Type())
		var data = make(map[string]any)
		data, err = StructDataToMap(rfv, fieldTag, fieldStruct, mustColumn...)
		if err != nil {
			return
		}
		datas = append(datas, data)
	case reflect.Slice:
		tag, fieldStruct, _ := StructsTypeParse(rfv.Type())
		for i := 0; i < rfv.Len(); i++ {
			var data = make(map[string]any)
			data, err = StructDataToMap(rfv.Index(i), tag, fieldStruct, mustColumn...)
			if err != nil {
				return
			}
			datas = append(datas, data)
		}
	default:
		return
	}
	return
}

func StructToUpdate(obj any, mustColumn ...string) (data map[string]any, pkTag string, pkValue any, err error) {
	tag, fieldStruct, pkField := StructsParse(obj)
	if len(tag) > 0 {
		data = make(map[string]any)
		rfv := reflect.Indirect(reflect.ValueOf(obj))
		data, err = structUpdateDataToMap(rfv, tag, fieldStruct, pkField, mustColumn...)
		if err != nil {
			return
		}
		if pkField != "" {
			pkTag = tag[slices.Index(fieldStruct, pkField)]
			pkValue = rfv.FieldByName(pkField).Interface()
		}
	}

	return
}
