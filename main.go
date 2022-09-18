package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	dbUrl := fmt.Sprintf("postgres://postgres:%s@localhost:5432/postgres", os.Getenv("DBPass"))
	ctx := context.Background()
	handler := NewHandler(ctx, dbUrl)

	router := gin.Default()
	router.GET("/ready", handler.health)

	router.POST("/register", handler.register)

	router.POST("/token/get", handler.getToken)
	router.POST("/token/validate", handler.validateToken)

	router.GET("/teams", handler.getTeams)
	router.DELETE("/teams/:id", handler.deleteTeam)
	router.Run()
}
