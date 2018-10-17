package main

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"logger/controller"
	"fmt"
	"github.com/revel/config"
	"logger/library"
)

var HttpPort int
var HttpHost string
var Config *config.Context

func init () {
	var err error
	Config, err = config.LoadContext("app.ini", []string{"config"})
	if err != nil {
		fmt.Println("配置文件读取失败");
		return ;
	}
	Config.SetSection("server")
	HttpPort = Config.IntDefault("http.port", 8080)
	HttpHost = Config.StringDefault("http.host", "127.0.0.1")
	Config.SetSection("DEFAULT")
	library.InitMysql(Config)
}

func main () {
	handle :=  mux.NewRouter()
	handle.HandleFunc("/", controller.Log)
	handle.HandleFunc("log", controller.Log)
	fmt.Print(HttpHost + strconv.Itoa(HttpPort))
	http.ListenAndServe(HttpHost + ":" + strconv.Itoa(HttpPort), handle)
}
