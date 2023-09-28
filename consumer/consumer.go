package consumer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	host string
}

func NewClient(host string) *Client {
	return &Client{
		host: host,
	}
}

type Discount struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Value       float64 `json:"value"`
}

func (c *Client) GetDiscount(id int) (*Discount, error) {
	url := fmt.Sprintf("%s/discounts/%d", c.host, id)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errNotFound
	}

	b, _ := io.ReadAll(resp.Body)

	var d Discount
	_ = json.Unmarshal(b, &d)

	return &d, nil
}

func (c *Client) GetDiscounts(filter map[string]string) ([]Discount, error) {
	url := fmt.Sprintf("%s/discounts", c.host)

	if filter["type"] != "" {
		url += "?type=" + filter["type"]
	}

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type response struct {
		Discounts []Discount `json:"discounts"`
	}

	b, _ := io.ReadAll(resp.Body)

	var r response
	_ = json.Unmarshal(b, &r)

	return r.Discounts, nil
}

func (c *Client) PutDiscount(d Discount) error {
	url := fmt.Sprintf("%s/discounts/%d", c.host, d.ID)

	var body struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Type        string  `json:"type"`
		Value       float64 `json:"value"`
	}

	body.Title = d.Title
	body.Description = d.Description
	body.Type = d.Type
	body.Value = d.Value

	rawBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", url, bytes.NewReader(rawBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errNotFound
	}

	return nil
}

func (c *Client) CreateDiscount(d Discount) error {
	url := fmt.Sprintf("%s/discounts", c.host)

	var body struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Type        string  `json:"type"`
		Value       float64 `json:"value"`
	}

	body.Title = d.Title
	body.Description = d.Description
	body.Type = d.Type
	body.Value = d.Value

	rawBody, _ := json.Marshal(body)

	resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewReader(rawBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return errInvalidRequest
	}

	return nil
}

func (c *Client) DeleteDiscount(id int) error {
	url := fmt.Sprintf("%s/discounts/%d", c.host, id)

	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return errNotFound
	}

	return nil
}

var (
	errNotFound       = errors.New("discount not found")
	errInvalidRequest = errors.New("invalid request")
)
