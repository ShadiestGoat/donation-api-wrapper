package donations

import (
	"net/http"
	"net/url"
)

type Fund struct {
	// The ID of the fund
	ID         string   `json:"id,omitempty"`
	// If this fund is the default or not. If Default, then it is what will be displayed on the main page.
	Default    *bool    `json:"default,omitempty"`
	// The goal of this fund. If 0, then there is no goal.
	Goal       float64  `json:"goal,omitempty"`
	// An alias for the fund. This is a name that will be used in the 'quick' url, ie. under /f/{alias}
	Alias      string   `json:"alias"`
	// A short fund title, in the format of 'Donate for {ShortTitle}'
	ShortTitle string   `json:"shortTitle"`
	// The official big title on the main page
	Title      string   `json:"title"`
	// The amount that has been received for this fund
	Amount     *float64 `json:"amount,omitempty"`
}

// Fetch funds, with query options.
// before, after are options for pagination
// if fetchAmounts, then it will fetch the amount donated for this fund 
func (c *Client) Funds(before, after string, fetchAmounts bool, complete *bool) ([]*Fund, error) {
	resp := []*Fund{}
	q := url.Values{}

	if before != "" {
		q.Set("before", before)
	}
	if after != "" {
		q.Set("after", after)
	}
	amt := "f"
	if fetchAmounts {
		amt = "t"
	}
	q.Set("amount", amt)
	if complete != nil {
		comp := "f"
		if *complete {
			comp = "t"
		}
		q.Set("complete", comp)
	}

	err := c.fetch(http.MethodGet, `/funds?` + q.Encode(), nil, &resp)

	return withErrorArr(resp, err)
}

// Create a new fund. For explanations on each of the arguments, see documentation for Fund
func (c *Client) NewFund(alias, shortTitle, title string, def bool, goal float64) (*Fund, error) {
	body := &Fund{
		Default:    &def,
		Goal:       goal,
		Alias:      alias,
		ShortTitle: shortTitle,
		Title:      title,
	}
	resp := &Fund{}
	err := c.fetch(http.MethodPost, `/funds`, body, &resp)

	return withError(resp, err)
}

func (c *Client) FundByID(id string) (*Fund, error) {
	resp := &Fund{}
	
	err := c.fetch(http.MethodGet, `/funds/` + id, nil, resp)

	return withError(resp, err)
}

func (c *Client) UpdateFund(f *Fund) error {
	if f == nil {
		return nil
	}

	return c.fetch(http.MethodPut, `/funds/` + f.ID, f, nil)
}

func (c *Client) MakeFundDefault(id string) (*Fund, error) {
	resp := &Fund{}
	
	err := c.fetch(http.MethodPut, `/funds/` + id, nil, &resp)

	return withError(resp, err)
}
