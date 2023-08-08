package main

import (
	"flag"

	"github.com/ashyfun/coffeezone"
	"github.com/gin-gonic/gin"
)

var connStr string

func main() {
	flag.StringVar(&connStr, "database", "", "")
	flag.Parse()

	coffeezone.SetConn(connStr)
	coffeezone.NewDatabasePool()
	defer coffeezone.CloseDatabasePool()

	r := gin.Default()
	r.GET("/cafes", coffeezone.CafesHandler)
	r.Run()
}
