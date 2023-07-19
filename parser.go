package coffeezone

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const LimitNodes = 300

func getLength(ctx context.Context, sel interface{}) (int, error) {
	var length int
	err := chromedp.Evaluate(
		fmt.Sprintf(`document.querySelectorAll('%v').length`, sel),
		&length,
	).Do(ctx)

	return length, err
}

func Run(url string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.ScrollIntoView("div.catalog-button-showMore", chromedp.NodeVisible),
		chromedp.WaitNotVisible("div.catalog-button-showMore > div.loading-box-img"),
		chromedp.WaitVisible("div.catalog-button-showMore > span.button.button-show-more"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			length, err := getLength(ctx, "li.minicard-item")
			if err != nil {
				log.Println(err)
				return err
			}

			for {
				log.Printf("Nodes length: %d\n", length)
				chromedp.Click(
					"div.catalog-button-showMore > span.button.button-show-more",
					chromedp.NodeVisible,
				).Do(ctx)

				for {
					len, err := getLength(ctx, "li.minicard-item")
					if err != nil {
						log.Println(err)
						return err
					}

					if len > length {
						length = len
						if len == LimitNodes {
							return nil
						}

						break
					}

					time.Sleep(time.Second)
				}
			}
		}),
		chromedp.Nodes("li.minicard-item", &nodes),
	); err != nil {
		log.Fatalf("Failed to open site: %v\n", err)
	}

	log.Printf("Total elements: %d\n", len(nodes))
}
