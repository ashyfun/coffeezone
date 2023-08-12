package main

import (
	"log"

	"github.com/ashyfun/coffeezone"
	"github.com/gin-gonic/gin"
)

func main() {
	flags := coffeezone.ParseFlags(nil)

	f, err := coffeezone.SetLogFileOutput(flags.LogFile)
	if err != nil {
		log.Fatalf(`SetLogFileOutput("%s"): %v`, flags.LogFile, err)
	}

	if f != nil {
		gin.DefaultWriter = f
	}

	coffeezone.SetConn(flags.ConnStr)
	coffeezone.NewDatabasePool()
	defer coffeezone.CloseDatabasePool()

	router := NewRouter()
	api1 := router.WrapGroup("/api/v1")
	{
		api1.View("GET", "/cafes", CafesHandler)
	}

	router.Run()
}
