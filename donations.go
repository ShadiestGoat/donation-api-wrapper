package donations

import (
	"net/http"
	"net/url"
)

type Donation struct {
	// The ID of the donation
	ID string `json:"id"`
	// The paypal order id
	OrderID string `json:"ppOrderID"`
	// The paypal capture id
	CaptureID string `json:"ppCaptureID"`
	// The id of the donor that made this donation
	Donor string `json:"donor"`
	// The message attached to this donation
	Message string `json:"message"`
	// The amount that this person donated (the amount that was sent, not the amount that was received!)
	Amount float64 `json:"amount"`
	// The fund that this was donated to
	FundID string `json:"fundID"`
}

// Fetch donations. The before & after values are used for pagination.
func (c *Client) Donations(before, after string) ([]*Donation, error) {
	resp := []*Donation{}
	q := url.Values{}
	if before != "" {
		q.Set("before", before)
	}
	if after != "" {
		q.Set("after", after)
	}

	err := c.fetch(http.MethodGet, `/donations?`+q.Encode(), nil, &resp)
	return withErrorArr(resp, err)
}
