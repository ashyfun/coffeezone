package main

import (
	"log"
	"net/http"

	"github.com/ashyfun/coffeezone"
	"github.com/gin-gonic/gin"
)

const ACAllow = "Access-Control-Allow"

func cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set(ACAllow+"-Origin", "*")
		ctx.Writer.Header().Set(ACAllow+"-Credentials", "true")
		ctx.Writer.Header().Set(ACAllow+"-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set(ACAllow+"-Methods", "POST, OPTIONS, GET, PUT")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func main() {
	flags := coffeezone.ParseFlags(nil)

	f, err := coffeezone.SetLogFileOutput(flags.LogFile)
	if err != nil {
		log.Fatalf(`SetLogFileOutput("%s"): %v`, flags.LogFile, err)
	}

	if f != nil {
		gin.DefaultWriter = f
	}

	res := coffeezone.SetAndCheckConn(flags.ConnStr)
	if res != nil {
		log.Fatalf(`SetAndCheckConn("%s"): %v`, flags.ConnStr, res)
	}

	coffeezone.NewDatabasePool()
	defer coffeezone.CloseDatabasePool()

	router := NewRouter()
	if gin.Mode() == "debug" {
		router.Use(cors())
	}
	api1 := router.WrapGroup("/api/v1")
	{
		api1.View("GET", "/cafes", CafesHandler)
	}

	router.Run()
}
