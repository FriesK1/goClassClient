package main

import (
	"github.com/sirupsen/logrus"

	config "github.com/FriesK1/goMartialConfig"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	config.SetAppName("myClientApp", "mca")
	config.SetAppVersion("1.0.0")

	viper.SetDefault("server.address", "0.0.0.0")
	pflag.StringP("bind.address", "a", viper.GetString("bind.address"), "Address to bind HTTP Server to")

	viper.SetDefault("server.port", 3000)
	pflag.IntP("server.port", "p", viper.GetInt("server.port"), "Port to bind HTTP Server to")

	viper.SetDefault("host.address", "127.0.0.1")
	pflag.StringP("host.address", "A", viper.GetString("host.address"), "HTTP Address to connect to")

	viper.SetDefault("host.port", 3001)
	pflag.IntP("host.port", "P", viper.GetInt("host.port"), "HTTP Port to connect to")

	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	config.GetConfigs()

	if err := StartWebServer(); err != nil {
		logrus.Fatal(err)
	}
}
