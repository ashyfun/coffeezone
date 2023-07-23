package coffeezone

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const LimitCafesLength = 300

func LoadMoreCafes(ctx context.Context) error {
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

	return nil
}

type Cafe struct {
	ID       string
	Title    string
	Topics   []string
	Location []float64
}

func (c *Cafe) String() string {
	str := fmt.Sprintf("%s (%s %v)", c.Title, c.ID, c.Location)
	for i, v := range c.Topics {
		if i > 0 {
			str += " | "
		} else {
			str += " - "
		}

		str += v
	}

	return str
}

func NewCafe(ctx context.Context, cafeNode *cdp.Node) *Cafe {
	var titleNodes []*cdp.Node
	err := chromedp.Nodes(cafeNode.FullXPath()+`//a[contains(@class, "title-link")]`, &titleNodes).Do(ctx)
	if err != nil || len(titleNodes) != 1 || titleNodes[0].ChildNodeCount != 1 {
		return nil
	}

	newCafe := &Cafe{}

	cafeID, _ := cafeNode.Attribute("data-id")
	newCafe.ID = cafeID
	newCafe.Title = strings.TrimSpace(titleNodes[0].Children[0].NodeValue)
	newCafe.Location = getLocation(cafeNode)

	topicsLen, err := GetLength(ctx, fmt.Sprintf(`li.minicard-item[data-id="%s"] div.minicard-item__features`, cafeID))
	if err == nil && topicsLen > 0 {
		var topicNodes []*cdp.Node
		chromedp.Nodes(
			cafeNode.FullXPath()+`
			//div[contains(@class, "minicard-item__features")]/*[not(contains(@class, "bullet"))]
			`,
			&topicNodes,
		).Do(ctx)
		for _, v := range topicNodes {
			if v.ChildNodeCount == 1 {
				newCafe.Topics = append(
					newCafe.Topics,
					strings.TrimSpace(v.Children[0].NodeValue),
				)
			}
		}
	}

	return newCafe
}

func getLocation(cafe *cdp.Node) []float64 {
	var (
		lonStr string
		latStr string
	)

	lonStr, lonExists := cafe.Attribute("data-lon")
	latStr, latExists := cafe.Attribute("data-lat")
	if !lonExists || !latExists {
		return nil
	}

	lon, lonErr := strconv.ParseFloat(lonStr, 64)
	lat, latErr := strconv.ParseFloat(latStr, 64)
	if lonErr != nil || latErr != nil {
		return nil
	}

	var location []float64
	location = append(location, lon, lat)

	return location
}
