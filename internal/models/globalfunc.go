package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Функция проверяет массив на пустые ячейки и удаляет их
func ChkArrayBySpace(array []string) []string {
	var tmpArray []string
	for _, v := range array {
		if len(v) != 0 {
			tmpArray = append(tmpArray, v)
			// array = append(array[:k], array[k+1:]...)
		}
	}
	return tmpArray
}

// Поиск значения в массиве и его удаление
func FindCellAndDelete(array []string, findValue string) []string {
	for k, v := range array {
		if v == findValue {
			array = append(array[:k], array[k+1:]...)
			break
		}
	}
	return array
}

// Функция конвертации времени в МСК часовой пояс
func ConvertDateTimeToMSK(iTime time.Time) (string, error) {
	dateTime := iTime
	dateTimeLocUTC3, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return "", fmt.Errorf("ConvertDateTimeToMSK:Time convert error-" + err.Error())
	}
	return dateTime.In(dateTimeLocUTC3).Format(layout), nil
}

func MarshalJSON[T any](object T) ([]byte, error) {
	buffer, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("marshalJSON:" + err.Error())
	}
	return buffer, nil
}
func UnmarshalJSON[T any](buffer []byte) (res T, err error) {
	err = json.Unmarshal(buffer, &res)
	if err != nil {
		return res, fmt.Errorf("unmarshalJSON:" + err.Error())
	}
	return res, nil
}
func GetStructInfo(in interface{}) (stInfo StructInfo, err error) {
	structInfo := StructInfo{}
	// Определяем информацию по структуре
	if err = structInfo.GetFieldInfo(in); err != nil {
		return stInfo, fmt.Errorf("GetStructInfo:" + err.Error())
	}
	return structInfo, err
}
func GetStructInfoPK(in interface{}) (fInfo fieldInfo, err error) {
	structInfo := StructInfo{}
	// Определяем информацию по структуре
	if err = structInfo.GetFieldInfo(in); err != nil {
		return fInfo, fmt.Errorf("GetStructInfoPK:" + err.Error())
	}
	// Получаем PK
	fInfo, err = structInfo.GetPrimaryKey()
	if err != nil {
		return fInfo, fmt.Errorf("GetStructInfoPK:" + err.Error())
	}

	return fInfo, err
}
func GetStructInfoFK(in interface{}) (fInfo fieldInfo, err error) {
	structInfo := StructInfo{}
	// Определяем информацию по структуре
	if err = structInfo.GetFieldInfo(in); err != nil {
		return fInfo, fmt.Errorf("GetStructInfoFK:" + err.Error())
	}
	// Получаем FK
	fInfo, err = structInfo.GetForeignKey()
	if err != nil {
		return fInfo, fmt.Errorf("GetStructInfoFK:" + err.Error())
	}

	return fInfo, err
}
func GetStructInfoSort(in interface{}) (fInfo fieldInfo, err error) {
	structInfo := StructInfo{}
	// Определяем информацию по структуре
	if err = structInfo.GetFieldInfo(in); err != nil {
		return fInfo, fmt.Errorf("GetStructInfoSort:" + err.Error())
	}
	// Получаем SortKey
	fInfo, err = structInfo.GetSortKey()
	if err != nil {
		return fInfo, fmt.Errorf("GetStructInfoSort:" + err.Error())
	}

	return fInfo, err
}

func FormatFloatToString(number float32) (format string) {
	// Установим формат общей длиной в 7 знаков
	if number >= 100000 {
		format = "%.1f"
	} else if number >= 10000 {
		format = "%.2f"
	} else if number >= 1000 {
		format = "%.3f"
	} else if number >= 100 {
		format = "%.4f"
	} else if number >= 10 {
		format = "%.5f"
	} else {
		format = "%.6f"
	}
	return format
}
func GetName(in interface{}) string {
	if t := reflect.TypeOf(in); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
