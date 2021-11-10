package gowd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
	u := fmt.Sprintf("%s/session/%s", b.driver.RemoteEndURL.String(), string(b.SessionID))
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

// FindElement finds the element user wants to get.
// https://www.w3.org/TR/webdriver/#find-element
func (b *Browser) FindElement(locator LocatorStrategy, value string) (*Element, error) {
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
		return &Element{}, fmt.Errorf("can't marshal json body: %w", err)
	}

	u := b.driver.RemoteEndURL.String() + "/session/" + string(b.SessionID) + "/element"
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return &Element{}, fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := b.driver.Client.Do(req)
	if err != nil {
		return &Element{}, fmt.Errorf("got http response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Fixme: define the struct of error response and handle it.
		// https://www.w3.org/TR/webdriver/#errors
		return &Element{}, fmt.Errorf("got invalid http status code: %d", resp.StatusCode)
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
		return &Element{}, fmt.Errorf("can't decode response: %w", err)
	}

	eID, ok := rb.Value[webElementIdentifier]
	if !ok {
		return &Element{}, fmt.Errorf("got empty element ID: value %#v", rb.Value)
	}

	return NewElement(ElementID(eID), b.driver, b), nil
}

// TakeScreenshot takes a screenshot of current browser.
// https://www.w3.org/TR/webdriver/#take-screenshot
func (b *Browser) TakeScreenshot() ([]byte, error) {
	u := fmt.Sprintf(
		"%s/session/%s/screenshot",
		b.driver.RemoteEndURL.String(),
		string(b.SessionID),
	)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := b.driver.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("got http response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Fixme: define the struct of error response and handle it.
		// https://www.w3.org/TR/webdriver/#errors
		return nil, fmt.Errorf("got invalid http status code: %d", resp.StatusCode)
	}

	/**
	 * Response format is like this.
	 *
	 * {
	 *  "value": "iVBORw0KGgoAAAANSUhEUgAABkAAAASwCAYAAACjAYaXAAABLWlDQ1BTa2lhA...(omit)"
	 * }
	 */
	type rf struct {
		Value string `json:"value"`
	}
	var rb rf
	if err := json.NewDecoder(resp.Body).Decode(&rb); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	// WebDriver returns a base64 encoded image.
	// https://www.w3.org/TR/webdriver/#take-screenshot
	// > Let encoding result be the result of trying encoding a canvas as Base64 canvas.
	// StdEncoding is the standard base64 encoding, as defined in RFC 4648.
	bt, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(rb.Value)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return bt, nil
}

// NewWindow opens a new window tab.
func (b *Browser) NewWindow() error {
	// https://www.w3.org/TR/webdriver/#new-window
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/session/%s/window/new", b.driver.RemoteEndURL.String(), string(b.SessionID)),
		strings.NewReader("{}"),
	)
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
