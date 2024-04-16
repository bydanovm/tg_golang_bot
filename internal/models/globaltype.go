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
	StructValue      interface{}
}
type StructInfo struct {
	StructName      string
	StructNumFields int
	StructFieldInfo map[string]fieldInfo
}

func (s *StructInfo) GetFieldInfo(in interface{}) error {
	val := reflect.ValueOf(in).Elem()
	fieldInfoMap := make(map[string]fieldInfo)
	structType := val.Type()

	s.StructName = structType.Name()
	s.StructNumFields = structType.NumField()

	for i := 0; i < s.StructNumFields; i++ {
		field := structType.Field(i)
		structValue := reflect.ValueOf(val.Field(i).Interface()).Interface()
		tag := field.Tag
		fieldInfoMap[field.Name] = fieldInfo{field.Name, field.Type.Kind(), tag.Get("sql_type"), structValue}
	}
	s.StructFieldInfo = fieldInfoMap
	return nil
}

// Объединить инфо полей в строки для SQL запроса
func (s *StructInfo) UnionFieldsSQL() (map[string]string, error) {
	var fld, val []string
	fieldInfoMap := make(map[string]string)
	for _, v := range s.StructFieldInfo {
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
