package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pgx "github.com/jackc/pgx/v4"
)

type Handler struct {
	ctx          context.Context
	dbConnection *pgx.Conn
}

func NewHandler(ctx context.Context, connUrl string) *Handler {
	var conn *pgx.Conn
	var err error
	for i := 0; i < 30; i++ {
		conn, err = pgx.Connect(ctx, connUrl)
		if err == nil {
			return &Handler{ctx: ctx, dbConnection: conn}
		}
		time.Sleep(100 * time.Millisecond)
	}
	panic("Could not connect to the DB for 3 seconds")
}

func (h *Handler) getToken(c *gin.Context) {
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
func (h *Handler) checkToken(c *gin.Context) {
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
func (h *Handler) getUser(c *gin.Context) {
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
func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
