package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func (c *Context) Response(res any, pagination *Pagination) {
	obj := gin.H{
		"success": true,
		"data":    res,
	}
	if pagination != nil && pagination.Total > 0 {
		obj["pagination"] = pagination
	}
	c.JSON(http.StatusOK, obj)
}

func (c *Context) BadRequest(err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}

func (c *Context) InternalError(err error) {
	log.Printf("Internal Server Error: \"%s\"\n", err)
	c.AbortWithStatus(http.StatusInternalServerError)
}
