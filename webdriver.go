package gowd

// WebDriver manages the settings for browser operations.
// See also https://www.selenium.dev/documentation/webdriver/.
// Fixme: for multi-browser support, it is better to use interface.
type WebDriver struct {
	// TODO: Add capabilities field
}

// New opens a new browser.
// Todo: Support multiple options such as "remote", "chromedriver", "geckodriver"...
func (wd *WebDriver) New() *Browser {
	return &Browser{
		// Fixme: fake it!
		SessionID: SessionID("faked-session-id"),
	}
}

// SessionID used for communication with WebDriver.
// https://www.w3.org/TR/webdriver/#sessions
type SessionID string

// Browser represents the state of a browser opened by WebDriver.
type Browser struct {
	SessionID SessionID
}
