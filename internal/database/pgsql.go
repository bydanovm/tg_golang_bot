package database

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Теперь нам нужно как то взаимодействовать с БД.
// Создаем переменные в которых мы будем хранить данные переменных окружения для подключению к БД.
var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

// Создаем таблицы в БД при подключении к ней
func createTable(tableStruct interface{}) error {
	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return fmt.Errorf("createTable:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()
	tableName, fields, fieldsTypes, err := getFieldsName(tableStruct)
	if err != nil {
		return fmt.Errorf("createTable:" + err.Error())
	}
	tableName = strings.ToLower(tableName)
	if exists, err := checkExistsTable(tableName, "public"); err == nil && !exists {
		var fieldsNameType []string
		for i := 0; i < len(fields); i++ {
			fieldsNameType = append(fieldsNameType, fields[i]+` `+fieldsTypes[i])
		}
		query := `CREATE TABLE ` + tableName + `(` + strings.Join(fieldsNameType, ",") + `);`
		if _, err = db.Exec(query); err != nil {
			return fmt.Errorf("createTable:" + sqlExecErr + ":" + err.Error())
		}
	} else if err != nil {
		return fmt.Errorf("createTable:" + err.Error())
	}
	return nil
}

// Создание таблицы
func CreateTables() error {
	users := Users{}
	dictCrypto := DictCrypto{}
	cryptoPrices := Cryptoprices{}
	settingsProject := SettingsProject{}

	if err := createTable(&users); err != nil {
		return fmt.Errorf("CreateTables:" + err.Error())
	}
	// Справочник криптовалют с последними ценами
	if err := createTable(&dictCrypto); err != nil {
		return fmt.Errorf("CreateTables:" + err.Error())
	}
	// Таблица всех цен по криптовалютам
	if err := createTable(&cryptoPrices); err != nil {
		return fmt.Errorf("CreateTables:" + err.Error())
	}
	// Таблица всех цен по криптовалютам
	if err := createTable(&settingsProject); err != nil {
		return fmt.Errorf("CreateTables:" + err.Error())
	}
	return nil
}

func checkExistsTable(tableName string, tableSchema string) (bool, error) {
	var count uint8
	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return false,
			fmt.Errorf("checkExistsTable:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()
	// Проверяем на существование таблицы в базе
	row := db.QueryRow(`select exists (select * from information_schema.tables where table_name = '` + tableName + `' and table_schema = '` + tableSchema + `')::int as "count";`)
	err = row.Scan(&count)
	if err != nil {
		return false,
			fmt.Errorf("checkExistsTable:" + sqlScanErr + ":" + err.Error())
	}
	if count == 0 {
		return false, nil
	}

	return true, nil
}

// Таблицу мы создали, и нам нужно заносить в нее данные, этим займется следующая функция.
// Собираем данные полученные ботом
func CollectData(username string, chatid int64, message string, answer []string) error {
	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return fmt.Errorf("CollectData:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	//Конвертируем срез с ответом в строку
	answ := strings.Join(answer, ", ")

	//Создаем SQL запрос
	data := `INSERT INTO users(username, chat_id, message, answer) VALUES($1, $2, $3, $4);`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data, `@`+username, chatid, message, answ); err != nil {
		return fmt.Errorf("CollectData:" + sqlExecErr + ":" + err.Error())
	}

	return nil
}
func WriteData(tableName string, Data map[string]string) error {
	var keys, values []string
	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return fmt.Errorf("WriteData:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	//Создаем SQL запрос
	data := `INSERT INTO ` + tableName + ` (`
	for key, value := range Data {
		keys = append(keys, key)
		values = append(values, value)
	}
	keysStr := strings.Join(keys, ", ")
	valuesStr := strings.Join(values, "', '")
	data += keysStr + `) VALUES ('` + valuesStr + `');`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data); err != nil {
		return fmt.Errorf("WriteData:" + sqlExecErr + ":" + err.Error())
	}

	return nil
}
func UpdateData(tableName string, Data map[string]string, expression []Expressions) error {
	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return fmt.Errorf("UpdateData:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	//Создаем SQL запрос
	data := `UPDATE ` + tableName + ` SET `
	for key, value := range Data {
		data += key + " = '" + value + "', "
	}
	data = data[:len(data)-2] + " WHERE "

	// keysStr := strings.Join(keys, ", ")
	// valuesStr := strings.Join(values, "', '")
	// data += keysStr + `) VALUES ('` + valuesStr + `');`

	// var str string
	for _, value := range expression {
		data += value.Join()
	}
	data = data[:len(data)-4] + `;`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data); err != nil {
		return fmt.Errorf("UpdateData:" + sqlExecErr + ":" + err.Error())
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
func ReadDataRow(fields interface{}, expression []Expressions, countIter int) ([]interface{}, bool, int, error) {
	returnValues := []interface{}{}
	cntIter := 0
	var str string
	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, false, 0,
			fmt.Errorf("ReadDataRow:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	columnsPtr := getFields(fields)
	// Опредение имени колонок
	tableName, columns, _, err := getFieldsName(fields)
	if err != nil {
		return nil, false, 0,
			fmt.Errorf("ReadDataRow:" + err.Error())
	}
	// var columnsArr []string
	// for k, _ := range columns {
	// 	columnsArr = append(columnsArr, k)
	// }
	//Создаем SQL запрос
	data := `SELECT ` + strings.Join(columns, ", ") + ` FROM ` + tableName + ` WHERE `
	for _, value := range expression {
		str += value.Join()
	}
	data += str[:len(str)-4] + `;`

	rows, err := db.Query(data)
	if err != nil {
		return nil, false, 0,
			fmt.Errorf("ReadDataRow:" + sqlExecErr + ":" + err.Error())
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(columnsPtr...)
		if err != nil {
			return nil, false, 0,
				fmt.Errorf("ReadDataRow:" + sqlScanErr + ":" + err.Error())
		}

		returnValue := clone(fields)

		returnValues = append(returnValues, returnValue)

		cntIter++
		if countIter == cntIter {
			break
		}
	}
	if err = rows.Err(); err != nil {
		return returnValues, true, cntIter,
			fmt.Errorf("ReadDataRow:" + sqlSomeOneErr + ":" + err.Error())
	}
	if cntIter > 0 {
		return returnValues, true, cntIter, nil
	}
	return nil, false, 0, nil
}

// Функция определения колонок (указателей) в структуре при передаче как интерфейс
func getFields(in interface{}) (out []interface{}) {
	val := reflect.ValueOf(in).Elem()
	// Проверка на поинтер
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	// // Проверка на структуру
	// if val.Kind() != reflect.Struct {
	// 	log.Fatal("неизвестный тип")
	// }

	numCols := val.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := val.Field(i).Addr().Interface()
		columns[i] = field
	}
	return columns
}

// Возврат имени структуры и имен полей
func getFieldsName(in interface{}) (string, []string, []string, error) {
	val := reflect.ValueOf(in).Elem()
	structType := val.Type()
	tableName := structType.Name()
	numCols := structType.NumField()
	columns := make([]string, numCols)
	fieldTypes := make([]string, numCols)
	for i := 0; i < numCols; i++ {
		field := structType.Field(i)
		tag := field.Tag
		fieldName := field.Name
		fieldType := tag.Get("sql_type")
		// columns[fieldName] = fieldType
		columns[i] = fieldName
		fieldTypes[i] = fieldType
		// fieldTypes[i] = fieldType
		// println(fmt.Sprintf("%v", columns[i]))
	}

	return tableName, columns, fieldTypes, nil
}
func clone(inter interface{}) interface{} {
	nInter := reflect.New(reflect.TypeOf(inter).Elem())

	val := reflect.ValueOf(inter).Elem()
	nVal := nInter.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		nvField.Set(val.Field(i))
	}

	return nInter.Interface()
}

// Также давайте напишем функцию которая будет считать количество уникальных пользователей
// которые писали боту, чтобы отдавать это число пользователям если они отправят боту нужную команду.
func GetNumberOfUsers() (int64, error) {
	var count int64

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return 0,
			fmt.Errorf("ReadDataRow:" + sqlConErr + ":" + err.Error())
	}
	defer db.Close()

	// Отправляем запрос в БД для подсчета числа уникальных пользователей
	row := db.QueryRow("SELECT COUNT(DISTINCT username) FROM users;")
	err = row.Scan(&count)
	if err != nil {
		return 0,
			fmt.Errorf("ReadDataRow:" + sqlScanErr + ":" + err.Error())
	}

	return count, nil
}
