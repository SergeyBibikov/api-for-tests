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
	defer handler.dbConnection.Close(ctx)

	router := gin.Default()
	router.GET("/ready", handler.health)
	router.POST("/token/get", handler.getToken)
	router.POST("/token/check", handler.checkToken)
	router.GET("/users/:id", handler.getUser)
	router.Run()
}
