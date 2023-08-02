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
	"github.com/jackc/pgx/v5"
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

type LocationType struct {
	Address   string
	Longitude float64
	Latitude  float64
}

func (c *LocationType) CreateOrUpdate() (string, []any) {
	values := []any{c.Address, c.Longitude, c.Latitude}
	return `
	insert into cz_cafe_locations (address_name, longitude, latitude)
	values ($1, $2, $3)
	on conflict (longitude, latitude) do update set address_name = $1, longitude = $2, latitude = $3
	returning id
	`, values
}

type Cafe struct {
	ID       string
	Title    string
	Topics   []string
	Location *LocationType
}

func (c *Cafe) CreateOrUpdate() (string, []any) {
	var locationID int32
	if c.Location != nil {
		sql, args := c.Location.CreateOrUpdate()
		QueryRowExec(func(r pgx.Row) {
			if err := r.Scan(&locationID); err != nil {
				log.Printf("Failed to create or update location: %v", err)
				locationID = 0
				return
			}

			log.Printf("%s: Location ID: %d", c.ID, locationID)
		}, sql, args...)
	}

	values := []any{c.ID, c.Title, locationID}
	return `
	insert into cz_cafes (code, title, location_id)
	values ($1, $2, $3) on conflict (code) do update set title = $2, location_id = $3, updated_at = now()
	returning code
	`, values
}

func (c *Cafe) String() string {
	return fmt.Sprintf("%s %s", c.ID, c.Title)
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
	newCafe.Location = getLocation(ctx, cafeNode)

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

func getLocation(ctx context.Context, cafe *cdp.Node) *LocationType {
	var (
		addrStr string
		lonStr  string
		latStr  string
	)

	var addrNodes []*cdp.Node
	err := chromedp.Nodes(cafe.FullXPath()+`//address/span[contains(@class, "address")]`, &addrNodes).Do(ctx)
	if err != nil || len(addrNodes) != 1 || addrNodes[0].ChildNodeCount != 1 {
		return nil
	}

	addrStr = strings.TrimSpace(addrNodes[0].Children[0].NodeValue)

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

	return &LocationType{addrStr, lon, lat}
}
