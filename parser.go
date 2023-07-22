package coffeezone

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const LimitCafesLength = 300

func Run(url string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.ScrollIntoView("div.catalog-button-showMore", chromedp.NodeVisible),
		chromedp.WaitNotVisible("div.catalog-button-showMore > div.loading-box-img"),
		chromedp.WaitVisible("div.catalog-button-showMore > span.button.button-show-more"),
		chromedp.Sleep(time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				cafesLen, err := GetLength(ctx, "li.minicard-item")
				if err != nil {
					return err
				}

				log.Printf("%d cafes\n", cafesLen)
				if cafesLen >= LimitCafesLength {
					break
				}

				chromedp.Click(
					"div.catalog-button-showMore > span.button.button-show-more",
					chromedp.NodeVisible,
				).Do(ctx)
				chromedp.Sleep(time.Second).Do(ctx)
			}

			var cafeNodes []*cdp.Node
			chromedp.Nodes("li.minicard-item", &cafeNodes).Do(ctx)
			for _, v := range cafeNodes {
				cafe := NewCafe(ctx, v)
				if cafe != nil {
					log.Println(cafe)
				}
			}

			return nil
		}),
	)

	if err != nil {
		log.Fatalf("Failed to open site: %v\n", err)
	}
}
