package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/mbydanov/tg_golang_bot/internal/models"
)

// Передавать данные в json структуре для сериализации
func CheckRecord[T any](in []byte) ([]byte, error) {
	// Десереализация
	data, err := models.UnmarshalJSON[T](in)
	if err != nil {
		return nil, fmt.Errorf("CheckRecord:" + err.Error())
	}
	// Проверка наличия записи в БД по PK
	if err := checkRecordPK(data); err != nil {
		return nil, fmt.Errorf("CheckRecord:" + err.Error())
	}
	// Сериализация обратно
	buffer, err := models.MarshalJSON(data)
	if err != nil {
		return nil, fmt.Errorf("CheckRecord:" + err.Error())
	}
	return buffer, nil
}
func WriteRecord[T any](in []byte) ([]byte, interface{}, error) {
	// Десереализация
	data, err := models.UnmarshalJSON[T](in)
	if err != nil {
		return nil, -1, fmt.Errorf("WriteRecord:" + err.Error())
	}
	// Запись в БД
	id, err := writeDataT(data)
	if err != nil {
		return nil, -1, fmt.Errorf("WriteRecord:" + err.Error())
	}
	// Сериализация обратно
	buffer, err := models.MarshalJSON(data)
	if err != nil {
		return nil, -1, fmt.Errorf("WriteRecord:" + err.Error())
	}
	return buffer, id, nil
}
func UpdateRecord[T any](in []byte) ([]byte, error) {
	// Десереализация
	data, err := models.UnmarshalJSON[T](in)
	if err != nil {
		return nil, fmt.Errorf("UpdateRecord:" + err.Error())
	}
	// Запись в БД
	if err := updateDataT(data); err != nil {
		return nil, fmt.Errorf("UpdateRecord:" + err.Error())
	}
	// Сериализация обратно
	buffer, err := models.MarshalJSON(data)
	if err != nil {
		return nil, fmt.Errorf("UpdateRecord:" + err.Error())
	}
	return buffer, nil
}
func checkRecordPK[T any](in T) error {
	// Определяем информацию по структуре
	structInfo := models.StructInfo{}
	if err := structInfo.GetFieldInfo(in); err != nil {
		return fmt.Errorf("CheckDataT:" + err.Error())
	}
	// Получаем PK
	fieldInfo, err := structInfo.GetPrimaryKey()
	if err != nil {
		return fmt.Errorf("CheckDataT:" + err.Error())
	}

	// Создаем SQL запрос
	data := `SELECT COUNT(*) FROM ` + structInfo.StructName + ` WHERE ` +
		fieldInfo.StructNameFields + ` = '` + fmt.Sprintf("%v", fieldInfo.StructValue) + `';`

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return fmt.Errorf("CheckDataT:Open:" + err.Error())
	}
	defer db.Close()

	//Выполняем наш SQL запрос
	var count int
	err = db.QueryRow(data).Scan(&count)
	if err != nil || count == 0 {
		if err == sql.ErrNoRows {
			return fmt.Errorf("QueryRow:Scan:NoRows" + data + ":" + err.Error())
		}
		return fmt.Errorf("QueryRow:Scan:" + data + ":" + func(in error) (out string) {
			if err != nil {
				out = err.Error()
			} else if count == 0 {
				out = "NoRows"
			} else {
				out = "OtherError"
			}
			return out
		}(err))
	}

	return nil
}
func writeDataT[T any](in T) (res interface{}, err error) {

	// Определяем информацию по структуре
	structInfo := models.StructInfo{}
	if err := structInfo.GetFieldInfo(in); err != nil {
		return -1, err
	}
	fieldinfoMap, err := structInfo.UnionFieldsSQL()
	if err != nil {
		return -1, err
	}

	// Создаем SQL запрос
	data := `INSERT INTO ` + structInfo.StructName + ` (` +
		fieldinfoMap["Fields"] + `) VALUES ('` + fieldinfoMap["Values"] + `')`

	fieldValue, err := structInfo.GetIncrement()
	if fieldValue.StructNameFields != "" {
		data += ` RETURNING ` + fieldValue.StructNameFields
	}
	data += `;`

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return -1, fmt.Errorf("WriteData:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	//Выполняем наш SQL запрос и получаем
	// sr, err := db.Exec(data)
	err = db.QueryRow(data).Scan(&res)
	if err != nil {
		return -1, fmt.Errorf("WriteData:" + data + ":" + err.Error())
	}

	return res, nil
}

func updateDataT[T any](in T) error {

	// Определяем информацию по структуре
	structInfo := models.StructInfo{}
	if err := structInfo.GetFieldInfo(in); err != nil {
		return err
	}

	// Получаем PK
	fieldInfo, err := structInfo.GetPrimaryKey()
	if err != nil {
		return fmt.Errorf("updateDataT:" + err.Error())
	}

	// Создаем SQL запрос
	data := `UPDATE ` + structInfo.StructName + ` SET `
	for _, value := range structInfo.StructFieldInfo {
		data += value.StructNameFields + " = '" + fmt.Sprintf("%v",
			func(in interface{}) (out interface{}) {
				switch inConv := in.(type) {
				case time.Time:
					out, _ = models.ConvertDateTimeToMSK(inConv)
				default:
					out = inConv
				}
				return out
			}(value.StructValue)) + "', "
	}
	data = data[:len(data)-2] + " WHERE "

	data += fieldInfo.StructNameFields + ` = '` + fmt.Sprintf("%v", fieldInfo.StructValue) + `';`

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return fmt.Errorf("updateDataT:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data); err != nil {
		return fmt.Errorf("updateDataT:" + data + ":" + err.Error())
	}

	return nil
}

// Функция чтения из БД
// Входные данные:
// - таблица
// - отображаемые поля
// - выражения
// * нужно добавить поддержку сортировки и группировки (через интерфейсы?)
// Выходные данные: массив-интерфейс (структура), ошибка
func ReadData[T any](fields T, expression []Expressions, countIter int) (map[int]T, int, error) {
	returnValues := make(map[int]T)
	cntIter := 0
	var str string
	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, 0,
			fmt.Errorf("ReadDataRow:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	columnsPtr := getFields(fields)
	// Опредение имени колонок
	tableName, columns, _, err := getFieldsName(fields)
	if err != nil {
		return nil, 0,
			fmt.Errorf("ReadDataRow:" + err.Error())
	}
	//Создаем SQL запрос
	data := `SELECT ` + strings.Join(columns, ", ") + ` FROM ` + tableName + ` WHERE `
	for _, value := range expression {
		str += value.Join()
	}
	data += str[:len(str)-4] + `;`

	rows, err := db.Query(data)
	if err != nil {
		return nil, 0,
			fmt.Errorf("ReadDataRow:" + sqlExecErr + ":" + err.Error())
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(columnsPtr...)
		if err != nil {
			return nil, 0,
				fmt.Errorf("ReadDataRow:" + sqlScanErr + ":" + err.Error())
		}

		returnValue := clone(fields)

		structInfo := models.StructInfo{}
		if err := structInfo.GetFieldInfo(returnValue); err != nil {
			return nil, 0, fmt.Errorf("ReadDataRow:" + err.Error())
		}
		// Получаем PK
		fieldValue, err := structInfo.GetPrimaryKey()
		if err != nil {
			return nil, 0, fmt.Errorf("ReadDataRow:" + err.Error())
		}

		key, ok := fieldValue.StructValue.(int)
		if !ok {
			return nil, cntIter,
				fmt.Errorf("ReadDataRow:Key type assertion error")
		}
		rValue, ok := returnValue.(T)
		if !ok {
			return nil, cntIter,
				fmt.Errorf("ReadDataRow:Value type assertion error")
		}
		returnValues[key] = rValue

		cntIter++
		if countIter == cntIter {
			break
		}
	}
	if err = rows.Err(); err != nil {
		return nil, cntIter,
			fmt.Errorf("ReadDataRow:" + sqlSomeOneErr + ":" + err.Error())
	}
	if cntIter > 0 {
		return returnValues, cntIter, nil
	}
	return returnValues, 0, nil
}
