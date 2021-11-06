# gowd

[![GoDoc](https://godoc.org/github.com/hgsgtk/gowd?status.svg)](https://godoc.org/github.com/hgsgtk/gowd)
[![MIT License](https://img.shields.io/github/license/hgsgtk/gowd)](https://github.com/hgsgtk/gowd/blob/main/LICENSE)

## Description

gowd is a WebDriver binding for Go. See [GoDoc](https://godoc.org/github.com/hgsgtk/gowd) for API specification.

## Getting started

For example, it can be used like this:

```go
package some

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
}
```
