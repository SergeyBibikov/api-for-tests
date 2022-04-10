package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ready", readinessCheck)
	router.POST("/token/get", getToken)
	router.POST("/token/check", checkToken)
	router.Run()
}

func readinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func getToken(c *gin.Context) {
	var body map[string]string
	c.BindJSON(&body)
	if body["username"] == "valid_user" && body["password"] == "valid_password" {
		c.JSON(http.StatusCreated, gin.H{
			"success":      true,
			"refreshToken": "ref_resh",
			"accessToken":  "access_t",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":      false,
			"refreshToken": "",
			"accessToken":  "",
		})
	}
}

func checkToken(c *gin.Context) {
	var body map[string]string
	c.BindJSON(&body)
	c.JSON(http.StatusOK, gin.H{
		"refreshToken": "ref_resh",
		"accessToken":  "access_t",
	})
}
