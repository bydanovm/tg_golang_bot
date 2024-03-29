package models

import (
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
