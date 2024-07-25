package database

import (
	"database/sql"
	"fmt"

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
func WriteRecord[T any](in []byte) ([]byte, error) {
	// Десереализация
	data, err := models.UnmarshalJSON[T](in)
	if err != nil {
		return nil, fmt.Errorf("WriteRecord:" + err.Error())
	}
	// Запись в БД
	if err := writeDataT(data); err != nil {
		return nil, fmt.Errorf("WriteRecord:" + err.Error())
	}
	// Сериализация обратно
	buffer, err := models.MarshalJSON(data)
	if err != nil {
		return nil, fmt.Errorf("WriteRecord:" + err.Error())
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
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("QueryRow:Scan:NoRows" + data + ":" + err.Error())
		}
		return fmt.Errorf("QueryRow:Scan:" + data + ":" + err.Error())
	}

	return nil
}
func writeDataT[T any](in T) error {

	// Определяем информацию по структуре
	structInfo := models.StructInfo{}
	if err := structInfo.GetFieldInfo(in); err != nil {
		return err
	}
	fieldinfoMap, err := structInfo.UnionFieldsSQL()
	if err != nil {
		return err
	}
	// Создаем SQL запрос
	data := `INSERT INTO ` + structInfo.StructName + ` (` +
		fieldinfoMap["Fields"] + `) VALUES ('` + fieldinfoMap["Values"] + `');`

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return fmt.Errorf("WriteData:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data); err != nil {
		return fmt.Errorf("WriteData:" + data + ":" + err.Error())
	}

	return nil
}
