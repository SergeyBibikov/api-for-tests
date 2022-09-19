package main

import "github.com/gin-gonic/gin"

func checkToken(c *gin.Context) {
	tok := c.GetHeader("Authrization")
	if tok == "" {
		c.JSON(401, gin.H{"error": "no token provided"})
	}
	c.Next()
}
