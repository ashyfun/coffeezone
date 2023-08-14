package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ashyfun/coffeezone"
)

const UsageHelp = `
Usage %s: OPTIONS <domain>

Options:
 --database Writing data to database PostgreSQL (postgresql://url)
 --pause    Pause after each scraper in seconds (default 3600)

Domain examples: zoon.ru/msk spb.zoon.ru
`

func start(domain string) {
	parser := coffeezone.NewParser(domain)

	log.Printf("Start parse %s\n", domain)
	parser.Run()

	for _, v := range parser.Cafes {
		if !coffeezone.DatabasePoolAvailable() {
			log.Println(v)
			continue
		}

		v.Handle()
	}

	log.Println("Done")
}

var pause int

func main() {
	coffeezone.SetUsage(UsageHelp)
	flags := coffeezone.ParseFlags(func() {
		flag.IntVar(&pause, "pause", 3600, "")
	})
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	if _, err := coffeezone.SetLogFileOutput(flags.LogFile); err != nil {
		log.Fatalf(`SetLogFileOutput("%s"): %v`, flags.LogFile, err)
	}

	res := coffeezone.SetAndCheckConn(flags.ConnStr)
	if res != nil {
		log.Fatalf(`SetAndCheckConn("%s"): %v`, flags.ConnStr, res)
	}

	coffeezone.NewDatabasePool()
	defer coffeezone.CloseDatabasePool()

	var (
		stopCh                   = make(chan os.Signal, 2)
		domain    string         = flag.Args()[0]
		waitGroup sync.WaitGroup = sync.WaitGroup{}
	)
	waitGroup.Add(1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		defer waitGroup.Done()

		for {
			start(domain)

			select {
			case <-stopCh:
				return
			case <-time.After(time.Duration(pause) * time.Second):
			}
		}
	}()

	waitGroup.Wait()
	close(stopCh)
}
