package main

import (
	"onysakura.fun/Server/commons"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/rest"
	"strconv"
)

var log = logrus.GetLogger()

func main() {
	port := commons.Configs.Port
	log.Info("server will start at http://127.0.0.1:" + strconv.Itoa(port))
	rest.Run(port)
}

func init() {
	//privateKey, _ := rsa.GenerateKey(rand.Reader, 512)
	//var privateKeyBase64 = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))
	//log.Info("privateKey: ", privateKeyBase64)
}
