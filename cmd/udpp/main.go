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
)

func main() {
	pflag.StringP("mongourl", "", "", "mongodb url")
	pflag.IntP("port", "p", 0, "listen port")
	pflag.Bool("enableCors", false, "enable cors")
	pflag.StringP("logfile", "l", "", "logfile path")
	pflag.Bool("auth", true, "enable auth")
	pflag.String("loglevel", "debug", "log level")

	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	//models.InitMysql(viper.GetString("dsn"))
	//defer models.CloseMysql()
	if !viper.GetBool("auth") {
		logrus.Infof("Disable auth")
	}

	models.InitMongo(viper.GetString("mongourl"))
	loglevel := viper.GetString("loglevel")
	if len(loglevel) == 0 {
		loglevel = "info"
	}
	switch loglevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

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
