package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ashyfun/coffeezone"
	"github.com/jackc/pgx/v5"
)

const UsageHelp = `
Usage %s: OPTIONS <domain>...

Options:
 --database Writing data to database PostgreSQL (postgresql://url)

Domain examples: zoon.ru/msk spb.zoon.ru
`

func usage() {
	fmt.Fprintf(os.Stderr, strings.TrimLeft(UsageHelp, "\n"), os.Args[0])
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
		links     []string       = flag.Args()
		waitGroup sync.WaitGroup = sync.WaitGroup{}
	)
	for _, v := range links {
		v := v

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			parser := coffeezone.NewParser(v)

			log.Printf("Start parse %s\n", v)
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
		}()
	}

	waitGroup.Wait()
}
