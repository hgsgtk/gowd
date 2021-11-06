package gowd_test

import (
	"fmt"

	"github.com/hgsgtk/gowd"
)

func Example() {
	// Assuming that chromedriver is already running in the local environment
	// > $ chromedriver
	driver := gowd.NewWebDriver()
	// Open browser
	browser, err := driver.New()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := browser.Close(); err != nil {
			panic(err)
		}
	}()

	// Navigate to example.com
	if err := browser.NavigateTo("https://example.com/"); err != nil {
		panic(err)
	}
	url, err := browser.GetCurrentURL()
	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	// Find the title
	titleElement, err := browser.FindElement(gowd.TagName, "h1")
	if err != nil {
		panic(err)
	}
	title, err := titleElement.GetText()
	if err != nil {
		panic(err)
	}
	fmt.Println(title)

	// Find the link to get more information and Click it
	link, err := browser.FindElement(gowd.LinkText, "More information...")
	if err != nil {
		panic(err)
	}
	if err = link.Click(); err != nil {
		panic(err)
	}

	// Confirm the title of the moved page
	titleElement, err = browser.FindElement(gowd.TagName, "h1")
	if err != nil {
		panic(err)
	}
	title, err = titleElement.GetText()
	if err != nil {
		panic(err)
	}
	fmt.Println(title)

	// Find RFC2606 page
	link, err = browser.FindElement(gowd.CSS, "[href=\"/go/rfc2606\"]")
	if err != nil {
		panic(err)
	}
	if err = link.Click(); err != nil {
		panic(err)
	}
	url, err = browser.GetCurrentURL()
	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	// Output:
	// https://example.com/
	// Example Domain
	// IANA-managed Reserved Domains
	// https://www.rfc-editor.org/rfc/rfc2606.html
}
