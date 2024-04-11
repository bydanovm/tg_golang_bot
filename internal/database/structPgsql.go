package database

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	sqlConErr       string = "SQL error connection"
	sqlExecErr      string = "SQL error exec query"
	sqlScanErr      string = "SQL error scan"
	sqlSomeOneErr   string = "SQL error"
	EQ              string = "="
	NotEQ           string = "!="
	Empty           string = ""
	Id              string = "id"
	Timestamp       string = "timestamp"
	CryptoId        string = "cryptoid"
	CryptoName      string = "cryptoname"
	CryptoLastPrice string = "cryptolastorice"
	CryptoUpdate    string = "cryptoupdate"
	Name            string = "name"
	Description     string = "description"
	Active          string = "active"
	Type            string = "type"
	Value           string = "value"
	Timestart       string = "timestart"
	Timelast        string = "timelast"
)

type LogMsg struct {
	Id        int       `sql_type:"SERIAL PRIMARY KEY"`
	Timestamp time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	UserName  string    `sql_type:"TEXT"`
	Chat_Id   int       `sql_type:"INTEGER"`
	Message   string    `sql_type:"TEXT"`
	Answer    string    `sql_type:"TEXT"`
}

type DictCrypto struct {
	Id              int       `sql_type:"SERIAL PRIMARY KEY"`
	Timestamp       time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	CryptoId        int       `sql_type:"INTEGER"`
	CryptoName      string    `sql_type:"TEXT"`
	CryptoLastPrice float32   `sql_type:"NUMERIC(15,9)"`
	CryptoUpdate    time.Time `sql_type:"TIMESTAMP"`
	Active          bool      `sql_type:"BOOLEAN NOT NULL DEFAULT TRUE"`
	CryptoCounter   int       `sql_type:"INTEGER NOT NULL DEFAULT 0"`
}

// Структура данных таблицы Cryptoprices
type Cryptoprices struct {
	Id           int       `sql_type:"SERIAL PRIMARY KEY"`
	Timestamp    time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	CryptoId     int       `sql_type:"INTEGER"`
	CryptoPrice  float32   `sql_type:"NUMERIC(15,9)"`
	CryptoUpdate time.Time `sql_type:"TIMESTAMP"`
}

// Настроечная таблица
type SettingsProject struct {
	Id          int       `sql_type:"SERIAL PRIMARY KEY"`
	Name        string    `sql_type:"TEXT"`
	Description string    `sql_type:"TEXT"`
	Active      bool      `sql_type:"BOOLEAN"`
	Type        string    `sql_type:"TEXT"`
	Value       string    `sql_type:"TEXT"`
	Timestart   time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	Timelast    time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
}
type LevelsSecureAdd struct {
	IdLvlSecAdd     int    `sql_type:"SERIAL PRIMARY KEY"`
	NameLvlSecAdd   string `sql_type:"TEXT"`
	ActiveLvlSecAdd bool   `sql_type:"BOOLEAN DEFAULT FALSE"`
	LvlSecId        int    `sql_type:"INTEGER REFERENCES LevelsSecure (idLvlSec)"`
}
type LevelsSecure struct {
	IdLvlSec   int    `sql_type:"SERIAL PRIMARY KEY"`
	NameLvlSec string `sql_type:"TEXT"`
}
type Groups struct {
	IdGrp    int    `sql_type:"SERIAL PRIMARY KEY"`
	NameGrp  string `sql_type:"TEXT"`
	LvlSecId int    `sql_type:"INTEGER REFERENCES LevelsSecure (idLvlSec)"`
}
type Users struct {
	IdUsr     int       `sql_type:"SERIAL PRIMARY KEY"`
	TsUsr     time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	NameUsr   string    `sql_type:"TEXT NOT NULL"`
	FirstName string    `sql_type:"TEXT NOT NULL"`
	LastName  string    `sql_type:"TEXT NOT NULL"`
	LangCode  string    `sql_type:"TEXT NOT NULL"`
	IsBot     bool      `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
	IsBanned  bool      `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
	ChatIdUsr int64     `sql_type:"NUMERIC(15,0) NOT NULL"`
	IdLvlSec  int       `sql_type:"INTEGER REFERENCES levelssecure (idlvlsec)"`
}

// Поиск пользователя в базе
func (u *Users) Find() (int, error) {
	fields := Users{}
	expLst := []Expressions{}

	expLst = append(expLst, Expressions{
		Key:      "idusr",
		Operator: EQ,
		Value:    `'` + fmt.Sprintf("%v", u.IdUsr) + `'`,
	})

	_, find, _, err := ReadDataRow(&fields, expLst, 1)
	if err != nil {
		return -1, fmt.Errorf("Find:" + err.Error())
	}

	if find {
		return u.IdUsr, nil
	}

	return -1, nil
}
func (u *Users) GetUserName(idUsr int) (string, error) {
	// Проверка на кеш
	if idUsr == u.IdUsr {
		return u.NameUsr, nil
	}
	// Если не закешировано, то новый поиск
	expLst := []Expressions{
		{Key: "idusr",
			Operator: EQ,
			Value:    `'` + fmt.Sprintf("%v", idUsr) + `'`},
	}
	rs, find, _, err := ReadDataRow(&Users{}, expLst, 1)
	if err != nil {
		return "", fmt.Errorf("GetUserName:" + err.Error())
	}

	if find {
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &u)
		}
		return u.NameUsr, nil
	}

	return "", nil
}

func (u *Users) GetChatId(idUsr int) (int64, error) {
	// Проверка на кеш
	if idUsr == u.IdUsr {
		return u.ChatIdUsr, nil
	}
	// Если не закешировано, то новый поиск
	expLst := []Expressions{
		{Key: "idusr",
			Operator: EQ,
			Value:    `'` + fmt.Sprintf("%v", idUsr) + `'`},
	}
	rs, find, _, err := ReadDataRow(&Users{}, expLst, 1)
	if err != nil {
		return 0, fmt.Errorf("GetChatId:" + err.Error())
	}

	if find {
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &u)
		}
		return u.ChatIdUsr, nil
	}
	return 0, nil
}

// Получение First+Last name
func (u *Users) GetFLName(idUsr int) (string, error) {
	// Проверка на кеш
	if idUsr == u.IdUsr {
		return u.FirstName + " " + u.LastName, nil
	}
	// Если не закешировано, то новый поиск
	expLst := []Expressions{
		{Key: "idusr",
			Operator: EQ,
			Value:    `'` + fmt.Sprintf("%v", idUsr) + `'`},
	}
	rs, find, _, err := ReadDataRow(&Users{}, expLst, 1)
	if err != nil {
		return "", fmt.Errorf("GetFLName:" + err.Error())
	}

	if find {
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &u)
		}
		return u.FirstName + " " + u.LastName, nil
	}
	return "", nil
}

// Добавление пользователя в базу
func (u *Users) Add() (int, error) {
	if err := WriteDataStruct(u); err != nil {
		return -1, err
	}
	return u.IdUsr, nil
}

type LimitsDict struct {
	IdLmtDct   int    `sql_type:"SERIAL PRIMARY KEY"`
	NameLmtDct string `sql_type:"TEXT"`
	StdValLmt  int    `sql_type:"INTEGER DEFAULT 0"`
}
type Limits struct {
	IdLmt       int       `sql_type:"SERIAL PRIMARY KEY"`
	ValAvailLmt int       `sql_type:"INTEGER DEFAULT 0"`
	ValUsedLmt  int       `sql_type:"INTEGER DEFAULT 0"`
	ActiveLmt   bool      `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
	TsLmtOn     time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	TsLmtOff    time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	UserId      int       `sql_type:"INTEGER REFERENCES Users (idUsr)"`
}
type TypeTrackingCrypto struct {
	IdTypTrkCrp       int    `sql_type:"SERIAL PRIMARY KEY"`
	NameTypeTrkCrp    string `sql_type:"TEXT NOT NULL UNIQUE"`
	DescTypTrkCrp     string `sql_type:"TEXT NOT NULL"`
	RisingTypTrkCrp   bool   `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
	CalcProcTypTrkCrp bool   `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
}
type TrackingCrypto struct {
	IdTrkCrp    int     `sql_type:"SERIAL PRIMARY KEY"`
	ValTrkCrp   float32 `sql_type:"NUMERIC(19,9)"`
	TypTrkCrpId int     `sql_type:"INTEGER REFERENCES TypeTrackingCrypto (idTypTrkCrp)"`
	DctCrpId    int     `sql_type:"INTEGER REFERENCES DictCrypto (CryptoId)"`
	UserId      int     `sql_type:"INTEGER REFERENCES Users (idUsr)"`
}

func (t *TrackingCrypto) GetTypeInfo() (interface{}, error) {
	// Возможно нужно один раз запустить и держать в кеше
	// Обновлять при обновлении настроек
	expLst := []Expressions{
		{Key: "IdTypTrkCrp", Operator: EQ, Value: fmt.Sprintf("%v", t.IdTrkCrp)},
	}
	rs, find, _, err := ReadDataRow(&TypeTrackingCrypto{}, expLst, 1)
	if err != nil {
		return nil, fmt.Errorf("GetTypeMode:" + err.Error())
	}
	if !find {
		return nil, fmt.Errorf("GetTypeMode:not find type")
	}
	for _, subRs := range rs {
		subFields := TypeTrackingCrypto{}
		mapstructure.Decode(subRs, &subFields)
		return subFields, nil
	}

	return nil, nil
}

func (t *TrackingCrypto) GetTypeForUser() error {

	return nil
}

// Типы не связанные с таблицами
type Expressions struct {
	Key      string
	Operator string
	Value    string
}

func (exp *Expressions) Join() string {
	return fmt.Sprintf("%s %s %s AND ", exp.Key, exp.Operator, exp.Value)
}

func (exp *Expressions) JoinForUpdate() string {
	return fmt.Sprintf("%s = '%s'", exp.Key, exp.Value)
}
