package models

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	RetrieverCoins string = `RetrieverCoins`
	Notificator    string = `Notificator`
	CoinMarketCap  string = `CoinMarketCap`
)

type StatusRetriever struct {
	MsgError error
}

type StatusChannel struct {
	Module string
	Start  bool
	Stop   bool
	Update bool
	Error  error
	Data   interface{}
}

type fieldInfo struct {
	StructNameFields string
	StructTypes      reflect.Kind
	StructSQLTypes   string
	StructTagIsPKey  string
	StructTagIsFKey  string
	StructTagIsMiss  string
	StructTagIsIncr  string
	StructTagIsSort  string
	StructTagIsExt   string
	StructValue      interface{}
}
type StructInfo struct {
	StructName      string
	StructNumFields int
	StructFieldInfo map[string]fieldInfo
}

func (s *StructInfo) GetFieldInfo(in interface{}) error {
	var val reflect.Value
	if reflect.TypeOf(in).Kind() == reflect.Struct {
		val = reflect.ValueOf(in)
	} else {
		val = reflect.ValueOf(in).Elem()
	}
	fieldInfoMap := make(map[string]fieldInfo)
	structType := val.Type()

	s.StructName = structType.Name()
	s.StructNumFields = structType.NumField()

	for i := 0; i < s.StructNumFields; i++ {
		field := structType.Field(i)
		structValue := reflect.ValueOf(val.Field(i).Interface()).Interface()
		tag := field.Tag
		fieldInfoMap[field.Name] = fieldInfo{
			StructNameFields: field.Name,
			StructTypes:      field.Type.Kind(),
			StructSQLTypes:   tag.Get("sql_type"),
			StructTagIsPKey:  tag.Get("pkey"),
			StructTagIsFKey:  tag.Get("fkey"),
			StructTagIsMiss:  tag.Get("miss"),
			StructTagIsIncr:  tag.Get("incr"),
			StructTagIsSort:  tag.Get("sortkey"),
			StructTagIsExt:   tag.Get("ext"),
			StructValue:      structValue}
	}
	s.StructFieldInfo = fieldInfoMap
	return nil
}

// Объединить инфо полей в строки для SQL запроса
func (s *StructInfo) UnionFieldsSQL() (map[string]string, error) {
	var fld, val []string
	fieldInfoMap := make(map[string]string)
	for _, v := range s.StructFieldInfo {
		// Пропустить это поле
		if v.StructTagIsMiss == "YES" {
			continue
		}
		fld = append(fld, v.StructNameFields)
		_, ok := v.StructValue.(time.Time)
		if ok {
			v.StructValue, _ = ConvertDateTimeToMSK(v.StructValue.(time.Time))
		}
		val = append(val, fmt.Sprintf("%v", v.StructValue))
	}

	fieldInfoMap["Fields"] = strings.Join(fld, ",")
	fieldInfoMap["Values"] = strings.Join(val, "','")

	return fieldInfoMap, nil
}

// Получить поле с PK для проверки в БД
func (s *StructInfo) GetPrimaryKey() (field fieldInfo, err error) {
	for _, v := range s.StructFieldInfo {
		if v.StructTagIsPKey == "YES" {
			field = v
		}
	}
	if field.StructNameFields == "" {
		err = fmt.Errorf("GetPrimaryKey:PK is not found")
	}

	return field, err
}

// Получить поле с FK для проверки в БД
func (s *StructInfo) GetForeignKey() (field fieldInfo, err error) {
	for _, v := range s.StructFieldInfo {
		if v.StructTagIsFKey == "YES" {
			field = v
		}
	}
	if field.StructNameFields == "" {
		err = fmt.Errorf("GetForeignKey:FK is not found")
	}

	return field, err
}

// Получить поле autoincrement
func (s *StructInfo) GetIncrement() (field fieldInfo, err error) {
	for _, v := range s.StructFieldInfo {
		if v.StructTagIsIncr == "YES" {
			field = v
			break
		}
	}
	if field.StructNameFields == "" {
		err = fmt.Errorf("GetIncrement:Incr is not found")
	}

	return field, err
}

// Получить доп. поле для сортировки
func (s *StructInfo) GetSortKey() (field fieldInfo, err error) {
	for _, v := range s.StructFieldInfo {
		if v.StructTagIsSort == "YES" {
			field = v
			break
		}
	}
	if field.StructNameFields == "" {
		err = fmt.Errorf("GetSortKey:Sort is not found")
	}

	return field, err
}
