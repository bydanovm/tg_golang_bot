package database

import (
	"fmt"
	"time"
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

type Users struct {
	Id        int       `sql_type:"SERIAL PRIMARY KEY"`
	Timestamp time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	UserName  string    `sql_type:"INT"`
	Chat_Id   int       `sql_type:"NUMERIC(15,3)"`
	Message   string    `sql_type:"TEXT"`
	Answer    string    `sql_type:"TEXT"`
}

type DictCrypto struct {
	Id              int       `sql_type:"SERIAL PRIMARY KEY"`
	Timestamp       time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	CryptoId        int       `sql_type:"INT"`
	CryptoName      string    `sql_type:"TEXT"`
	CryptoLastPrice float32   `sql_type:"NUMERIC(15,9)"`
	CryptoUpdate    time.Time `sql_type:"TIMESTAMP"`
}

// Структура данных таблицы Cryptoprices
type Cryptoprices struct {
	Id           int       `sql_type:"SERIAL PRIMARY KEY"`
	Timestamp    time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	CryptoId     int       `sql_type:"INT"`
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
