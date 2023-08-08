package main

import (
	"flag"

	"github.com/ashyfun/coffeezone"
)

var connStr string

func main() {
	flag.StringVar(&connStr, "database", "", "")
	flag.Parse()

	coffeezone.SetConn(connStr)
	coffeezone.NewDatabasePool()
	defer coffeezone.CloseDatabasePool()

	router := NewRouter()
	api1 := router.WrapGroup("/api/v1")
	{
		api1.View("GET", "/cafes", CafesHandler)
	}

	router.Run()
}
