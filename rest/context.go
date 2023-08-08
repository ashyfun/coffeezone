package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func (c *Context) Response(res interface{}, err error) {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

func (c *Context) InternalError(err error) {
	log.Printf("Internal Server Error: \"%s\"\n", err)
	c.AbortWithStatus(http.StatusInternalServerError)
}
