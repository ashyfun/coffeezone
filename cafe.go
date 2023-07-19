package coffeezone

import (
	"strconv"

	"github.com/chromedp/cdproto/cdp"
)

func GetLocation(n *cdp.Node) ([]float64, bool) {
	var (
		lonStr string
		latStr string
	)

	lonStr, lonExists := n.Attribute("data-lon")
	latStr, latExists := n.Attribute("data-lat")
	if !lonExists || !latExists {
		return nil, false
	}

	lon, lonErr := strconv.ParseFloat(lonStr, 64)
	lat, latErr := strconv.ParseFloat(latStr, 64)
	if lonErr != nil || latErr != nil {
		return nil, false
	}

	var location []float64
	location = append(location, lon, lat)

	return location, true
}
