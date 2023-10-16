package twitter

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// Fetches tweets with News from certain teams as well as sports.
func Fetch() {
	fmt.Println("Getting News from Twitter")

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"
	// Set up options to customize the browser
	// opts := append(chromedp.DefaultExecAllocatorOptions[:],
	// 	chromedp.Flag("headless", true), // Run in non-headless mode for visibility (change to true for headless)
	// 	chromedp.UserAgent(userAgent),   // Set the custom user-agent string
	// )

	// // Create an allocator and context with the options
	// allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	// defer cancelAlloc()

	// Create a new context with the allocator options
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Navigate to a website and do something (e.g., print the user-agent)
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://twitter.com"),
		chromedp.WaitVisible(`body > footer`),
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
