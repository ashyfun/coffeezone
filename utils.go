package coffeezone

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func GetLength(ctx context.Context, sel interface{}) (int, error) {
	var length int
	err := chromedp.Evaluate(
		fmt.Sprintf(`document.querySelectorAll('%v').length`, sel),
		&length,
	).Do(ctx)

	return length, err
}
