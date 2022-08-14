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
	ctx     context.Context
	connUrl string
}

func NewHandler(ctx context.Context, connUrl string) *Handler {
	var err error
	for i := 0; i < 30; i++ {
		_, err = pgx.Connect(ctx, connUrl)
		if err == nil {
			return &Handler{ctx, connUrl}
		}
		time.Sleep(100 * time.Millisecond)
	}
	panic("Could not connect to the DB for 3 seconds")
}

func (h *Handler) getConnection(c *gin.Context) *pgx.Conn {
	conn, err := pgx.Connect(h.ctx, h.connUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error establishing a connections: \n%v", err),
		})
	}
	return conn
}

func (h *Handler) getToken(c *gin.Context) {
	conn := h.getConnection(c)
	if conn == nil {
		return
	}
	defer conn.Close(h.ctx)
	var body Token
	c.BindJSON(&body)
	if body.Password == "" {
		c.JSON(400, gin.H{"error": "Password is a required field"})
		return
	}
	if body.Username == "" {
		c.JSON(400, gin.H{"error": "Username is a required field"})
		return
	}
	var role string
	row := conn.QueryRow(
		h.ctx,
		"Select r.name from users u "+
			"join roles r on u.roleid = r.id "+
			"where username=$1 and password=$2",
		body.Username,
		body.Password)
	row.Scan(&role)
	if role == "" {
		c.JSON(400, gin.H{"error": "invalid username and/or password"})
		return
	}
	c.JSON(200, gin.H{"token": fmt.Sprintf("%s_token", role)})
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
