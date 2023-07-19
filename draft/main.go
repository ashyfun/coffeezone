package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/ashyfun/coffeezone"
)

func main() {
	var links = []string{
		"zoon.ru/msk",
		"spb.zoon.ru",
		"vladimir.zoon.ru",
	}
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
