package coffeezone

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Cafe struct {
	ID       string
	Title    string
	Location []float64
}

func (c *Cafe) String() string {
	return fmt.Sprintf("%s (%s %v)", c.Title, c.ID, c.Location)
}

func NewCafe(ctx context.Context, cafeNode *cdp.Node) *Cafe {
	var titleNodes []*cdp.Node
	err := chromedp.Nodes(cafeNode.FullXPath()+`//a[contains(@class, "title-link")]`, &titleNodes).Do(ctx)
	if err != nil || len(titleNodes) > 1 || titleNodes[0].ChildNodeCount > 1 {
		return nil
	}

	cafeID, _ := cafeNode.Attribute("data-id")
	cafeTitle := strings.TrimSpace(titleNodes[0].Children[0].NodeValue)
	cafeLocation := getLocation(cafeNode)

	return &Cafe{ID: cafeID, Title: cafeTitle, Location: cafeLocation}
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
