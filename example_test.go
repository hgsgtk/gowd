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

	fmt.Println(browser.SessionID)
}
