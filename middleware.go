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

func checkIsAdmin(c *gin.Context) {
	t := c.GetHeader("Authrization")
	tok, err := parseToken(t)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}
	if tok.Role != "Admin" {
		c.AbortWithStatusJSON(403, gin.H{"error": "user is not an admin"})
		return
	}
	c.Next()
}
