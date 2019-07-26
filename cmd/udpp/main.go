package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	gh "github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/sunliang711/udpp/handlers"
	"github.com/sunliang711/udpp/models"
	"github.com/sunliang711/udpp/utils"
)

func main() {
	pflag.StringP("mongourl", "", "", "mongodb url")
	pflag.StringP("blockDB", "", "", "mongodb url")
	pflag.IntP("port", "p", 0, "listen port")
	pflag.Bool("enableCors", false, "enable cors")
	pflag.StringP("logfile", "l", "", "logfile path")
	pflag.Bool("auth", true, "enable auth")
	pflag.String("loglevel", "debug", "log level")

	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	if !viper.GetBool("auth") {
		logrus.Infof("Disable auth")
	}

	models.InitMongo(viper.GetString("mongourl"))
	models.InitBlockDb(viper.GetString("blockDB"))
	logrus.SetLevel(utils.LogLevel(viper.GetString("loglevel")))

	rt := handlers.Router(viper.GetBool("enableCors"))
	addr := fmt.Sprintf(":%d", viper.GetInt("port"))
	logrus.Infof("Listen on %v", addr)

	var output io.Writer
	logfile := viper.GetString("logfile")
	output = os.Stdout
	if len(logfile) > 0 {
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		output = io.MultiWriter(os.Stdout, f)
		logrus.SetOutput(output)
	}

	//TODO https
	http.ListenAndServe(addr, gh.LoggingHandler(output, rt))
}
