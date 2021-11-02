package gowd_test

import "github.com/hgsgtk/gowd"

func ExampleBrowser_TakeScreenshot() {
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

	// Navigate to example.com
	if err := browser.NavigateTo("https://example.com/"); err != nil {
		panic(err)
	}
	_, err = browser.TakeScreenshot()
	if err != nil {
		panic(err)
	}

	// Output:
}
