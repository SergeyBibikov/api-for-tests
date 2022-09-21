package main

import "github.com/gin-gonic/gin"

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
