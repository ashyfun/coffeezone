package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ashyfun/coffeezone"
	"github.com/jackc/pgx/v5"
)

const UsageHelp = `
Usage %s: OPTIONS <domain>

Options:
 --database Writing data to database PostgreSQL (postgresql://url)

Domain examples: zoon.ru/msk spb.zoon.ru
`

func usage() {
	fmt.Fprintf(os.Stderr, strings.TrimLeft(UsageHelp, "\n"), os.Args[0])
}

func start(domain string) {
	parser := coffeezone.NewParser(domain)

	log.Printf("Start parse %s\n", domain)
	parser.Run()

	for _, v := range parser.Cafes {
		if !coffeezone.DatabasePoolAvailable() {
			log.Println(v)
			continue
		}

		sql, args := v.CreateOrUpdate()
		coffeezone.QueryRowExec(func(r pgx.Row) {
			var code string
			if err := r.Scan(&code); err != nil {
				log.Printf("Failed to add/update entry %s: %v", v.ID, err)
				return
			}

			log.Printf("Entry %s added/updated", code)
		}, sql, args...)
	}
}

var connStr string

func main() {
	flag.Usage = usage
	flag.StringVar(&connStr, "database", "", "")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	coffeezone.SetConn(connStr)
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
			case <-time.After(20 * time.Second):
			}
		}
	}()

	waitGroup.Wait()
}
