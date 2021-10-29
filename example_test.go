package gowd_test

import (
	"fmt"
	"github.com/hgsgtk/gowd"
)

func Example() {
	driver := gowd.WebDriver{}
	browser := driver.New()

	fmt.Println(browser.SessionID)
	// Output: faked-session-id
}
