package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ashyfun/coffeezone"
)

const UsageOf = `
Usage of %s: <domain>...

e.g. zoon.ru/msk spb.zoon.ru
`

func usage() {
	fmt.Fprintf(os.Stderr, UsageOf, os.Args[0])
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	var links []string = flag.Args()
	var waitGroup sync.WaitGroup = sync.WaitGroup{}
	for _, v := range links {
		v := v

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			var url string = fmt.Sprintf("https://%s/restaurants/", v)
			log.Printf("Start parse %s\n", url)
			coffeezone.Run(url)
		}()
	}

	waitGroup.Wait()
}
