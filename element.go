package gowd

import (
	"encoding/json"
	"fmt"
	"net/http"
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
