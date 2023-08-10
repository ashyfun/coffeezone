package main

import (
	"flag"
	"log"

	"github.com/ashyfun/coffeezone"
	"github.com/gin-gonic/gin"
)

var (
	connStr string
	logFile string
)

func main() {
	flag.StringVar(&connStr, "database", "", "")
	flag.StringVar(&logFile, "logfile", "", "")
	flag.Parse()

	f, err := coffeezone.SetLogFileOutput(logFile)
	if err != nil {
		log.Fatalf(`SetLogFileOutput("%s"): %v`, logFile, err)
	}

	if f != nil {
		gin.DefaultWriter = f
	}

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
