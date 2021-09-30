package sqlite

import (
	"database/sql"
	"onysakura.fun/Server/commons"
	"onysakura.fun/Server/commons/logrus"
)

var log = logrus.GetLogger()
var DB *sql.DB

func init() {
	database, err := sql.Open("sqlite3", commons.Configs.SQLitePath)
	if err != nil {
		log.Panic("open database fail ", err)
	}
	DB = database
}

type Dao interface {
	QueryFirst()
	QueryAll()
	Insert()
	Update()
}
