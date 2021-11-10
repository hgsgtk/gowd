package gowd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ElementID is an identifier of elements (ex. 84b10d39-94f5-4768-8457-dd218597a1e5).
// https://www.w3.org/TR/webdriver/#elements
type ElementID string

// Element is an identified element by Browser.
type Element struct {
	ID      ElementID
	browser *Browser
	driver  *WebDriver
}

// NewElement is the constructor of Element.
func NewElement(ID ElementID, driver *WebDriver, browser *Browser) *Element {
	return &Element{ID: ID, driver: driver, browser: browser}
}

// GetText gets the text of Element.
// https://www.w3.org/TR/webdriver/#get-element-text
func (e *Element) GetText() (string, error) {
	u := e.driver.RemoteEndURL.String() +
		"/session/" + string(e.browser.SessionID) +
		"/element/" + string(e.ID) + "/text"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := e.driver.Client.Do(req)
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

// Click clicks the element.
// https://www.w3.org/TR/webdriver/#element-click
func (e *Element) Click() error {
	u := fmt.Sprintf(
		"%s/session/%s/element/%s/click",
		e.driver.RemoteEndURL.String(),
		string(e.browser.SessionID),
		string(e.ID),
	)
	// The request body should be JSON empty object.
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader("{}"))
	if err != nil {
		return fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := e.driver.Client.Do(req)
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

// TakeScreenshot takes a screenshot of current element.
// https://www.w3.org/TR/webdriver/#take-element-screenshot
func (e *Element) TakeScreenshot() ([]byte, error) {
	u := fmt.Sprintf("%s/session/%s/element/%s/screenshot",
		e.driver.RemoteEndURL.String(),
		string(e.browser.SessionID),
		string(e.ID),
	)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create a request: %w", err)
	}

	resp, err := e.driver.Client.Do(req)
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

	// Fixme: It is implemented similarly to Browser.TakeScreenshot.
	// 	I want to refactor it to be common to one.
	// WebDriver returns a base64 encoded image.
	// https://www.w3.org/TR/webdriver/#take-element-screenshot
	// > Let encoding result be the result of trying encoding a canvas as Base64 canvas.
	// StdEncoding is the standard base64 encoding, as defined in RFC 4648.
	bt, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(rb.Value)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return bt, nil
}
