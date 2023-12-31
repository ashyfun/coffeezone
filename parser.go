package coffeezone

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const Concurrency = 6

type Parser struct {
	url    string
	ctx    context.Context
	cancel context.CancelFunc
	Cafes  []*Cafe
}

func NewParser(url string) *Parser {
	ctx, cancel := chromedp.NewContext(context.Background())
	return &Parser{
		url:    fmt.Sprintf("https://%s/restaurants/", url),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Parser) Run() {
	defer p.cancel()

	err := chromedp.Run(p.ctx,
		chromedp.Navigate(p.url),
		chromedp.WaitReady("body"),
		chromedp.ScrollIntoView("div.catalog-button-showMore", chromedp.NodeVisible),
		chromedp.WaitNotVisible("div.catalog-button-showMore > div.loading-box-img"),
		chromedp.WaitVisible("div.catalog-button-showMore > span.button.button-show-more"),
		chromedp.Sleep(time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := LoadMoreCafes(ctx)
			if err != nil {
				return err
			}

			var (
				waitGroup sync.WaitGroup
				cafeNodes []*cdp.Node
			)
			chromedp.Nodes("li.minicard-item", &cafeNodes).Do(ctx)

			cafeCh := make(chan *cdp.Node)
			for i := 0; i < Concurrency; i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()

					for v := range cafeCh {
						cafe := NewCafe(ctx, v)
						if cafe != nil {
							p.Cafes = append(p.Cafes, cafe)
						}
					}
				}()
			}

			for _, v := range cafeNodes {
				cafeCh <- v
			}
			close(cafeCh)

			waitGroup.Wait()

			return nil
		}),
	)

	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
}
