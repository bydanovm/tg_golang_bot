package database

import (
	"database/sql"
	"fmt"

	"github.com/mbydanov/tg_golang_bot/internal/models"
)

// Передавать данные в json структуре для сериализации

func WriteRecord[T any](in []byte) ([]byte, error) {
	// Десереализация
	data, err := models.UnmarshalJSON[T](in)
	if err != nil {
		return nil, fmt.Errorf("WriteRecord:" + err.Error())
	}
	// Запись в БД
	if err := WriteDataT(data); err != nil {
		return nil, fmt.Errorf("WriteRecord:" + err.Error())
	}
	// Сериализация обратно
	buffer, err := models.MarshalJSON(data)
	if err != nil {
		return nil, fmt.Errorf("WriteRecord:" + err.Error())
	}
	return buffer, nil
}

func WriteDataT[T any](in T) error {

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
		return fmt.Errorf("WriteData:" + sqlExecErr + ":" + err.Error())
	}

	return nil
}
