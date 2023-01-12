package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctx context.Context
}

func NewHandler(ctx context.Context, connUrl string) *Handler {
	for i := 0; i < 30; i++ {
		_, err := getDbConnection()
		if err == nil {
			return &Handler{ctx}
		}
		time.Sleep(100 * time.Millisecond)
	}
	panic("Could not connect to the DB for 3 seconds")
}

func (h *Handler) getToken(c *gin.Context) {
	conn, _ := getDbConnection()
	if conn == nil {
		return
	}
	defer conn.Close(h.ctx)

	var body GetTokenBody
	c.BindJSON(&body)
	username := body.Username
	password := body.Password

	if password == "" {
		c.JSON(400, gin.H{"error": "Password is a required field"})
		return
	}
	if username == "" {
		c.JSON(400, gin.H{"error": "Username is a required field"})
		return
	}

	var role string
	row := conn.QueryRow(
		h.ctx,
		"Select r.name from users u "+
			"join roles r on u.roleid = r.id "+
			"where username=$1 and password=$2",
		username,
		password)
	row.Scan(&role)
	if role == "" {
		c.JSON(400, gin.H{"error": "invalid username and/or password"})
		return
	}

	c.JSON(200, gin.H{"token": fmt.Sprintf("%s_token_%s", role, username)})
}

func (h *Handler) validateToken(c *gin.Context) {

	conn, _ := getDbConnection()
	if conn == nil {
		return
	}
	defer conn.Close(h.ctx)

	var body ValidateTokenBody
	c.BindJSON(&body)
	t := body.Token
	token, err := parseToken(t)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	var role string
	row := conn.QueryRow(
		h.ctx,
		"Select r.name from users u "+
			"join roles r on u.roleid = r.id "+
			"where username=$1",
		token.Username)
	row.Scan(&role)

	if role == "" {
		c.JSON(401, gin.H{"error": "invalid username"})
		return
	}
	if role != token.Role {
		c.JSON(401, gin.H{"error": "incorrect user role"})
		return
	}
	c.JSON(200, gin.H{})
}

func (h *Handler) register(c *gin.Context) {
	conn, _ := getDbConnection()
	if conn == nil {
		return
	}
	defer conn.Close(h.ctx)

	var body RegisterBody
	c.BindJSON(&body)

	username := body.Username
	password := body.Password
	email := body.Email

	if username == "" || password == "" || email == "" {
		c.JSON(400, gin.H{"error": "Username, password and email are required"})
		return
	}
	if err := conn.QueryRow(h.ctx, "select id from users where username = $1", username).Scan(nil); err == nil {
		c.JSON(400, gin.H{"error": "The username is already taken"})
		return
	}
	if err := conn.QueryRow(h.ctx, "select id from users where email = $1", email).Scan(nil); err == nil {
		c.JSON(400, gin.H{"error": "Other user has already used this email."})
		return
	}
	if err := validatePassword(password); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := validateEmail(email); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	row := conn.QueryRow(h.ctx,
		"insert into users values(DEFAULT, $1, $2, $3, 2, null, null) returning id",
		username, email, password)

	var id int
	err := row.Scan(&id)

	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("User insertion error: %s", err.Error())})
		return
	}

	c.JSON(201, gin.H{"message": "user created", "userId": id})
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func (h *Handler) getTeams(c *gin.Context) {
	err := validateTeamsQueryParams(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query := getTeamsSQLQuery(c)
	conn, _ := getDbConnection()
	defer conn.Close(c)
	rows, _ := conn.Query(c, query)
	defer rows.Close()

	type Team struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Conf string `json:"conference"`
		Div  string `json:"division"`
		Year int    `json:"est_year"`
	}
	results := []Team{}

	for rows.Next() {
		t := Team{}
		rows.Scan(&t.Id, &t.Name, &t.Conf, &t.Div, &t.Year)
		results = append(results, t)
	}

	c.JSON(200, results)
}

// WIP
// TODO: validate token before processing a request
func (h *Handler) deleteTeam(c *gin.Context) {
	conn, _ := getDbConnection()
	_, err := conn.Exec(c, "delete from teams where id = $1", c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}
