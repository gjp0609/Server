package user

import (
	_ "github.com/mattn/go-sqlite3"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/commons/sqlite"
)

var log = logrus.GetLogger()

type User struct {
	Id       int
	Username string
	Password string
}

func init() {

	user := GetUser("admin")
	log.Println("user", user)
	if user != nil {
		return
	}
}

func GetUser(username string) *User {
	var err error
	defer func() {
		if err != nil {
			log.Warning(err)
		}
	}()
	stmt, err := sqlite.DB.Prepare("select * from user where username = ?")
	if err != nil {
		return nil
	}
	rows, err := stmt.Query(username)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			log.Warn("get user info fail. ", err)
			return nil
		}
		return &user
	}
	return nil
}
