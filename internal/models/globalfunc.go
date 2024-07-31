package models

import (
	"encoding/json"
	"fmt"
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
		return fInfo, fmt.Errorf("GetStructInfoPK:" + err.Error())
	}
	// Получаем FK
	fInfo, err = structInfo.GetForeignKey()
	if err != nil {
		return fInfo, fmt.Errorf("GetStructInfoPK:" + err.Error())
	}

	return fInfo, err
}
