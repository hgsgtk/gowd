package gowd_test

import (
	"fmt"
	"github.com/hgsgtk/gowd"
)

func Example() {
	// Assuming that chromedriver is already running in the local environment
	// > $ chromedriver
	driver := gowd.NewWebDriver()
	browser, err := driver.New()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := browser.Close(); err != nil {
			panic(err)
		}
	}()

	if err := browser.NavigateTo("https://example.com/"); err != nil {
		panic(err)
	}
	url, err := browser.GetCurrentURL()
	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	titleElement, err := browser.FindElement(gowd.TagName, "h1")
	if err != nil {
		panic(err)
	}

	title, err := titleElement.GetText()
	if err != nil {
		panic(err)
	}
	fmt.Println(title)
	// Output:
	// https://example.com/
	// Example Domain
}
