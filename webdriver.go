package gowd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WebDriver manages the settings for browser operations.
// See also https://www.selenium.dev/documentation/webdriver/.
// Fixme: for multi-browser support, it is better to use interface.
type WebDriver struct {
	// RemoteEndURL is a URL of "Remote end" node.
	// https://www.w3.org/TR/webdriver/#nodes
	// > The remote end hosts the server side of the protocol.
	RemoteEndURL *url.URL
	// Client communicates with WebDriver Remote end.
	// https://pkg.go.dev/net/http#Client
	// > The Client's Transport typically has internal state (cached TCP connections),
	// > so Clients should be reused instead of created as needed.
	// > Clients are safe for concurrent use by multiple goroutines.
	Client *http.Client
	// TODO: Add capabilities field
}

// responseNewSession represents the structure of response "New Session" (POST /session)
// https://www.w3.org/TR/webdriver/#dfn-new-sessions
type responseNewSession struct {
	Value struct {
		SessionID    SessionID   `json:"sessionId"`
		Capabilities interface{} `json:"capabilities"` // very flexible
	} `json:"value"`
	Raw json.RawMessage `json:"-"` // for marshaling raw response
}

// UnmarshalJSON is implemented to map raw data into struct.
// https://budougumi0617.github.io/2021/10/25/smart_saving_json_raw_message/
func (r *responseNewSession) UnmarshalJSON(data []byte) error {
	type response responseNewSession
	var rs response
	if err := json.Unmarshal(data, &rs); err != nil {
		return err
	}

	*r = (responseNewSession)(rs)
	r.Raw = data
	return nil
}

func NewWebDriver() *WebDriver {
	c := http.DefaultClient
	c.Timeout = 5 * time.Second // choose 3s at random
	// Todo: set a proxy between local end and remote end
	// https://www.w3.org/TR/webdriver/#nodes
	ru, err := url.Parse("http://localhost:9515")
	if err != nil {
		return nil
	}
	return &WebDriver{
		RemoteEndURL: ru,
		Client:       new(http.Client),
	}
}

// New opens a new browser.
// Todo: Support multiple options such as "remote", "chromedriver", "geckodriver"...
func (wd *WebDriver) New() (*Browser, error) {
	// TODO: support other capabilities options
	// TODO: enable users to set chromeOptions https://github.com/hgsgtk/gowd/pull/8#issuecomment-955629879
	rb := `
{
	"capabilities": {
		"alwaysMatch": {
			"goog:chromeOptions": {
				"args": ["--no-sandbox", "--headless"]
			}
		}
	}
}
`

	// https://www.w3.org/TR/webdriver/#new-session
	u := wd.RemoteEndURL.String() + "/session"
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(rb))
	if err != nil {
		return nil, fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := wd.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("got http response error: %w", err)
	}
	defer resp.Body.Close()

	var rns responseNewSession
	if err := json.NewDecoder(resp.Body).Decode(&rns); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	if resp.StatusCode != 200 {
		// Fixme: define the struct of error response and handle it.
		// https://www.w3.org/TR/webdriver/#errors
		return nil, fmt.Errorf("got invalid http status code: %d, body: %s", resp.StatusCode, rns.Raw)
	}

	if rns.Value.SessionID == "" {
		return nil, fmt.Errorf("got empty sessionId, response body: %s", rns.Raw)
	}

	return &Browser{
		SessionID: rns.Value.SessionID,
		driver:    wd,
	}, nil
}

// SessionID used for communication with WebDriver.
// https://www.w3.org/TR/webdriver/#sessions
// ex. 15f0a07b906057033a40c9946005c86d
type SessionID string

// Browser represents the state of a browser opened by WebDriver.
type Browser struct {
	SessionID SessionID
	driver    *WebDriver
}

func (b *Browser) Close() error {
	// https://www.w3.org/TR/webdriver/#delete-session
	u := b.driver.RemoteEndURL.String() + "/session/" + string(b.SessionID)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := b.driver.Client.Do(req)
	if err != nil {
		return fmt.Errorf("got http response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Fixme: define the struct of error response and handle it.
		// https://www.w3.org/TR/webdriver/#errors
		return fmt.Errorf("got invalid http status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *Browser) NavigateTo(url string) error {
	type rp struct {
		URL string `json:"url"`
	}
	rpb := rp{
		URL: url,
	}
	body, err := json.Marshal(rpb)
	if err != nil {
		return fmt.Errorf("can't marshal json body: %w", err)
	}

	// https://www.w3.org/TR/webdriver/#navigate-to
	u := b.driver.RemoteEndURL.String() + "/session/" + string(b.SessionID) + "/url"
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := b.driver.Client.Do(req)
	if err != nil {
		return fmt.Errorf("got http response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Fixme: define the struct of error response and handle it.
		// https://www.w3.org/TR/webdriver/#errors
		return fmt.Errorf("got invalid http status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *Browser) GetCurrentURL() (string, error) {
	// https://www.w3.org/TR/webdriver/#get-current-url
	u := b.driver.RemoteEndURL.String() + "/session/" + string(b.SessionID) + "/url"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := b.driver.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("got http response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Fixme: define the struct of error response and handle it.
		// https://www.w3.org/TR/webdriver/#errors
		return "", fmt.Errorf("got invalid http status code: %d", resp.StatusCode)
	}

	type rf struct {
		Value string `json:"value"`
	}
	var rb rf
	if err := json.NewDecoder(resp.Body).Decode(&rb); err != nil {
		return "", fmt.Errorf("can't decode response: %w", err)
	}

	return rb.Value, nil
}

// webElementIdentifier is the string constant "element-6066-11e4-a52e-4f735466cecf".
// https://www.w3.org/TR/webdriver/#elements
// The old WebDriver JSON protocol uses `ELEMENT` key.
const webElementIdentifier = "element-6066-11e4-a52e-4f735466cecf"

// FindElement finds the element user wants to get.
// https://www.w3.org/TR/webdriver/#find-element
func (b *Browser) FindElement(locator LocatorStrategy, value string) (ElementID, error) {
	type rp struct {
		Using LocatorStrategy `json:"using"`
		Value string          `json:"value"`
	}
	rpb := rp{
		Using: locator,
		Value: value,
	}
	body, err := json.Marshal(rpb)
	if err != nil {
		return "", fmt.Errorf("can't marshal json body: %w", err)
	}

	u := b.driver.RemoteEndURL.String() + "/session/" + string(b.SessionID) + "/element"
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := b.driver.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("got http response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Fixme: define the struct of error response and handle it.
		// https://www.w3.org/TR/webdriver/#errors
		return "", fmt.Errorf("got invalid http status code: %d", resp.StatusCode)
	}

	/**
	 * Response format is like this.
	 *
	 * {
	 *  "value": {
	 *      "element-6066-11e4-a52e-4f735466cecf": "84b10d39-94f5-4768-8457-dd218597a1e5"
	 *  }
	 * }
	 */
	type rf struct {
		Value map[string]string `json:"value"`
	}
	var rb rf
	if err := json.NewDecoder(resp.Body).Decode(&rb); err != nil {
		return "", fmt.Errorf("can't decode response: %w", err)
	}

	eID, ok := rb.Value[webElementIdentifier]
	if !ok {
		return "", fmt.Errorf("got empty element ID: value %#v", rb.Value)
	}

	return ElementID(eID), nil
}

// ElementID is an identifier of elements (ex. 84b10d39-94f5-4768-8457-dd218597a1e5).
// https://www.w3.org/TR/webdriver/#elements
type ElementID string

// LocatorStrategy is the keyword used to search for elements in the current browsing context.
// https://www.w3.org/TR/webdriver/#locator-strategies
type LocatorStrategy string

const (
	CSS             LocatorStrategy = "css selector"
	LinkText        LocatorStrategy = "link text"
	PartialLinkText LocatorStrategy = "partial link text"
	TagName         LocatorStrategy = "tag name"
	Xpath           LocatorStrategy = "xpath"
)
