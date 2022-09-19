package main

import "github.com/gin-gonic/gin"

func checkTokenIsPresent(c *gin.Context) {
	tok := c.GetHeader("Authrization")
	if tok == "" {
		c.JSON(401, gin.H{"error": "no token provided"})
		return
	}
	_, err := parseToken(tok)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error})
		return
	}
	c.Next()
}
