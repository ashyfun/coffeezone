package coffeezone

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

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

func SetUsage(help string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, strings.TrimLeft(help, "\n"), os.Args[0])
	}
}

type Options struct {
	ConnStr string
	LogFile string
}

func ParseFlags(before func()) *Options {
	opts := &Options{}

	flag.StringVar(&opts.ConnStr, "database", "", "")
	flag.StringVar(&opts.LogFile, "logfile", "", "")

	if before != nil {
		before()
	}

	flag.Parse()
	return opts
}
