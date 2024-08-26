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
	Id        int       `sql_type:"SERIAL PRIMARY KEY" incr:"YES"`
	Timestamp time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	UserName  string    `sql_type:"TEXT"`
	Chat_Id   int       `sql_type:"INTEGER"`
	Message   string    `sql_type:"TEXT"`
	Answer    string    `sql_type:"TEXT"`
}

type DictCrypto struct {
	Id              int       `sql_type:"SERIAL PRIMARY KEY" incr:"YES"`
	Timestamp       time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	CryptoId        int       `sql_type:"INTEGER" pkey:"YES"`
	CryptoName      string    `sql_type:"TEXT" fkey:"YES"`
	CryptoLastPrice float32   `sql_type:"NUMERIC(15,9)"`
	CryptoUpdate    time.Time `sql_type:"TIMESTAMP"`
	Active          bool      `sql_type:"BOOLEAN NOT NULL DEFAULT TRUE"`
	CryptoCounter   int       `sql_type:"INTEGER NOT NULL DEFAULT 0"`
	CryptoRank      int       `sql_type:"INTEGER NOT NULL DEFAULT 0" sortkey:"YES"`
	CoinMrktId      int       `sql_type:"INTEGER NOT NULL DEFAULT 0"`
}

// Структура данных таблицы Cryptoprices
type Cryptoprices struct {
	Id           int       `sql_type:"SERIAL PRIMARY KEY" miss:"YES" incr:"YES" pkey:"YES"`
	Timestamp    time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	CryptoId     int       `sql_type:"INTEGER" fkey:"YES"`
	CryptoPrice  float32   `sql_type:"NUMERIC(15,9)"`
	CryptoUpdate time.Time `sql_type:"TIMESTAMP"`
}

// Настроечная таблица
type SettingsProject struct {
	Id          int       `sql_type:"SERIAL PRIMARY KEY" incr:"YES"`
	Name        string    `sql_type:"TEXT"`
	Description string    `sql_type:"TEXT"`
	Active      bool      `sql_type:"BOOLEAN"`
	Type        string    `sql_type:"TEXT"`
	Value       string    `sql_type:"TEXT"`
	Timestart   time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	Timelast    time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
}
type LevelsSecureAdd struct {
	IdLvlSecAdd     int    `sql_type:"SERIAL PRIMARY KEY" incr:"YES"`
	NameLvlSecAdd   string `sql_type:"TEXT"`
	ActiveLvlSecAdd bool   `sql_type:"BOOLEAN DEFAULT FALSE"`
	LvlSecId        int    `sql_type:"INTEGER REFERENCES LevelsSecure (idLvlSec)"`
}
type LevelsSecure struct {
	IdLvlSec   int    `sql_type:"SERIAL PRIMARY KEY" incr:"YES"`
	NameLvlSec string `sql_type:"TEXT"`
}
type Groups struct {
	IdGrp    int    `sql_type:"SERIAL PRIMARY KEY" incr:"YES"`
	NameGrp  string `sql_type:"TEXT"`
	LvlSecId int    `sql_type:"INTEGER REFERENCES LevelsSecure (idLvlSec)"`
}
type Users struct {
	IdUsr     int       `sql_type:"SERIAL PRIMARY KEY" pkey:"YES" incr:"YES"`
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

func (u *Users) CheckUser() error {
	if u.IdUsr == 0 {
		return fmt.Errorf("CheckUser:Field IdUsr is empty")
	}
	idUser, err := u.Find()
	if err != nil {
		return fmt.Errorf("CheckUser:" + err.Error())
	}
	if idUser < 0 {
		// Если пользователя нет в базе - добавляем
		_, err = u.Add()
		if err != nil {
			return fmt.Errorf("CheckUser:" + err.Error())
		}
	} else {
		// Иначе - считываем
		expLst := []Expressions{
			{Key: "idusr",
				Operator: EQ,
				Value:    `'` + fmt.Sprintf("%v", u.IdUsr) + `'`},
		}
		rs, find, _, err := ReadDataRow(&Users{}, expLst, 1)
		if err != nil {
			return fmt.Errorf("CheckUser:" + err.Error())
		}
		if find {
			for _, subRs := range rs {
				mapstructure.Decode(subRs, &u)
			}
		}
	}
	return nil
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

// func (u *Users) GetUserName() (string, error) {
// 	// Проверка на кеш
// 	if u.IdUsr != 0 {
// 		return u.NameUsr, nil
// 	}
// 	// Если не закешировано, то новый поиск
// 	if err := u.CheckUser(); err != nil {
// 		return "", fmt.Errorf("GetUserName:" + err.Error())
// 	}
// 	return u.NameUsr, nil
// }

// func (u *Users) GetChatId() (int64, error) {
// 	// Проверка на кеш
// 	if u.IdUsr != 0 {
// 		return u.ChatIdUsr, nil
// 	}
// 	// Если не закешировано, то новый поиск
// 	if err := u.CheckUser(); err != nil {
// 		return 0, fmt.Errorf("GetChatId:" + err.Error())
// 	}
// 	return u.ChatIdUsr, nil
// }

// // Получение First+Last name
// func (u *Users) GetFLName() (string, error) {
// 	// Проверка на кеш
// 	if u.IdUsr != 0 {
// 		return u.FirstName + " " + u.LastName, nil
// 	}
// 	// Если не закешировано, то новый поиск
// 	if err := u.CheckUser(); err != nil {
// 		return "", fmt.Errorf("GetFLName:" + err.Error())
// 	}
// 	return u.FirstName + " " + u.LastName, nil
// }

// Добавление пользователя в базу
func (u *Users) Add() (int, error) {
	if u.NameUsr == "" {
		u.NameUsr = "Анонимный пользователь"
	}
	if u.IdUsr == 0 || u.NameUsr == "" || u.FirstName == "" ||
		u.LangCode == "" || int(u.ChatIdUsr) == 0 || u.IdLvlSec == 0 {
		return -1, fmt.Errorf("CheckUser:Some field is empty")
	}
	if err := WriteDataStruct(u); err != nil {
		return -1, err
	}
	return u.IdUsr, nil
}

type LimitsDict struct {
	IdLmtDct   int    `sql_type:"SERIAL PRIMARY KEY" pkey:"YES" incr:"YES" miss:"YES"`
	NameLmtDct string `sql_type:"TEXT NOT NULL UNIQUE" fkey:"YES"`
	DescLmtDct string `sql_type:"TEXT NOT NULL"`
	StdValLmt  int    `sql_type:"INTEGER DEFAULT 0"`
}
type Limits struct {
	IdLmt       int       `sql_type:"SERIAL PRIMARY KEY" pkey:"YES" incr:"YES" miss:"YES"`
	ValAvailLmt int       `sql_type:"INTEGER DEFAULT 0"`
	ValUsedLmt  int       `sql_type:"INTEGER DEFAULT 0"`
	ActiveLmt   bool      `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
	TsLmtOn     time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	TsLmtOff    time.Time `sql_type:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	UserId      int       `sql_type:"INTEGER REFERENCES Users (idUsr)" fkey:"YES"`
	LtmDctId    int       `sql_type:"INTEGER REFERENCES LimitsDict (idLmtDct)"`
}

// Функция заполнения лимита по имени лимита и ИД пользователя
func (l *Limits) GetLimit(nameLmt string, usrId int) error {
	// Поиск ИД лимита по имени лимита
	expLst := []Expressions{
		{Key: "NameLmtDct",
			Operator: EQ,
			Value:    `'` + fmt.Sprintf("%v", nameLmt) + `'`},
	}
	rs, find, _, err := ReadDataRow(&LimitsDict{}, expLst, 1)
	if err != nil {
		return fmt.Errorf("GetLimit:" + err.Error())
	}
	lmtDct := LimitsDict{}
	if find {
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &lmtDct)
		}
	} else {
		return fmt.Errorf("GetLimit:Limit type not found")
	}
	// Поиск лимита по ИД лимита и пользователю
	expLst = []Expressions{
		{Key: "UserId",
			Operator: EQ,
			Value:    `'` + fmt.Sprintf("%v", usrId) + `'`},
		{Key: "LtmDctId",
			Operator: EQ,
			Value:    `'` + fmt.Sprintf("%v", lmtDct.IdLmtDct) + `'`},
	}
	rs, find, _, err = ReadDataRow(&Limits{}, expLst, 1)
	if err != nil {
		return fmt.Errorf("GetLimit:" + err.Error())
	}
	if find {
		for _, subRs := range rs {
			mapstructure.Decode(subRs, &l)
		}
	} else {
		return fmt.Errorf("GetLimit:Limit %s for user %s id:%v not found", lmtDct.NameLmtDct, "username", usrId)
	}

	return nil
}
func (l *Limits) SetLimit() error {
	if l.IdLmt == 0 || l.LtmDctId == 0 || l.UserId == 0 ||
		l.ValAvailLmt == 0 {
		return fmt.Errorf("SetTracking:Some field is empty")
	}
	if err := WriteDataStruct(l); err != nil {
		return err
	}
	//Кеширование
	LmtCache[l.IdLmt] = *l
	return nil
}

// Функция инкрементации лимитированного значения и возврата оставщегося лимита
func (l *Limits) IncrLimit(valIncr int) (int, error) {
	if l.IdLmt == 0 {
		return 0, fmt.Errorf("IncrLimit:Limit not initialised")
	}
	if l.ValAvailLmt == l.ValUsedLmt {
		return 0, nil
	}
	l.ValUsedLmt += valIncr

	// Обновляем поле в БД
	data := map[string]string{
		"ValUsedLmt": fmt.Sprintf("%v", l.ValUsedLmt),
	}
	expLst := []Expressions{
		{Key: "IdLmt", Operator: EQ, Value: `'` + fmt.Sprintf("%v", l.IdLmt) + `'`},
	}
	if err := UpdateData("Limits", data, expLst); err != nil {
		return 0, fmt.Errorf("IncrLimit:" + err.Error())
	}
	//Кеширование
	LmtCache[l.IdLmt] = *l
	return l.ValAvailLmt - l.ValUsedLmt, nil
}

type TypeTrackingCrypto struct {
	IdTypTrkCrp       int    `sql_type:"SERIAL PRIMARY KEY" pkey:"YES" incr:"YES"`
	NameTypeTrkCrp    string `sql_type:"TEXT NOT NULL UNIQUE" fkey:"YES"`
	DescTypTrkCrp     string `sql_type:"TEXT NOT NULL"`
	RisingTypTrkCrp   bool   `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
	CalcProcTypTrkCrp bool   `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
}

// Получение всех типов отслеживаний из БД
func (t *TypeTrackingCrypto) GetAllTypeInfo() ([]interface{}, error) {
	expLst := []Expressions{
		{Key: "IdTypTrkCrp", Operator: NotEQ, Value: "0"},
	}
	rs, find, _, err := ReadDataRow(&TypeTrackingCrypto{}, expLst, 0)
	if err != nil {
		return nil, fmt.Errorf("GetAllTypeInfo:" + err.Error())
	}
	if !find {
		return nil, fmt.Errorf("GetAllTypeInfo:not find types")
	}
	res := []interface{}{}
	subFields := TypeTrackingCrypto{}
	for _, subRs := range rs {
		mapstructure.Decode(subRs, &subFields)
		res = append(res, subFields)
	}

	return res, nil
}
func (t *TypeTrackingCrypto) GetTypeInfo() (interface{}, error) {
	expLst := []Expressions{
		{Key: "IdTypTrkCrp", Operator: EQ, Value: fmt.Sprintf("%v", t.IdTypTrkCrp)},
	}
	rs, find, _, err := ReadDataRow(&TypeTrackingCrypto{}, expLst, 1)
	if err != nil {
		return nil, fmt.Errorf("GetTypeInfo:" + err.Error())
	}
	if !find {
		return nil, fmt.Errorf("GetTypeInfo:not find type")
	}
	for _, subRs := range rs {
		subFields := TypeTrackingCrypto{}
		mapstructure.Decode(subRs, &subFields)
		return subFields, nil
	}

	return nil, nil
}

type TrackingCrypto struct {
	IdTrkCrp    int     `sql_type:"SERIAL PRIMARY KEY" pkey:"YES" incr:"YES" miss:"YES"`
	ValTrkCrp   float32 `sql_type:"NUMERIC(19,9)"`
	OnTrkCrp    bool    `sql_type:"BOOLEAN NOT NULL DEFAULT FALSE"`
	LmtId       int     `sql_type:"INTEGER REFERENCES Limits (lmtid)"`
	TypTrkCrpId int     `sql_type:"INTEGER REFERENCES TypeTrackingCrypto (idTypTrkCrp)"`
	DctCrpId    int     `sql_type:"INTEGER REFERENCES DictCrypto (CryptoId)"`
	UserId      int     `sql_type:"INTEGER REFERENCES Users (idUsr)"  fkey:"YES"`
}

func (t *TrackingCrypto) GetTypeInfo() (interface{}, error) {
	// Возможно нужно один раз запустить и держать в кеше
	// Обновлять при обновлении настроек
	expLst := []Expressions{
		{Key: "IdTypTrkCrp", Operator: EQ, Value: fmt.Sprintf("%v", t.TypTrkCrpId)},
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
func (t *TrackingCrypto) OffTracking() error {
	if t.IdTrkCrp == 0 {
		return fmt.Errorf("OffTracking:tracking not initialised")
	}
	t.OnTrkCrp = false
	// Обновляем поле в БД
	data := map[string]string{
		"OnTrkCrp": fmt.Sprintf("%v", t.OnTrkCrp),
	}
	expLst := []Expressions{
		{Key: "IdTrkCrp", Operator: EQ, Value: `'` + fmt.Sprintf("%v", t.IdTrkCrp) + `'`},
	}
	if err := UpdateData("TrackingCrypto", data, expLst); err != nil {
		return fmt.Errorf("OffTracking:" + err.Error())
	}
	TCCache[t.IdTrkCrp] = *t
	return nil
}
func (t *TrackingCrypto) SetTracking() error {
	if t.DctCrpId == 0 || t.TypTrkCrpId == 0 || t.ValTrkCrp == 0 ||
		t.UserId == 0 {
		return fmt.Errorf("SetTracking:Some field is empty")
	}
	if err := WriteDataStruct(t); err != nil {
		return err
	}
	TCCache[t.IdTrkCrp] = *t
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

type CoinMarkets struct {
	IdMrkt          int    `sql_type:"SERIAL PRIMARY KEY" pkey:"YES" incr:"YES"` // ИД
	Url             string `sql_type:"TEXT NOT NULL"`                            // Ссылка
	ShortDesc       string `sql_type:"TEXT NOT NULL" fkey:"YES"`                 // Мнемоника
	Description     string `sql_type:"TEXT NOT NULL"`                            // Описание
	HeadApiKey      string `sql_type:"TEXT NOT NULL"`                            // req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("API_CMC"))
	HeadAccept      string `sql_type:"TEXT NOT NULL"`                            // req.Header.Set("Accepts", "application/json")
	HeadContentType string `sql_type:"TEXT NOT NULL"`                            // req.Header.Add("content-type", "application/json")
	Timeout         int    `sql_type:"INTEGER DEFAULT 0"`                        // Период опроса
}

type CoinMarketsEndpoint struct {
	IdMrktEnd   int    `sql_type:"SERIAL PRIMARY KEY" pkey:"YES" incr:"YES"`           // ИД
	Endpoint    string `sql_type:"TEXT NOT NULL"`                                      // Ресурс ендпоинта
	Method      string `sql_type:"TEXT NOT NULL"`                                      // Метод опроса
	Description string `sql_type:"TEXT NOT NULL"`                                      // Описание
	CoinMrktId  int    `sql_type:"INTEGER REFERENCES CoinMarkets (idMrkt)" fkey:"YES"` // ИД койн маркета
}

type CoinMarketsHand struct {
	IdMrktHand    int    `sql_type:"SERIAL PRIMARY KEY" pkey:"YES" incr:"YES"`                       // ИД
	Key           string `sql_type:"TEXT NOT NULL"`                                                  // Ключ
	Type          string `sql_type:"TEXT NOT NULL"`                                                  // Тип данных
	Value         string `sql_type:"TEXT NOT NULL"`                                                  // Значение
	Description   string `sql_type:"TEXT NOT NULL"`                                                  // Описание
	CoinMrktEndId int    `sql_type:"INTEGER REFERENCES CoinMarketsEndpoint (IdMrktEnd)"  fkey:"YES"` // ИД ендпоинта
}
