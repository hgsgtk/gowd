package gowd_test

import (
	"os"
	"testing"

	"github.com/hgsgtk/gowd"
)

func TestElement_TakeScreenshot(t *testing.T) {
	// Assuming that chromedriver is already running in the local environment
	// > $ chromedriver
	driver := gowd.NewWebDriver()
	browser, err := driver.New()
	if err != nil {
		t.Fatalf("failed to new WebDriver: %#v", err)
	}
	defer func() {
		if err := browser.Close(); err != nil {
			t.Fatalf("failed to close a browser: %#v", err)
		}
	}()

	// Navigate to example.com
	if err := browser.NavigateTo("https://example.com/"); err != nil {
		t.Fatalf("failed to navigate a page: %#v", err)
	}
	div, err := browser.FindElement(gowd.TagName, "h1")
	if err != nil {
		t.Fatalf("failed to find <h1> element: %#v", err)
	}

	screen, err := div.TakeScreenshot()
	if err != nil {
		t.Fatalf("failed to take screenshot: %#v", err)
	}

	if err := os.WriteFile("dist/example_title.png", screen, 0664); err != nil {
		t.Fatalf("failed to write a file: %#v", err)
	}

	// Fixme: assertions on png files generated by this test
}
