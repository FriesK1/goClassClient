package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"net/http"

	DataObject "github.com/FriesK1/goClassData"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func StartWebServer() (err error) {
	var bindUrl string

	bindUrl = strings.Join([]string{viper.GetString("server.host"), viper.GetString("server.port")}, ":")
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	if err = router.SetTrustedProxies([]string{"10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12", "127.0.0.0/8"}); err != nil {
		logrus.Fatal(err)
	}

	router.LoadHTMLGlob("templates/*.tmpl")

	router.GET("/users", GetUserList)
	router.GET("/user/:name", GetUser)
	router.GET("/health-check", HealthCheck)

	_ = router.Run(bindUrl)
	return
}

func GetUserList(c *gin.Context) {
	var data DataObject.Employees

	res, err := http.Get(
		fmt.Sprintf("http://%s:%s/users", viper.GetString("host.address"), viper.GetString("host.port")),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	if res.StatusCode > 299 {
		logrus.Errorln("Response failed with code: ", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.Fatal(err)
	}

	c.HTML(http.StatusOK, "show_user_list.tmpl", data)
}

func GetUser(c *gin.Context) {
	var data DataObject.Employee
	var name = c.Param("name")

	res, err := http.Get(
		fmt.Sprintf("http://%s:%s/user/%s", viper.GetString("host.address"), viper.GetString("host.port"), name),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	if res.StatusCode > 299 {
		logrus.Errorf("Response failed with code: %d", res.StatusCode)
		c.HTML(http.StatusBadRequest, "user_not_found.tmpl", name)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.Fatal(err)
	}

	c.HTML(http.StatusOK, "show_user.tmpl", data)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
