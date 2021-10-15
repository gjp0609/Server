package user

import (
	md5 "crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"onysakura.fun/Server/commons"
	"onysakura.fun/Server/commons/data"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/commons/sqlite"
	"strings"
	"time"
)

var log = logrus.GetLogger()

type User struct {
	Id       int
	Username string
	Password string
}

type MyCustomClaims struct {
	Foo string
	jwt.RegisteredClaims
}

func Login(writer http.ResponseWriter, request *http.Request) {
	var returnMsg = data.NewErrorMsg()
	var err error
	defer func() {
		if err != nil {
			log.Warning(err)
			returnMsg.Msg = returnMsg.Msg + ", Reason: " + err.Error()
		}
		_, _ = writer.Write(returnMsg.ToString())
	}()
	var loginUser User
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&loginUser)
	if err != nil {
		returnMsg.Msg = "decode param error"
		return
	}
	log.Info("login user: ", loginUser)

	user := GetUser(loginUser.Username)
	log.Info("User: ", user)
	if user == nil {
		returnMsg.Msg = "user not found"
		return
	}
	if !strings.EqualFold(getMD5Hash(user.Password), loginUser.Password) {
		returnMsg.Msg = "password invalid"
		return
	}
	var privateKey *rsa.PrivateKey
	privateKey, err = loadPrivateKeyBase64(commons.Configs.PrivateKey)
	if err != nil {
		returnMsg.Msg = "load private key error"
		return
	}
	//claims := MyCustomClaims{
	//	Foo: "sadasd",
	//}
	//token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	now := time.Now()
	expire := now.Add(time.Minute * 30)
	token := jwt.New(jwt.SigningMethodRS256)
	// 获取 claims
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = now.Unix()
	claims["exp"] = expire.Unix()
	claims["username"] = user.Username
	var tokenStr string
	tokenStr, err = token.SignedString(privateKey)
	if err != nil {
		returnMsg.Msg = "issue token error"
		return
	}
	log.Info("Token: ", tokenStr)
	returnMsg = data.Msg{Code: data.MsgOk, Msg: "ok", Data: tokenStr}
}

func Auth(authorization string) (*string, error) {
	privateKey, err := loadPrivateKeyBase64(commons.Configs.PrivateKey)
	if err != nil {
		return nil, err
	}
	parse, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		return privateKey.Public(), nil
	})
	if err != nil {
		return nil, err
	}
	//log.Debug("parse: ", parse)
	if !parse.Valid {
		return nil, errors.New("errrr")
	}
	claims := parse.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	return &username, nil
}

func GetUser(username string) *User {
	var err error
	defer func() {
		if err != nil {
			log.Warn("get user fail: ", err)
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
			return nil
		}
		return &user
	}
	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func loadPrivateKeyBase64(base64key string) (*rsa.PrivateKey, error) {
	decodeString, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(decodeString)
	if err != nil {
		return nil, errors.New("parse private key error")
	}

	return privateKey, nil
}

func getParam(param string, request *http.Request, defaultValue string) *string {
	var value string
	if request.URL.Query()[param] != nil {
		value = request.URL.Query()[param][0]
		return &value
	}
	if &defaultValue != nil {
		return &defaultValue
	}
	return nil
}
