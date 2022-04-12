package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ready", readinessCheck)
	router.POST("/token/get", getToken)
	router.POST("/token/check", checkToken)
	router.GET("/users/:id", getUser)
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
	if body["accessToken"] == "access_t" {
		c.JSON(http.StatusOK, gin.H{
			"valid": true,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
		})
	}
}

func getUser(c *gin.Context) {
	params := c.Request.URL.Query()
	user, isPresent := params["id"]
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"user":    "",
		})
	}
	var body map[string]string
	c.BindJSON(&body)
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"user":    fmt.Sprintf("user%s", user),
	})
}
